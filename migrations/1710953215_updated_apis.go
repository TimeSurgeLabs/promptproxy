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

		collection, err := dao.FindCollectionByNameOrId("0zmba94robz8v5t")
		if err != nil {
			return err
		}

		// remove
		collection.Schema.RemoveField("fmr0gswi")

		return dao.SaveCollection(collection)
	}, func(db dbx.Builder) error {
		dao := daos.New(db);

		collection, err := dao.FindCollectionByNameOrId("0zmba94robz8v5t")
		if err != nil {
			return err
		}

		// add
		del_model := &schema.SchemaField{}
		if err := json.Unmarshal([]byte(`{
			"system": false,
			"id": "fmr0gswi",
			"name": "model",
			"type": "text",
			"required": true,
			"presentable": false,
			"unique": false,
			"options": {
				"min": null,
				"max": null,
				"pattern": ""
			}
		}`), del_model); err != nil {
			return err
		}
		collection.Schema.AddField(del_model)

		return dao.SaveCollection(collection)
	})
}
