package middleware

import (
	"strings"

	"github.com/labstack/echo/v5"
	"github.com/pocketbase/dbx"
	"github.com/pocketbase/pocketbase"
)

func APIKey(app *pocketbase.PocketBase) func(next echo.HandlerFunc) echo.HandlerFunc {
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

	return apiKeyMiddleware
}
