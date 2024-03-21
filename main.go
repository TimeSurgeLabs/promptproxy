package main

import (
	"context"
	"encoding/json"
	"errors"
	"log"
	"os"
	"strings"

	"github.com/labstack/echo/v5"
	"github.com/pocketbase/dbx"
	"github.com/pocketbase/pocketbase"
	"github.com/pocketbase/pocketbase/core"
	"github.com/pocketbase/pocketbase/models"
	"github.com/pocketbase/pocketbase/plugins/migratecmd"
	openai "github.com/sashabaranov/go-openai"

	_ "github.com/TimeSurgeLabs/promptproxy/migrations"
)

func ProcessCompletionPrompt(req *openai.CompletionRequest, systemPrompt string) error {
	// req.Prompt can be a string or an array of strings
	// if its a string, prepend the systemPrompt to it
	// if its an array, prepend the systemPrompt to each element
	// then reassign it to req.Prompt
	// if its neither, return an error

	switch prompt := req.Prompt.(type) {
	case string:
		req.Prompt = systemPrompt + "\n\n" + prompt
	case []string:
		for i, p := range prompt {
			prompt[i] = systemPrompt + "\n\n" + p
		}
		req.Prompt = prompt
	default:
		return errors.New("invalid prompt type")
	}

	return nil
}

func ProcessChatCompletionPrompt(req *openai.ChatCompletionRequest, systemPrompt string) {
	// loop through messages, check if a prompt with role "system" exists
	// if it does, replace with the systemPrompt
	// if it doesn't, prepend the systemPrompt to the messages

	for i, message := range req.Messages {
		if message.Role == "system" {
			req.Messages[i].Content = systemPrompt
			return
		}
	}

	// add it to the beginning of the messages
	req.Messages = append([]openai.ChatCompletionMessage{{Role: "system", Content: systemPrompt}}, req.Messages...)
}

func main() {
	app := pocketbase.New()

	// loosely check if it was executed using "go run"
	isGoRun := strings.HasPrefix(os.Args[0], os.TempDir())

	migratecmd.MustRegister(app, app.RootCmd, migratecmd.Config{
		// enable auto creation of migration files when making collection changes in the Admin UI
		// (the isGoRun check is to enable it only during development)
		Automigrate: isGoRun,
	})

	app.OnBeforeServe().Add(func(e *core.ServeEvent) error {
		// e.Router.GET("/*", apis.StaticDirectoryHandler(os.DirFS("./public"), false))
		apiKeyMiddleware := func(next echo.HandlerFunc) echo.HandlerFunc {
			return func(c echo.Context) error {
				apiKey := c.Request().Header.Get("Authorization")
				if apiKey == "" {
					return c.JSON(401, map[string]interface{}{
						"error": "missing api key",
					})
				}

				// remove "Bearer " token
				apiKey = strings.TrimPrefix(apiKey, "Bearer ")
				// check if the api key exists in the database
				apiRecord, err := app.Dao().FindRecordsByExpr("keys", dbx.HashExp{
					"key": apiKey,
				})
				if err != nil {
					return err
				}

				if len(apiRecord) == 0 {
					return c.JSON(401, map[string]interface{}{
						"error": "invalid api key",
					})
				}

				return next(c)
			}
		}
		// model is of format model:prompt
		// prompt is optional, api and model are required.
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
				err = ProcessCompletionPrompt(&req, systemPrompt)
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
		}, apiKeyMiddleware)

		e.Router.POST("/v1/chat/completions", func(c echo.Context) error {
			var req openai.ChatCompletionRequest

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
				ProcessChatCompletionPrompt(&req, systemPrompt)
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

			// dump to json and print
			b, err := json.Marshal(req)
			if err != nil {
				return err
			}

			log.Println(string(b))

			resp, err := client.CreateChatCompletion(context.TODO(), req)
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
		}, apiKeyMiddleware)

		return nil
	})

	if err := app.Start(); err != nil {
		log.Fatal(err)
	}
}
