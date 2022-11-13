package migrations

import (
	"encoding/json"

	"github.com/pocketbase/dbx"
	"github.com/pocketbase/pocketbase/daos"
	m "github.com/pocketbase/pocketbase/migrations"
	"github.com/pocketbase/pocketbase/models"
)

// Auto generated migration with the most recent collections configuration.
func init() {
	m.Register(func(db dbx.Builder) error {
		jsonData := `[
    {
        "id": "_pb_users_auth_",
        "name": "users",
        "type": "auth",
        "system": false,
        "schema": [
            {
                "id": "users_name",
                "name": "name",
                "type": "text",
                "system": false,
                "required": false,
                "unique": false,
                "options": {
                    "min": null,
                    "max": null,
                    "pattern": ""
                }
            },
            {
                "id": "users_avatar",
                "name": "avatar",
                "type": "file",
                "system": false,
                "required": false,
                "unique": false,
                "options": {
                    "maxSelect": 1,
                    "maxSize": 5242880,
                    "mimeTypes": [
                        "image/jpg",
                        "image/jpeg",
                        "image/png",
                        "image/svg+xml",
                        "image/gif"
                    ],
                    "thumbs": null
                }
            }
        ],
        "listRule": "id = @request.auth.id",
        "viewRule": "id = @request.auth.id",
        "createRule": "",
        "updateRule": "id = @request.auth.id",
        "deleteRule": "id = @request.auth.id",
        "options": {
            "allowEmailAuth": true,
            "allowOAuth2Auth": true,
            "allowUsernameAuth": true,
            "exceptEmailDomains": null,
            "manageRule": null,
            "minPasswordLength": 8,
            "onlyEmailDomains": null,
            "requireEmail": false
        }
    },
    {
        "id": "5lt8hsy53jvfhwt",
        "name": "posts_public",
        "type": "base",
        "system": false,
        "schema": [
            {
                "id": "nqfl3tgr",
                "name": "field",
                "type": "text",
                "system": false,
                "required": false,
                "unique": false,
                "options": {
                    "min": null,
                    "max": null,
                    "pattern": ""
                }
            }
        ],
        "listRule": "",
        "viewRule": "",
        "createRule": "",
        "updateRule": "",
        "deleteRule": "",
        "options": {}
    },
    {
        "id": "4ea6djy249fzqfg",
        "name": "posts_admin",
        "type": "base",
        "system": false,
        "schema": [
            {
                "id": "i6fk95ln",
                "name": "field",
                "type": "text",
                "system": false,
                "required": false,
                "unique": false,
                "options": {
                    "min": null,
                    "max": null,
                    "pattern": ""
                }
            }
        ],
        "listRule": null,
        "viewRule": null,
        "createRule": null,
        "updateRule": null,
        "deleteRule": null,
        "options": {}
    },
    {
        "id": "78h06xtuhv8dqzx",
        "name": "posts_user",
        "type": "base",
        "system": false,
        "schema": [
            {
                "id": "y0sordxe",
                "name": "field",
                "type": "text",
                "system": false,
                "required": false,
                "unique": false,
                "options": {
                    "min": null,
                    "max": null,
                    "pattern": ""
                }
            }
        ],
        "listRule": "@request.auth.id != \"\"",
        "viewRule": "@request.auth.id != \"\"",
        "createRule": "@request.auth.id != \"\"",
        "updateRule": "@request.auth.id != \"\"",
        "deleteRule": "@request.auth.id != \"\"",
        "options": {}
    }
]`

		collections := []*models.Collection{}
		if err := json.Unmarshal([]byte(jsonData), &collections); err != nil {
			return err
		}

		return daos.New(db).ImportCollections(collections, true, nil)
	}, func(db dbx.Builder) error {
		// no revert since the configuration on the environment, on which
		// the migration was executed, could have changed via the UI/API
		return nil
	})
}
