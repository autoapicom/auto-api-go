# Claude Instructions — auto-api-go

## Language

All code comments, documentation, and commit messages must be in **English**.

## Commands

- Build: `go build ./...`
- Run tests: `go test ./...`

## Key Files

- `client.go` — Client struct with 6 public methods
- `types.go` — all request/response types
- `errors.go` — ApiError and AuthError types
- `go.mod` — module configuration

## Code Style

- Go 1.21+ required
- Standard library only — zero external dependencies
- `context.Context` as the first parameter of every public method
- Errors via return values `(value, error)`, use `errors.As()` to check type
- PascalCase for exported names
- `OfferItem.Data` is `json.RawMessage` — varies between sources
- `api_key` passed in query string for GET requests, `x-api-key` header for POST
- Keep the codebase simple — avoid unnecessary abstractions
- English godoc comments only
