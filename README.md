# auto-api-client-go

[![Go Reference](https://pkg.go.dev/badge/github.com/autoapicom/auto-api-go.svg)](https://pkg.go.dev/github.com/autoapicom/auto-api-go)
[![Go Version](https://img.shields.io/github/go-mod/go-version/autoapicom/auto-api-go)](go.mod)
[![License](https://img.shields.io/github/license/autoapicom/auto-api-go)](LICENSE)

Go client for [auto-api.com](https://auto-api.com) — car listings API across multiple marketplaces.

One API to access car listings from 8 marketplaces: encar, mobile.de, autoscout24, che168, dongchedi, guazi, dubicars, dubizzle. Search offers, track price changes, and get listing data in a unified format.

## Installation

```bash
go get github.com/autoapicom/auto-api-go
```

## Usage

```go
import "github.com/autoapicom/auto-api-go"

client := autoapi.NewClient("your-api-key")
ctx := context.Background()
```

### Get filters

```go
filters, err := client.GetFilters(ctx, "encar")
```

### Search offers

```go
offers, err := client.GetOffers(ctx, "mobilede", &autoapi.OffersParams{
    Page:     1,
    Brand:    "BMW",
    YearFrom: 2020,
})

// Pagination
fmt.Println(offers.Meta.Page)
fmt.Println(offers.Meta.NextPage)
```

### Get single offer

```go
offer, err := client.GetOffer(ctx, "encar", "40427050")
```

### Track changes

```go
changeID, err := client.GetChangeID(ctx, "encar", "2025-01-15")
changes, err := client.GetChanges(ctx, "encar", changeID)

// Next batch
nextBatch, err := client.GetChanges(ctx, "encar", changes.Meta.NextChangeID)
```

### Get offer by URL

```go
info, err := client.GetOfferByURL(ctx, "https://www.encar.com/dc/dc_cardetailview.do?carid=40427050")
```

### Decode offer data

Offer data varies between sources, so it's stored as `json.RawMessage`. Decode into `OfferData` or your own struct:

```go
for _, item := range offers.Result {
    var d autoapi.OfferData
    json.Unmarshal(item.Data, &d)
    fmt.Printf("%s %s %s — $%s\n", d.Mark, d.Model, d.Year, d.Price)
}
```

### Error handling

```go
import "errors"

offers, err := client.GetOffers(ctx, "encar", &autoapi.OffersParams{Page: 1})
if err != nil {
    var authErr *autoapi.AuthError
    var apiErr *autoapi.ApiError
    if errors.As(err, &authErr) {
        // 401/403 — invalid API key
    } else if errors.As(err, &apiErr) {
        fmt.Println(apiErr.StatusCode, apiErr.Message)
    }
}
```

## Supported sources

| Source | Platform | Region |
|--------|----------|--------|
| `encar` | [encar.com](https://www.encar.com) | South Korea |
| `mobilede` | [mobile.de](https://www.mobile.de) | Germany |
| `autoscout24` | [autoscout24.com](https://www.autoscout24.com) | Europe |
| `che168` | [che168.com](https://www.che168.com) | China |
| `dongchedi` | [dongchedi.com](https://www.dongchedi.com) | China |
| `guazi` | [guazi.com](https://www.guazi.com) | China |
| `dubicars` | [dubicars.com](https://www.dubicars.com) | UAE |
| `dubizzle` | [dubizzle.com](https://www.dubizzle.com) | UAE |

## Other languages

| Language | Package |
|----------|---------|
| PHP | [auto-api/client](https://github.com/autoapicom/auto-api-php) |
| TypeScript | [@auto-api/client](https://github.com/autoapicom/auto-api-node) |
| Python | [auto-api-client](https://github.com/autoapicom/auto-api-python) |
| C# | [AutoApi.Client](https://github.com/autoapicom/auto-api-dotnet) |
| Java | [auto-api-client](https://github.com/autoapicom/auto-api-java) |
| Ruby | [auto-api-client](https://github.com/autoapicom/auto-api-ruby) |
| Rust | [auto-api-client](https://github.com/autoapicom/auto-api-rust) |

## Documentation

[auto-api.com](https://auto-api.com)
