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

		collection, err := dao.FindCollectionByNameOrId("2lsx7ksvpbih4cu")
		if err != nil {
			return err
		}

		// add
		new_model := &schema.SchemaField{}
		if err := json.Unmarshal([]byte(`{
			"system": false,
			"id": "1w3owj09",
			"name": "model",
			"type": "relation",
			"required": true,
			"presentable": false,
			"unique": false,
			"options": {
				"collectionId": "k10mgnjjc2zafen",
				"cascadeDelete": false,
				"minSelect": null,
				"maxSelect": 1,
				"displayFields": null
			}
		}`), new_model); err != nil {
			return err
		}
		collection.Schema.AddField(new_model)

		// add
		new_system_prompt := &schema.SchemaField{}
		if err := json.Unmarshal([]byte(`{
			"system": false,
			"id": "slhdtvqo",
			"name": "system_prompt",
			"type": "text",
			"required": false,
			"presentable": false,
			"unique": false,
			"options": {
				"min": null,
				"max": null,
				"pattern": ""
			}
		}`), new_system_prompt); err != nil {
			return err
		}
		collection.Schema.AddField(new_system_prompt)

		return dao.SaveCollection(collection)
	}, func(db dbx.Builder) error {
		dao := daos.New(db);

		collection, err := dao.FindCollectionByNameOrId("2lsx7ksvpbih4cu")
		if err != nil {
			return err
		}

		// remove
		collection.Schema.RemoveField("1w3owj09")

		// remove
		collection.Schema.RemoveField("slhdtvqo")

		return dao.SaveCollection(collection)
	})
}
