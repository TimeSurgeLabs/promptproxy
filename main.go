package main

import (
	"log"
	"os"
	"strings"

	"github.com/pocketbase/pocketbase"
	"github.com/pocketbase/pocketbase/plugins/migratecmd"

	_ "github.com/TimeSurgeLabs/promptproxy/migrations"
	"github.com/TimeSurgeLabs/promptproxy/routes"
)

func main() {
	app := pocketbase.New()

	// loosely check if it was executed using "go run"
	isGoRun := strings.HasPrefix(os.Args[0], os.TempDir())

	migratecmd.MustRegister(app, app.RootCmd, migratecmd.Config{
		// enable auto creation of migration files when making collection changes in the Admin UI
		// (the isGoRun check is to enable it only during development)
		Automigrate: isGoRun,
	})

	routes.BindChatCompletionRoute(app)
	routes.BindCompletionRoute(app)

	if err := app.Start(); err != nil {
		log.Fatal(err)
	}
}
