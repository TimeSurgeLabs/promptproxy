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

		// remove
		collection.Schema.RemoveField("5jowoxa3")

		// add
		new_user_prompt := &schema.SchemaField{}
		if err := json.Unmarshal([]byte(`{
			"system": false,
			"id": "6omyowgu",
			"name": "user_prompt",
			"type": "text",
			"required": true,
			"presentable": false,
			"unique": false,
			"options": {
				"min": null,
				"max": null,
				"pattern": ""
			}
		}`), new_user_prompt); err != nil {
			return err
		}
		collection.Schema.AddField(new_user_prompt)

		// add
		new_assistant_response := &schema.SchemaField{}
		if err := json.Unmarshal([]byte(`{
			"system": false,
			"id": "rfrsxv1w",
			"name": "assistant_response",
			"type": "text",
			"required": true,
			"presentable": false,
			"unique": false,
			"options": {
				"min": null,
				"max": null,
				"pattern": ""
			}
		}`), new_assistant_response); err != nil {
			return err
		}
		collection.Schema.AddField(new_assistant_response)

		// update
		edit_input := &schema.SchemaField{}
		if err := json.Unmarshal([]byte(`{
			"system": false,
			"id": "4dyabmvp",
			"name": "input",
			"type": "json",
			"required": true,
			"presentable": false,
			"unique": false,
			"options": {
				"maxSize": 2000000
			}
		}`), edit_input); err != nil {
			return err
		}
		collection.Schema.AddField(edit_input)

		// update
		edit_output := &schema.SchemaField{}
		if err := json.Unmarshal([]byte(`{
			"system": false,
			"id": "lnwaukln",
			"name": "output",
			"type": "json",
			"required": true,
			"presentable": false,
			"unique": false,
			"options": {
				"maxSize": 2000000
			}
		}`), edit_output); err != nil {
			return err
		}
		collection.Schema.AddField(edit_output)

		return dao.SaveCollection(collection)
	}, func(db dbx.Builder) error {
		dao := daos.New(db);

		collection, err := dao.FindCollectionByNameOrId("2lsx7ksvpbih4cu")
		if err != nil {
			return err
		}

		// add
		del_api := &schema.SchemaField{}
		if err := json.Unmarshal([]byte(`{
			"system": false,
			"id": "5jowoxa3",
			"name": "api",
			"type": "relation",
			"required": true,
			"presentable": false,
			"unique": false,
			"options": {
				"collectionId": "0zmba94robz8v5t",
				"cascadeDelete": false,
				"minSelect": null,
				"maxSelect": 1,
				"displayFields": null
			}
		}`), del_api); err != nil {
			return err
		}
		collection.Schema.AddField(del_api)

		// remove
		collection.Schema.RemoveField("6omyowgu")

		// remove
		collection.Schema.RemoveField("rfrsxv1w")

		// update
		edit_input := &schema.SchemaField{}
		if err := json.Unmarshal([]byte(`{
			"system": false,
			"id": "4dyabmvp",
			"name": "input",
			"type": "json",
			"required": false,
			"presentable": false,
			"unique": false,
			"options": {
				"maxSize": 2000000
			}
		}`), edit_input); err != nil {
			return err
		}
		collection.Schema.AddField(edit_input)

		// update
		edit_output := &schema.SchemaField{}
		if err := json.Unmarshal([]byte(`{
			"system": false,
			"id": "lnwaukln",
			"name": "output",
			"type": "json",
			"required": false,
			"presentable": false,
			"unique": false,
			"options": {
				"maxSize": 2000000
			}
		}`), edit_output); err != nil {
			return err
		}
		collection.Schema.AddField(edit_output)

		return dao.SaveCollection(collection)
	})
}
