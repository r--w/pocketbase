package migrations

import (
	"database/sql"
	"errors"
	"log"
	"strings"

	"github.com/pocketbase/dbx"
	"github.com/pocketbase/pocketbase/daos"
	m "github.com/pocketbase/pocketbase/migrations"
	"github.com/pocketbase/pocketbase/models"
)

func init() {
	m.Register(func(db dbx.Builder) error {
		dao := daos.New(db)
		collection, err := dao.FindCollectionByNameOrId("users")
		if err != nil {
			return err
		}

		_, err = dao.FindFirstRecordByData(collection.Name, "email", UserEmailPassword)
		exists := true
		if err != nil {
			if errors.Is(err, sql.ErrNoRows) {
				exists = false
			} else {
				return err
			}
		}

		if exists {
			return nil
		}

		log.Println("inserting normal user: ", UserEmailPassword)

		r := models.NewRecord(collection)
		if err := r.SetEmail(UserEmailPassword); err != nil {
			return err
		}
		if err := r.SetUsername(strings.Split(UserEmailPassword, "@")[0]); err != nil {
			return err
		}
		if err := r.SetVerified(true); err != nil {
			return err
		}
		if err := r.SetPassword(UserEmailPassword); err != nil {
			return err
		}

		if err := dao.SaveRecord(r); err != nil {
			return err
		}

		return nil
	}, func(db dbx.Builder) error {
		return nil
	})
}
