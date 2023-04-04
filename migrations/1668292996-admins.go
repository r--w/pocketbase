package migrations

import (
	"database/sql"
	"errors"
	"log"

	"github.com/pocketbase/dbx"
	"github.com/pocketbase/pocketbase/daos"
	m "github.com/pocketbase/pocketbase/migrations"
	"github.com/pocketbase/pocketbase/models"
)

func init() {
	m.Register(func(db dbx.Builder) error {
		dao := daos.New(db)

		exists := true
		_, err := dao.FindAdminByEmail(AdminEmailPassword)
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

		log.Println("inserting admin user: ", AdminEmailPassword)

		admin := models.Admin{
			Email: AdminEmailPassword,
		}

		if err := admin.SetPassword(AdminEmailPassword); err != nil {
			return err
		}

		return dao.SaveAdmin(&admin)
	}, func(db dbx.Builder) error {
		return nil
	})
}
