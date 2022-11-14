package migrations

import (
	"log"

	"github.com/pocketbase/dbx"
	"github.com/pocketbase/pocketbase/daos"
	m "github.com/pocketbase/pocketbase/migrations"
	"github.com/pocketbase/pocketbase/models"
)

func init() {
	m.Register(func(db dbx.Builder) error {
		dao := daos.New(db)
		for _, c := range []string{PostsAdmin, PostsUser, PostsPublic} {
			collection, err := dao.FindCollectionByNameOrId(c)
			if err != nil {
				return err
			}

			log.Println("inserting post to: ", c)

			r := models.NewRecord(collection)
			r.Set("field", "test")

			if err := dao.SaveRecord(r); err != nil {
				return err
			}
		}

		return nil
	}, func(db dbx.Builder) error {
		return nil
	})
}
