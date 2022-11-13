![GithubActions](https://github.com/coinpaprika/news/workflows/Check%20&%20test%20&%20build/badge.svg)

## Go
At least version 1.19 is required

## Building
`make build`

## Run tests:
`make test`

## Quality & linteners:
`make check`

## Other targets
`make help`

## Usage:
```go
  client := pocketbase.NewClient("http://localhost:8090", "admin@admin.com", "admin@admin.com")
  respBytes, _ := client.List("news", pocketbase.Params{Size: 2, Filters: "title~'Bitcoin'"})
  var respParsed response
  json.Unmarshal(respBytes, &respParsed)
  log.Print(respParsed)
```