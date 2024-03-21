package routes

import (
	"context"
	"strings"

	"github.com/labstack/echo/v5"
	"github.com/pocketbase/dbx"
	"github.com/pocketbase/pocketbase"
	"github.com/pocketbase/pocketbase/core"
	"github.com/pocketbase/pocketbase/models"
	"github.com/sashabaranov/go-openai"

	"github.com/TimeSurgeLabs/promptproxy/middleware"
	"github.com/TimeSurgeLabs/promptproxy/utils"
)

func BindCompletionRoute(app *pocketbase.PocketBase) {
	app.OnBeforeServe().Add(func(e *core.ServeEvent) error {
		e.Router.POST("/v1/completions", func(c echo.Context) error {
			var req openai.CompletionRequest
			if err := c.Bind(&req); err != nil {
				return err
			}

			// parse the req.Model field
			parts := strings.Split(req.Model, ":")

			// if the length is less than 1, return an error
			if len(parts) < 1 {
				return c.JSON(400, map[string]interface{}{
					"error": "invalid model format. Must be of format api:model:prompt",
				})
			}

			systemPrompt := ""
			// if a prompt is present, query from the database
			if len(parts) == 2 {
				promptRecord, err := app.Dao().FindRecordById("prompts", parts[1])
				if err != nil {
					return err
				}

				// if the record is not found, return an error
				if promptRecord == nil {
					return c.JSON(404, map[string]interface{}{
						"error": "prompt not found",
					})
				}

				systemPrompt = promptRecord.GetString("instructions") + "\n\n"
				err = utils.ProcessCompletionPrompt(&req, systemPrompt)
				if err != nil {
					return c.JSON(400, map[string]interface{}{
						"error": err.Error(),
					})
				}
			}

			modelId := parts[0]
			// first get the model name
			// then get the api as the relation "api" from the model
			modelRecord, err := app.Dao().FindRecordById("models", modelId)
			if err != nil {
				return err
			}

			// if the record is not found, return an error
			if modelRecord == nil {
				return c.JSON(404, map[string]interface{}{
					"error": "model not found",
				})
			}

			apiId := modelRecord.GetString("api")
			apiRecord, err := app.Dao().FindRecordById("apis", apiId)
			if err != nil {
				return err
			}

			if apiRecord == nil {
				return c.JSON(404, map[string]interface{}{
					"error": "api not found",
				})
			}

			apiUrl := apiRecord.GetString("url")
			apiKey := apiRecord.GetString("api_key")
			model := modelRecord.GetString("model")

			conf := openai.DefaultConfig(apiKey)
			conf.BaseURL = apiUrl
			conf.HTTPClient.Timeout = 0
			// create a new openai client with the given url and key
			client := openai.NewClientWithConfig(conf)

			req.Model = model

			if req.Stream {
				return c.JSON(400, map[string]interface{}{
					"error": "streaming not supported",
				})
			}

			resp, err := client.CreateCompletion(context.TODO(), req)
			if err != nil {
				return c.JSON(500, map[string]interface{}{
					"error": err.Error(),
				})
			}

			// get the api key from the header
			apiKey = c.Request().Header.Get("Authorization")

			keyRecords, err := app.Dao().FindRecordsByExpr("keys", dbx.HashExp{
				"key": apiKey,
			})

			if err != nil {
				return err
			}

			if len(keyRecords) == 0 {
				return c.JSON(401, map[string]interface{}{
					"error": "invalid api key",
				})
			}

			keyRecord := keyRecords[0]
			userId := keyRecord.GetString("user")

			collection, err := app.Dao().FindCollectionByNameOrId("requests")
			if err != nil {
				return err
			}

			// add a new record to the requests table
			record := models.NewRecord(collection)
			record.Set("api", apiId)
			record.Set("user", userId)
			if len(parts) == 2 {
				record.Set("prompt", parts[1])
			}
			record.Set("input_tokens", resp.Usage.PromptTokens)
			record.Set("output_tokens", resp.Usage.CompletionTokens)

			if err := app.Dao().SaveRecord(record); err != nil {
				return err
			}

			return c.JSON(200, resp)
		}, middleware.APIKey(app))

		return nil
	})
}
