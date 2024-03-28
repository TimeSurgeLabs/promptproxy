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

		collection, err := dao.FindCollectionByNameOrId("5sx1sxlg7nnhija")
		if err != nil {
			return err
		}

		// add
		new_track := &schema.SchemaField{}
		if err := json.Unmarshal([]byte(`{
			"system": false,
			"id": "c5wbtdcm",
			"name": "track",
			"type": "bool",
			"required": false,
			"presentable": false,
			"unique": false,
			"options": {}
		}`), new_track); err != nil {
			return err
		}
		collection.Schema.AddField(new_track)

		return dao.SaveCollection(collection)
	}, func(db dbx.Builder) error {
		dao := daos.New(db);

		collection, err := dao.FindCollectionByNameOrId("5sx1sxlg7nnhija")
		if err != nil {
			return err
		}

		// remove
		collection.Schema.RemoveField("c5wbtdcm")

		return dao.SaveCollection(collection)
	})
}
