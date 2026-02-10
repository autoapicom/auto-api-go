# auto-api Go Client

Go client for [auto-api.com](https://auto-api.com) — car listings API across 8 marketplaces.

## Quick Start

```bash
go get github.com/autoapicom/auto-api-go
```

```go
import "github.com/autoapicom/auto-api-go"

client := autoapi.NewClient("your-api-key")
offers, err := client.GetOffers(ctx, "encar", &autoapi.OffersParams{Page: 1})
```

## Build & Test

```bash
go build ./...
go test ./...
```

## Key Files

- `client.go` — Client struct, NewClient(), 6 public methods, HTTP helpers
- `types.go` — All types: OffersParams, OffersResponse, Meta, OfferItem, OfferData, etc.
- `errors.go` — ApiError and AuthError types
- `go.mod` — Module: github.com/autoapicom/auto-api-go, Go 1.21+

## Conventions

- Go 1.21+, stdlib only (net/http, encoding/json), zero dependencies
- context.Context as first parameter of every public method
- Errors via return (value, error), use errors.As() to check type
- PascalCase for exported names (Go convention)
- OfferItem.Data is json.RawMessage — structure varies between sources
- All comments are godoc-compatible, written in English

## API Methods

| Method | Params | Returns |
|--------|--------|---------|
| `GetFilters(ctx, source)` | source name | `map[string]interface{}` |
| `GetOffers(ctx, source, params)` | source + *OffersParams | `*OffersResponse` |
| `GetOffer(ctx, source, innerID)` | source + inner_id | `*OffersResponse` |
| `GetChangeID(ctx, source, date)` | source + yyyy-mm-dd | `int` |
| `GetChanges(ctx, source, changeID)` | source + change_id | `*ChangesResponse` |
| `GetOfferByURL(ctx, url)` | marketplace URL | `map[string]interface{}` |
