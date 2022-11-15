package main

import (
	"log"
	"time"

	"github.com/mitchellh/mapstructure"
	"github.com/r--w/pocketbase"
	"go.uber.org/multierr"
)

type Post struct {
	Field string `json:"field"`
	ID    string `json:"id"`
}

func main() {
	// TODO add documentation how to run local Pocketbase server
	var errors error
	client := pocketbase.NewClient("http://localhost:8090")
	// Config options:
	// pocketbase.WithAdminEmailPassword("admin@admin.com", "admin@admin.com")
	// pocketbase.WithUserEmailPassword("user@user.com", "user@user.com")
	// pocketbase.WithDebug()

	response, err := client.List("posts_public", pocketbase.Params{
		Size:    1,
		Page:    1,
		Sort:    "-created",
		Filters: "field~'test'",
	})
	errors = multierr.Append(errors, err)

	log.Printf("Total items: %d, total pages: %d\n", response.TotalItems, response.TotalPages)
	for _, item := range response.Items {
		var test Post
		err := mapstructure.Decode(item, &test)
		errors = multierr.Append(errors, err)

		log.Printf("Item: %#v\n", test)
	}

	log.Println("Inserting new item")
	// use can use struct type - just make sure it has json tags
	err = client.Create("posts_public", Post{
		Field: "test_" + time.Now().Format(time.Stamp),
	})
	errors = multierr.Append(errors, err)

	// or use simple map[string]interface{}
	err = client.Create("posts_public", map[string]interface{}{
		"field": "test_" + time.Now().Format(time.Stamp),
	})
	errors = multierr.Append(errors, err)

	if errors != nil {
		log.Fatal(errors)
	}
}
