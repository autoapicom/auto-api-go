# Copilot Instructions — auto-api-go

This is a Go 1.21+ client library for the auto-api.com vehicle data API.

## Architecture

- Single `Client` struct (`client.go`) with 6 public methods
- Every method takes `context.Context` as its first parameter and `source` as second
- Standard library only for HTTP (`net/http`, `encoding/json`) — zero dependencies
- Methods return `(value, error)` pairs
- `OfferItem.Data` is `json.RawMessage` — varies between sources
- Authentication: `api_key` in query string for GET, `x-api-key` header for POST

## Error Types

- `ApiError` — base error type for all API errors, implements `Error()` method
- `AuthError` — returned on 401/403 responses, implements `Error()` method
- Use `errors.As()` to check error types

## Conventions

- PascalCase for all exported names
- All code comments and documentation must be in English
- Keep it simple — no unnecessary abstractions or wrapper types
