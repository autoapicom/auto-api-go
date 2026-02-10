package autoapi

import "encoding/json"

// OffersParams holds query parameters for GetOffers.
type OffersParams struct {
	Page          int    `url:"page"`
	Brand         string `url:"brand,omitempty"`
	Model         string `url:"model,omitempty"`
	Configuration string `url:"configuration,omitempty"`
	Complectation string `url:"complectation,omitempty"`
	Transmission  string `url:"transmission,omitempty"`
	Color         string `url:"color,omitempty"`
	BodyType      string `url:"body_type,omitempty"`
	EngineType    string `url:"engine_type,omitempty"`
	YearFrom      int    `url:"year_from,omitempty"`
	YearTo        int    `url:"year_to,omitempty"`
	MileageFrom   int    `url:"mileage_from,omitempty"`
	MileageTo     int    `url:"mileage_to,omitempty"`
	PriceFrom     int    `url:"price_from,omitempty"`
	PriceTo       int    `url:"price_to,omitempty"`
}

// OffersResponse is the response from GetOffers.
type OffersResponse struct {
	Result []OfferItem `json:"result"`
	Meta   Meta        `json:"meta"`
}

// ChangesResponse is the response from GetChanges.
type ChangesResponse struct {
	Result []ChangeItem `json:"result"`
	Meta   ChangesMeta  `json:"meta"`
}

// Meta holds pagination info for offers.
type Meta struct {
	Page     int `json:"page"`
	NextPage int `json:"next_page"`
	Limit    int `json:"limit"`
}

// ChangesMeta holds pagination info for changes.
type ChangesMeta struct {
	CurChangeID  int `json:"cur_change_id"`
	NextChangeID int `json:"next_change_id"`
	Limit        int `json:"limit"`
}

// OfferItem is a single item in the offers result array.
type OfferItem struct {
	ID         int             `json:"id"`
	InnerID    string          `json:"inner_id"`
	ChangeType string          `json:"change_type"`
	CreatedAt  string          `json:"created_at"`
	Data       json.RawMessage `json:"data"`
}

// ChangeItem is a single item in the changes result array.
type ChangeItem struct {
	ID         int             `json:"id"`
	InnerID    string          `json:"inner_id"`
	ChangeType string          `json:"change_type"`
	CreatedAt  string          `json:"created_at"`
	Data       json.RawMessage `json:"data"`
}

// ChangeIDResponse is the response from GetChangeID.
type ChangeIDResponse struct {
	ChangeID int `json:"change_id"`
}

// OfferData holds common offer fields shared across all sources.
// Since each source may have additional fields, use json.RawMessage
// from OfferItem.Data and unmarshal into your own struct if needed.
type OfferData struct {
	InnerID          string   `json:"inner_id"`
	URL              string   `json:"url"`
	Mark             string   `json:"mark"`
	Model            string   `json:"model"`
	Generation       string   `json:"generation"`
	Configuration    string   `json:"configuration"`
	Complectation    string   `json:"complectation"`
	Year             string   `json:"year"`
	Color            string   `json:"color"`
	Price            string   `json:"price"`
	KmAge            string   `json:"km_age"`
	EngineType       string   `json:"engine_type"`
	TransmissionType string   `json:"transmission_type"`
	BodyType         string   `json:"body_type"`
	Address          string   `json:"address"`
	SellerType       string   `json:"seller_type"`
	IsDealer         bool     `json:"is_dealer"`
	Displacement     string   `json:"displacement"`
	OfferCreated     string   `json:"offer_created"`
	Images           []string `json:"images"`
}
