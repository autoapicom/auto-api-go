// Auto API Go Client — Complete usage example.
//
// Replace "your-api-key" with your actual API key from https://auto-api.com
//
// Run: go run examples/main.go
package main

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log"

	"github.com/autoapicom/auto-api-go"
)

func main() {
	client := autoapi.NewClient("your-api-key")
	ctx := context.Background()
	source := "encar"

	// --- Get available filters ---

	filters, err := client.GetFilters(ctx, source)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("Filter keys: %d\n", len(filters))

	// --- Search offers with filters ---

	offers, err := client.GetOffers(ctx, source, &autoapi.OffersParams{
		Page:     1,
		Brand:    "Hyundai",
		YearFrom: 2020,
		PriceTo:  50000,
	})
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("\n--- Offers (page %d) ---\n", offers.Meta.Page)
	for _, item := range offers.Result {
		var d autoapi.OfferData
		json.Unmarshal(item.Data, &d)
		fmt.Printf("%s %s %s — $%s (%s km)\n", d.Mark, d.Model, d.Year, d.Price, d.KmAge)
	}

	// Pagination
	if offers.Meta.NextPage > 0 {
		nextPage, err := client.GetOffers(ctx, source, &autoapi.OffersParams{
			Page:     offers.Meta.NextPage,
			Brand:    "Hyundai",
			YearFrom: 2020,
		})
		if err == nil {
			fmt.Printf("Next page has %d offers\n", len(nextPage.Result))
		}
	}

	// --- Get single offer ---

	innerID := "40427050"
	if len(offers.Result) > 0 {
		innerID = offers.Result[0].InnerID
	}

	offer, err := client.GetOffer(ctx, source, innerID)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("\n--- Single offer ---\n")
	if len(offer.Result) > 0 {
		var d autoapi.OfferData
		json.Unmarshal(offer.Result[0].Data, &d)
		fmt.Printf("URL: %s\nSeller: %s\nImages: %d\n", d.URL, d.SellerType, len(d.Images))
	}

	// --- Track changes ---

	changeID, err := client.GetChangeID(ctx, source, "2025-01-15")
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("\n--- Changes from 2025-01-15 (change_id: %d) ---\n", changeID)

	changes, err := client.GetChanges(ctx, source, changeID)
	if err != nil {
		log.Fatal(err)
	}

	for _, change := range changes.Result {
		fmt.Printf("[%s] %s\n", change.ChangeType, change.InnerID)
	}

	if changes.Meta.NextChangeID > 0 {
		moreChanges, err := client.GetChanges(ctx, source, changes.Meta.NextChangeID)
		if err == nil {
			fmt.Printf("Next batch: %d changes\n", len(moreChanges.Result))
		}
	}

	// --- Get offer by URL ---

	info, err := client.GetOfferByURL(ctx, "https://www.encar.com/dc/dc_cardetailview.do?carid=40427050")
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("\n--- Offer by URL ---\n")
	fmt.Printf("%s %s — $%v\n", info["mark"], info["model"], info["price"])

	// --- Error handling ---

	badClient := autoapi.NewClient("invalid-key")
	_, err = badClient.GetOffers(ctx, "encar", &autoapi.OffersParams{Page: 1})
	if err != nil {
		var authErr *autoapi.AuthError
		var apiErr *autoapi.ApiError
		if errors.As(err, &authErr) {
			fmt.Printf("\nAuth error: %s (HTTP %d)\n", authErr.Message, authErr.StatusCode)
		} else if errors.As(err, &apiErr) {
			fmt.Printf("\nAPI error: %s (HTTP %d)\n", apiErr.Message, apiErr.StatusCode)
		} else {
			fmt.Printf("\nError: %s\n", err)
		}
	}
}
