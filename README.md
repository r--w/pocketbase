[![Check & test & build](https://github.com/r--w/pocketbase/actions/workflows/main.yml/badge.svg)](https://github.com/r--w/pocketbase/actions/workflows/main.yml)

## Development

### Makefile targets 
* `make serve` - builds all binaries and runs local PocketBase server
* `make test` - runs tests (make sure that PocketBase server is running - `make serve` before)
* `make check` - runs linters and security checks (run this before commit)
* `make build` - builds all binaries (examples and PocketBase server) 
* `make help` - shows help and other targets

## Contributing
* Go 1.19+ (for making changes in the Go code)
* Make sure that all checks are green (run `make check` before commit)
* Make sure that all tests pass (run `make test` before commit)
* Create a PR with your changes and wait for review

## Usage:
```go
  client := pocketbase.NewClient("http://localhost:8090", "admin@admin.com", "admin@admin.com")
  respBytes, _ := client.List("news", pocketbase.Params{Size: 2, Filters: "title~'Bitcoin'"})
  var respParsed response
  json.Unmarshal(respBytes, &respParsed)
  log.Print(respParsed)
```