package main

import (
	"log"
	"time"

	"pocketbase"

	"github.com/hashicorp/go-multierror"
	"github.com/mitchellh/mapstructure"
)

type Post struct {
	Field string `json:"field"`
	ID    string `json:"id"`
}

func main() {
	// TODO add documentation how to run local Pocketbase server
	// TODO example with custom options
	// TODO example with Params
	// TODO add documentation about mapstructure usage here

	var errors error
	client := pocketbase.NewClient("http://localhost:8090",
		pocketbase.WithAdminEmailPassword("admin1@admin.com", "admin@admin.com"),
		// pocketbase.WithDebug(),
	)
	// client := pocketbase.NewClient("http://localhost:8090")
	response, err := client.List("posts_public", pocketbase.Params{
		Size:    1,
		Page:    1,
		Sort:    "-created",
		Filters: "field~'test'",
	})
	errors = multierror.Append(errors, err)

	log.Printf("Total items: %d, total pages: %d\n", response.TotalItems, response.TotalPages)
	for _, item := range response.Items {
		var test Post
		err := mapstructure.Decode(item, &test)
		errors = multierror.Append(errors, err)

		log.Printf("Item: %#v\n", test)
	}

	log.Println("Inserting new item")
	// use can use struct type - just make sure it has json tags
	err = client.Create("posts_public", Post{
		Field: "test_" + time.Now().Format(time.Stamp),
	})
	errors = multierror.Append(errors, err)

	// or use map[string]interface{}
	err = client.Create("posts_public", map[string]interface{}{
		"field": "test_" + time.Now().Format(time.Stamp),
	})

	if errors != nil {
		log.Fatal(errors)
	}
}
