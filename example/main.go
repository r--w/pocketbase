package main

import (
	"errors"
	"log"
	"time"

	"github.com/mitchellh/mapstructure"
	"github.com/r--w/pocketbase"
)

type Post struct {
	Field string `json:"field"`
	ID    string `json:"id"`
}

func main() {
	// REMEMBER to start the Pocketbase before running this example with `make serve` command
	var errs error
	client := pocketbase.NewClient("http://localhost:8090",
		pocketbase.WithAdminEmailPassword("admin@admin.com", "admin@admin.com"),
		pocketbase.WithDebug(),
	)

	client.Authorize()

	client = pocketbase.NewClient("http://localhost:8090",
		pocketbase.WithAdminToken(client.AuthStore().Token()),
		pocketbase.WithDebug(),
	)

	// Other configuration options:
	// pocketbase.WithAdminEmailPassword("admin@admin.com", "admin@admin.com")
	// pocketbase.WithUserEmailPassword("user@user.com", "user@user.com")
	// pocketbase.WithDebug()
	response, err := client.List("posts_public", pocketbase.ParamsList{
		Size:    1,
		Page:    1,
		Sort:    "-created",
		Filters: "field~'test'",
	})

	errs = errors.Join(errs, err)

	log.Printf("Total items: %d, total pages: %d\n", response.TotalItems, response.TotalPages)
	for _, item := range response.Items {
		var test Post
		err := mapstructure.Decode(item, &test)
		errs = errors.Join(errs, err)

		log.Printf("Item: %#v\n", test)
	}

	log.Println("Inserting new item")
	// you can use struct type - just make sure it has JSON tags
	_, err = client.Create("posts_public", Post{
		Field: "test_" + time.Now().Format(time.Stamp),
	})
	errs = errors.Join(errs, err)

	// or you can use simple map[string]any
	r, err := client.Create("posts_public", map[string]any{
		"field": "test_" + time.Now().Format(time.Stamp),
	})
	errs = errors.Join(errs, err)

	err = client.Delete("posts_public", r.ID)
	errs = errors.Join(errs, err)

	if errs != nil {
		log.Fatal(errs)
	}
}
