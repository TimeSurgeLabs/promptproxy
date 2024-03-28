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

		// update
		edit_assistant := &schema.SchemaField{}
		if err := json.Unmarshal([]byte(`{
			"system": false,
			"id": "n2udvkii",
			"name": "assistant",
			"type": "relation",
			"required": false,
			"presentable": false,
			"unique": false,
			"options": {
				"collectionId": "hfxmy2ruwi4uryw",
				"cascadeDelete": false,
				"minSelect": null,
				"maxSelect": 1,
				"displayFields": null
			}
		}`), edit_assistant); err != nil {
			return err
		}
		collection.Schema.AddField(edit_assistant)

		return dao.SaveCollection(collection)
	}, func(db dbx.Builder) error {
		dao := daos.New(db);

		collection, err := dao.FindCollectionByNameOrId("2lsx7ksvpbih4cu")
		if err != nil {
			return err
		}

		// update
		edit_assistant := &schema.SchemaField{}
		if err := json.Unmarshal([]byte(`{
			"system": false,
			"id": "n2udvkii",
			"name": "prompt",
			"type": "relation",
			"required": false,
			"presentable": false,
			"unique": false,
			"options": {
				"collectionId": "hfxmy2ruwi4uryw",
				"cascadeDelete": false,
				"minSelect": null,
				"maxSelect": 1,
				"displayFields": null
			}
		}`), edit_assistant); err != nil {
			return err
		}
		collection.Schema.AddField(edit_assistant)

		return dao.SaveCollection(collection)
	})
}
