[![Check & test & build](https://github.com/r--w/pocketbase/actions/workflows/main.yml/badge.svg)](https://github.com/r--w/pocketbase/actions/workflows/main.yml)

### Project
This repository contains community-maintained Go SDK for Pocketbase API.
It's well-tested and used in production in [Coinpaprika](https://coinpaprika.com), but not all endpoints are covered yet.

Currently supported operations:
* **Authentication** - anonymous, admin and user via email/password)
* **Create** 
* **Update**
* **Delete**
* **List** - with pagination, filtering, sorting
* **Other** - are planned to be implemented in the future after reaching 50 stars, but feel free to clone or contribute.

### PocketBase
[Pocketbase](https://pocketbase.io) is a simple, self-hosted, open-source, no-code, database for your personal data. 
It's a great alternative to Airtable, Notion, and Google Sheets. Source code is available on [github.com/pocketbase/pocketbase](https://github.com/pocketbase/pocketbase)

### Usage:
```go
  client := pocketbase.NewClient("http://localhost:8090", "admin@admin.com", "admin@admin.com")
  respBytes, _ := client.List("news", pocketbase.Params{Size: 2, Filters: "title~'Bitcoin'"})
  var respParsed response
  json.Unmarshal(respBytes, &respParsed)
  log.Print(respParsed)
```

## Development

### Makefile targets 
* `make serve` - builds all binaries and runs local PocketBase server, it will create collections and sample data based on [migration files](./migrations)
* `make test` - runs tests (make sure that PocketBase server is running - `make serve` before)
* `make check` - runs linters and security checks (run this before commit)
* `make build` - builds all binaries (examples and PocketBase server) 
* `make help` - shows help and other targets

## Contributing
* Go 1.19+ (for making changes in the Go code)
* Make sure that all checks are green (run `make check` before commit)
* Make sure that all tests pass (run `make test` before commit)
* Create a PR with your changes and wait for review