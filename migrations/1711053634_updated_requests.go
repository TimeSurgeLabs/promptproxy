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
		new_input := &schema.SchemaField{}
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
		}`), new_input); err != nil {
			return err
		}
		collection.Schema.AddField(new_input)

		// add
		new_output := &schema.SchemaField{}
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
		}`), new_output); err != nil {
			return err
		}
		collection.Schema.AddField(new_output)

		return dao.SaveCollection(collection)
	}, func(db dbx.Builder) error {
		dao := daos.New(db);

		collection, err := dao.FindCollectionByNameOrId("2lsx7ksvpbih4cu")
		if err != nil {
			return err
		}

		// remove
		collection.Schema.RemoveField("4dyabmvp")

		// remove
		collection.Schema.RemoveField("lnwaukln")

		return dao.SaveCollection(collection)
	})
}
