[![Check & test & build](https://github.com/r--w/pocketbase/actions/workflows/main.yml/badge.svg)](https://github.com/r--w/pocketbase/actions/workflows/main.yml)

## Development
## Go
At least version 1.19 is required

### Makefile targets 
* `make serve` - build all binaries and run local PocketBase server
* `make test` - run tests (make sure that PocketBase server is running - `make serve` before)
* `make check` - run linters and security checks (run this before commit)
* `make build` - build all binaries (examples and PocketBase server) 
* `make help` - show help and other targets

## Contributing
* Go 1.19+ (for making changes in the Go code)
* Make sure that all checks are green (locally run `make check`)
* Make sure that all tests pass (locally run `make test`)
* Create a PR with your changes and wait for review

## Usage:
```go
  client := pocketbase.NewClient("http://localhost:8090", "admin@admin.com", "admin@admin.com")
  respBytes, _ := client.List("news", pocketbase.Params{Size: 2, Filters: "title~'Bitcoin'"})
  var respParsed response
  json.Unmarshal(respBytes, &respParsed)
  log.Print(respParsed)
```