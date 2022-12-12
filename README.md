[![Check & test & build](https://github.com/r--w/pocketbase/actions/workflows/main.yml/badge.svg)](https://github.com/r--w/pocketbase/actions/workflows/main.yml)
[![PocketBase](https://pocketbase.io/images/logo.svg)](https://pocketbase.io)

### Project
This repository contains community-maintained Go SDK for Pocketbase API.
It's well-tested and used in production in [Coinpaprika](https://coinpaprika.com), but not all endpoints are covered yet.

### Compatibility
* `v0.9.2` version of SDK is compatible with Pocketbase v0.9.x (SSE & generics support introduced)
* `v0.8.0` version of SDK is compatible with Pocketbase v0.8.x

### PocketBase
[Pocketbase](https://pocketbase.io) is a simple, self-hosted, open-source, no-code, database for your personal data.
It's a great alternative to Airtable, Notion, and Google Sheets. Source code is available on [GitHub](https://github.com/pocketbase/pocketbase)

### Currently supported operations
This SDK doesn't have feature parity with official SDKs and supports the following operations:

* **Authentication** - anonymous, admin and user via email/password
* **Create** 
* **Update**
* **Delete**
* **List** - with pagination, filtering, sorting
* **Other** - feel free to create an issue or contribute

### Usage & examples

Simple list example without authentication (assuming your collections are public):

```go
package main

import (
	"log"

	"github.com/r--w/pocketbase"
)

func main() {
	client := pocketbase.NewClient("http://localhost:8090")
	response, err := client.List("posts_public", pocketbase.ParamsList{
		Page: 1, Size: 10, Sort: "-created", Filters: "field~'test'",
	})
	if err != nil {
		log.Fatal(err)
	}
	log.Print(response.TotalItems)
}
```
Creating an item with admin user (auth via email/pass). 
Please note that you can pass `map[string]any` or `struct with JSON tags` as a payload:

```go
package main

import (
	"log"

	"github.com/r--w/pocketbase"
)

func main() {
	client := pocketbase.NewClient("http://localhost:8090", 
		pocketbase.WithAdminEmailPassword("admin@admin.com", "admin@admin.com"))
	response, err := client.Create("posts_admin", map[string]any{
		"field": "test",
	})
	if err != nil {
		log.Fatal(err)
	}
	log.Print(response.ID)
}
```
For even easier interaction with collection results as user-defined types, you can go with `CollectionSet`:

```go
package main

import (
	"log"

	"github.com/r--w/pocketbase"
)

type post struct {
	ID      string
	Field   string
	Created string
}

func main() {
	client := pocketbase.NewClient("http://localhost:8090")
	collection := pocketbase.CollectionSet[post](client, "posts_public")
	response, err := collection.List(pocketbase.ParamsList{
		Page: 1, Size: 10, Sort: "-created", Filters: "field~'test'",
	})
	if err != nil {
		log.Fatal(err)
	}
	
    log.Printf("%+v", response.Items)
}
```

Realtime API via Server-Sent Events (SSE) is also supported:

```go
package main

import (
	"log"

	"github.com/r--w/pocketbase"
)

type post struct {
	ID      string
	Field   string
	Created string
}

func main() {
	client := pocketbase.NewClient("http://localhost:8090")
	collection := pocketbase.CollectionSet[post](client, "posts_public")
	response, err := collection.List(pocketbase.ParamsList{
		Page: 1, Size: 10, Sort: "-created", Filters: "field~'test'",
	})
	if err != nil {
		log.Fatal(err)
	}
	
	stream, err := collection.Subscribe()
	if err != nil {
		log.Fatal(err)
	}
	defer stream.Unsubscribe()
	if err = stream.WaitAuthReady(); err != nil {
		log.Fatal(err)
	}
	for ev := range stream.Events() {
		log.Print(ev.Action, ev.Record)
	}
}
```

More examples can be found in:
* [example file](./example/main.go)
* [tests for the client](./client_test.go)
* [tests for the collection](./collection_test.go)
* remember to start the Pocketbase before running examples with `make serve` command

## Development

### Makefile targets 
* `make serve` - builds all binaries and runs local PocketBase server, it will create collections and sample data based on [migration files](./migrations)
* `make test` - runs tests (make sure that PocketBase server is running - `make serve` before)
* `make check` - runs linters and security checks (run this before commit)
* `make build` - builds all binaries (examples and PocketBase server) 
* `make help` - shows help and other targets

## Contributing
* Go 1.19+ (for making changes in the Go code)
* While developing use `WithDebug()` client option to see HTTP requests and responses
* Make sure that all checks are green (run `make check` before commit)
* Make sure that all tests pass (run `make test` before commit)
* Create a PR with your changes and wait for review