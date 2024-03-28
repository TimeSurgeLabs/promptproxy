package migrations

import (
	"encoding/json"

	"github.com/pocketbase/dbx"
	"github.com/pocketbase/pocketbase/daos"
	m "github.com/pocketbase/pocketbase/migrations"
	"github.com/pocketbase/pocketbase/models/schema"
)

func init() {
	m.Register(func(db dbx.Builder) error {
		dao := daos.New(db);

		collection, err := dao.FindCollectionByNameOrId("hfxmy2ruwi4uryw")
		if err != nil {
			return err
		}

		collection.Name = "assistants"

		// add
		new_tools := &schema.SchemaField{}
		if err := json.Unmarshal([]byte(`{
			"system": false,
			"id": "ngt0ysze",
			"name": "tools",
			"type": "json",
			"required": false,
			"presentable": false,
			"unique": false,
			"options": {
				"maxSize": 2000000
			}
		}`), new_tools); err != nil {
			return err
		}
		collection.Schema.AddField(new_tools)

		return dao.SaveCollection(collection)
	}, func(db dbx.Builder) error {
		dao := daos.New(db);

		collection, err := dao.FindCollectionByNameOrId("hfxmy2ruwi4uryw")
		if err != nil {
			return err
		}

		collection.Name = "prompts"

		// remove
		collection.Schema.RemoveField("ngt0ysze")

		return dao.SaveCollection(collection)
	})
}
