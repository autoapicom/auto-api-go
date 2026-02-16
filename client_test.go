package autoapi

import (
	"context"
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
)

// newTestClient creates a client pointed at the given test server.
func newTestClient(serverURL string) *Client {
	c := NewClient("test-key")
	c.SetBaseURL(serverURL)
	return c
}

// jsonHandler returns an http.HandlerFunc that responds with the given status and body.
func jsonHandler(status int, body interface{}) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(status)
		json.NewEncoder(w).Encode(body)
	}
}

// ── GetFilters ──────────────────────────────────────────────────

func TestGetFilters(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte(`{"brands":["Toyota","Honda"],"body_types":["sedan","suv"]}`))
	}))
	defer srv.Close()

	client := newTestClient(srv.URL)
	result, err := client.GetFilters(context.Background(), "encar")

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	brands, ok := result["brands"].([]interface{})
	if !ok || len(brands) != 2 {
		t.Fatalf("expected 2 brands, got %v", result["brands"])
	}
}

func TestGetFiltersEndpoint(t *testing.T) {
	var capturedPath string
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		capturedPath = r.URL.Path
		w.Write([]byte(`{}`))
	}))
	defer srv.Close()

	client := newTestClient(srv.URL)
	client.GetFilters(context.Background(), "encar")

	if capturedPath != "/api/v2/encar/filters" {
		t.Errorf("expected path /api/v2/encar/filters, got %s", capturedPath)
	}
}

func TestGetFiltersApiKey(t *testing.T) {
	var capturedQuery string
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		capturedQuery = r.URL.RawQuery
		w.Write([]byte(`{}`))
	}))
	defer srv.Close()

	client := newTestClient(srv.URL)
	client.GetFilters(context.Background(), "encar")

	if q := r(capturedQuery, "api_key=test-key"); !q {
		t.Errorf("expected api_key=test-key in query, got %s", capturedQuery)
	}
}

func TestGetFiltersUsesGET(t *testing.T) {
	var capturedMethod string
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		capturedMethod = r.Method
		w.Write([]byte(`{}`))
	}))
	defer srv.Close()

	client := newTestClient(srv.URL)
	client.GetFilters(context.Background(), "encar")

	if capturedMethod != "GET" {
		t.Errorf("expected GET, got %s", capturedMethod)
	}
}

// ── GetOffers ───────────────────────────────────────────────────

func TestGetOffers(t *testing.T) {
	srv := httptest.NewServer(jsonHandler(200, map[string]interface{}{
		"result": []map[string]interface{}{{"inner_id": "a1"}},
		"meta":   map[string]interface{}{"page": 1, "next_page": 2, "limit": 20},
	}))
	defer srv.Close()

	client := newTestClient(srv.URL)
	result, err := client.GetOffers(context.Background(), "encar", &OffersParams{Page: 1})

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(result.Result) != 1 {
		t.Errorf("expected 1 result, got %d", len(result.Result))
	}
	if result.Meta.Page != 1 {
		t.Errorf("expected page 1, got %d", result.Meta.Page)
	}
}

func TestGetOffersWithFilters(t *testing.T) {
	var capturedQuery string
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		capturedQuery = r.URL.RawQuery
		w.Write([]byte(`{"result":[],"meta":{"page":2,"next_page":3,"limit":20}}`))
	}))
	defer srv.Close()

	client := newTestClient(srv.URL)
	client.GetOffers(context.Background(), "mobile_de", &OffersParams{
		Page:  2,
		Brand: "BMW",
	})

	if !contains(capturedQuery, "page=2") {
		t.Errorf("expected page=2 in query, got %s", capturedQuery)
	}
	if !contains(capturedQuery, "brand=BMW") {
		t.Errorf("expected brand=BMW in query, got %s", capturedQuery)
	}
}

func TestGetOffersNilParams(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte(`{"result":[],"meta":{"page":0,"next_page":0,"limit":20}}`))
	}))
	defer srv.Close()

	client := newTestClient(srv.URL)
	_, err := client.GetOffers(context.Background(), "encar", nil)

	if err != nil {
		t.Fatalf("unexpected error with nil params: %v", err)
	}
}

// ── GetOffer ────────────────────────────────────────────────────

func TestGetOffer(t *testing.T) {
	var capturedQuery string
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		capturedQuery = r.URL.RawQuery
		w.Write([]byte(`{"result":[{"inner_id":"abc123"}],"meta":{"page":1,"next_page":0,"limit":20}}`))
	}))
	defer srv.Close()

	client := newTestClient(srv.URL)
	result, err := client.GetOffer(context.Background(), "encar", "abc123")

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !contains(capturedQuery, "inner_id=abc123") {
		t.Errorf("expected inner_id=abc123 in query, got %s", capturedQuery)
	}
	if len(result.Result) != 1 {
		t.Errorf("expected 1 result, got %d", len(result.Result))
	}
}

// ── GetChangeID ─────────────────────────────────────────────────

func TestGetChangeID(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte(`{"change_id":42567}`))
	}))
	defer srv.Close()

	client := newTestClient(srv.URL)
	result, err := client.GetChangeID(context.Background(), "encar", "2024-01-15")

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if result != 42567 {
		t.Errorf("expected 42567, got %d", result)
	}
}

func TestGetChangeIDPassesDate(t *testing.T) {
	var capturedQuery string
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		capturedQuery = r.URL.RawQuery
		w.Write([]byte(`{"change_id":0}`))
	}))
	defer srv.Close()

	client := newTestClient(srv.URL)
	client.GetChangeID(context.Background(), "encar", "2024-01-15")

	if !contains(capturedQuery, "date=2024-01-15") {
		t.Errorf("expected date=2024-01-15 in query, got %s", capturedQuery)
	}
}

func TestGetChangeIDZero(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte(`{"change_id":0}`))
	}))
	defer srv.Close()

	client := newTestClient(srv.URL)
	result, err := client.GetChangeID(context.Background(), "encar", "2024-01-01")

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if result != 0 {
		t.Errorf("expected 0, got %d", result)
	}
}

// ── GetChanges ──────────────────────────────────────────────────

func TestGetChanges(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte(`{"result":[{"id":1,"inner_id":"new1","change_type":"added","created_at":"2024-01-15","data":{}}],"meta":{"cur_change_id":42567,"next_change_id":42568,"limit":500}}`))
	}))
	defer srv.Close()

	client := newTestClient(srv.URL)
	result, err := client.GetChanges(context.Background(), "encar", 42567)

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(result.Result) != 1 {
		t.Errorf("expected 1 result, got %d", len(result.Result))
	}
	if result.Meta.CurChangeID != 42567 {
		t.Errorf("expected cur_change_id 42567, got %d", result.Meta.CurChangeID)
	}
}

func TestGetChangesPassesChangeID(t *testing.T) {
	var capturedQuery string
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		capturedQuery = r.URL.RawQuery
		w.Write([]byte(`{"result":[],"meta":{"cur_change_id":0,"next_change_id":0,"limit":500}}`))
	}))
	defer srv.Close()

	client := newTestClient(srv.URL)
	client.GetChanges(context.Background(), "encar", 42567)

	if !contains(capturedQuery, "change_id=42567") {
		t.Errorf("expected change_id=42567 in query, got %s", capturedQuery)
	}
}

// ── GetOfferByURL ───────────────────────────────────────────────

func TestGetOfferByURL(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte(`{"brand":"BMW","model":"X5","price":45000}`))
	}))
	defer srv.Close()

	client := newTestClient(srv.URL)
	result, err := client.GetOfferByURL(context.Background(), "https://www.encar.com/car/123")

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if result["brand"] != "BMW" {
		t.Errorf("expected brand BMW, got %v", result["brand"])
	}
}

func TestGetOfferByURLUsesPOST(t *testing.T) {
	var capturedMethod string
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		capturedMethod = r.Method
		w.Write([]byte(`{}`))
	}))
	defer srv.Close()

	client := newTestClient(srv.URL)
	client.GetOfferByURL(context.Background(), "https://example.com/car/123")

	if capturedMethod != "POST" {
		t.Errorf("expected POST, got %s", capturedMethod)
	}
}

func TestGetOfferByURLUsesV1Endpoint(t *testing.T) {
	var capturedPath string
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		capturedPath = r.URL.Path
		w.Write([]byte(`{}`))
	}))
	defer srv.Close()

	client := newTestClient(srv.URL)
	client.GetOfferByURL(context.Background(), "https://example.com/car/123")

	if capturedPath != "/api/v1/offer/info" {
		t.Errorf("expected /api/v1/offer/info, got %s", capturedPath)
	}
}

func TestGetOfferByURLSendsApiKeyInHeader(t *testing.T) {
	var capturedHeader string
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		capturedHeader = r.Header.Get("x-api-key")
		w.Write([]byte(`{}`))
	}))
	defer srv.Close()

	client := newTestClient(srv.URL)
	client.GetOfferByURL(context.Background(), "https://example.com/car/123")

	if capturedHeader != "test-key" {
		t.Errorf("expected x-api-key header test-key, got %s", capturedHeader)
	}
}

func TestGetOfferByURLSendsURLInBody(t *testing.T) {
	var capturedBody map[string]string
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		body, _ := io.ReadAll(r.Body)
		json.Unmarshal(body, &capturedBody)
		w.Write([]byte(`{}`))
	}))
	defer srv.Close()

	client := newTestClient(srv.URL)
	client.GetOfferByURL(context.Background(), "https://example.com/car/123")

	if capturedBody["url"] != "https://example.com/car/123" {
		t.Errorf("expected url in body, got %v", capturedBody)
	}
}

func TestGetOfferByURLSendsContentType(t *testing.T) {
	var capturedCT string
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		capturedCT = r.Header.Get("Content-Type")
		w.Write([]byte(`{}`))
	}))
	defer srv.Close()

	client := newTestClient(srv.URL)
	client.GetOfferByURL(context.Background(), "https://example.com/car/123")

	if capturedCT != "application/json" {
		t.Errorf("expected Content-Type application/json, got %s", capturedCT)
	}
}

// ── Custom API version ──────────────────────────────────────────

func TestCustomAPIVersion(t *testing.T) {
	var capturedPath string
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		capturedPath = r.URL.Path
		w.Write([]byte(`{}`))
	}))
	defer srv.Close()

	client := newTestClient(srv.URL)
	client.SetAPIVersion("v3")
	client.GetFilters(context.Background(), "encar")

	if capturedPath != "/api/v3/encar/filters" {
		t.Errorf("expected /api/v3/encar/filters, got %s", capturedPath)
	}
}

// ── Error handling ──────────────────────────────────────────────

func TestApiErrorOnServerError(t *testing.T) {
	srv := httptest.NewServer(jsonHandler(500, map[string]string{"message": "Internal server error"}))
	defer srv.Close()

	client := newTestClient(srv.URL)
	_, err := client.GetFilters(context.Background(), "encar")

	if err == nil {
		t.Fatal("expected error")
	}

	var apiErr *ApiError
	if !errors.As(err, &apiErr) {
		t.Fatalf("expected *ApiError, got %T", err)
	}
	if apiErr.StatusCode != 500 {
		t.Errorf("expected status 500, got %d", apiErr.StatusCode)
	}
}

func TestApiErrorContainsBody(t *testing.T) {
	srv := httptest.NewServer(jsonHandler(422, map[string]interface{}{"message": "Validation failed", "errors": []string{"invalid page"}}))
	defer srv.Close()

	client := newTestClient(srv.URL)
	_, err := client.GetFilters(context.Background(), "encar")

	var apiErr *ApiError
	if !errors.As(err, &apiErr) {
		t.Fatalf("expected *ApiError, got %T", err)
	}
	if apiErr.Body == nil {
		t.Error("expected non-nil body")
	}
}

func TestApiErrorUsesMessageFromBody(t *testing.T) {
	srv := httptest.NewServer(jsonHandler(500, map[string]string{"message": "Custom error message"}))
	defer srv.Close()

	client := newTestClient(srv.URL)
	_, err := client.GetFilters(context.Background(), "encar")

	var apiErr *ApiError
	if !errors.As(err, &apiErr) {
		t.Fatalf("expected *ApiError, got %T", err)
	}
	if apiErr.Message != "Custom error message" {
		t.Errorf("expected 'Custom error message', got '%s'", apiErr.Message)
	}
}

func TestApiErrorFallbackMessage(t *testing.T) {
	srv := httptest.NewServer(jsonHandler(500, map[string]string{"error": "something"}))
	defer srv.Close()

	client := newTestClient(srv.URL)
	_, err := client.GetFilters(context.Background(), "encar")

	var apiErr *ApiError
	if !errors.As(err, &apiErr) {
		t.Fatalf("expected *ApiError, got %T", err)
	}
	if apiErr.Message != "API error: 500" {
		t.Errorf("expected fallback message, got '%s'", apiErr.Message)
	}
}

func TestAuthErrorOn401(t *testing.T) {
	srv := httptest.NewServer(jsonHandler(401, map[string]string{"message": "Unauthorized"}))
	defer srv.Close()

	client := newTestClient(srv.URL)
	_, err := client.GetFilters(context.Background(), "encar")

	var authErr *AuthError
	if !errors.As(err, &authErr) {
		t.Fatalf("expected *AuthError, got %T", err)
	}
	if authErr.StatusCode != 401 {
		t.Errorf("expected status 401, got %d", authErr.StatusCode)
	}
}

func TestAuthErrorOn403(t *testing.T) {
	srv := httptest.NewServer(jsonHandler(403, map[string]string{"message": "Forbidden"}))
	defer srv.Close()

	client := newTestClient(srv.URL)
	_, err := client.GetOffers(context.Background(), "encar", nil)

	var authErr *AuthError
	if !errors.As(err, &authErr) {
		t.Fatalf("expected *AuthError, got %T", err)
	}
}

func TestApiErrorOnInvalidJSON(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("not json at all"))
	}))
	defer srv.Close()

	client := newTestClient(srv.URL)
	_, err := client.GetFilters(context.Background(), "encar")

	if err == nil {
		t.Fatal("expected error for invalid JSON")
	}

	var apiErr *ApiError
	if !errors.As(err, &apiErr) {
		t.Fatalf("expected *ApiError, got %T", err)
	}
	if !contains(apiErr.Message, "invalid JSON response") {
		t.Errorf("expected 'invalid JSON response' in message, got '%s'", apiErr.Message)
	}
}

func TestNotFoundReturnsApiErrorNotAuthError(t *testing.T) {
	srv := httptest.NewServer(jsonHandler(404, map[string]string{"message": "Source not found"}))
	defer srv.Close()

	client := newTestClient(srv.URL)
	_, err := client.GetFilters(context.Background(), "unknown_source")

	var authErr *AuthError
	if errors.As(err, &authErr) {
		t.Error("expected ApiError, not AuthError")
	}

	var apiErr *ApiError
	if !errors.As(err, &apiErr) {
		t.Fatalf("expected *ApiError, got %T", err)
	}
	if apiErr.StatusCode != 404 {
		t.Errorf("expected status 404, got %d", apiErr.StatusCode)
	}
}

// ── encodeParams ────────────────────────────────────────────────

func TestEncodeParamsOmitsEmpty(t *testing.T) {
	client := NewClient("test-key")
	params := &OffersParams{Page: 1, Brand: "BMW"}
	query := client.encodeParams(params)

	if query.Get("page") != "1" {
		t.Errorf("expected page=1, got %s", query.Get("page"))
	}
	if query.Get("brand") != "BMW" {
		t.Errorf("expected brand=BMW, got %s", query.Get("brand"))
	}
	if query.Get("model") != "" {
		t.Errorf("expected model to be omitted, got %s", query.Get("model"))
	}
	if query.Get("year_from") != "" {
		t.Errorf("expected year_from to be omitted, got %s", query.Get("year_from"))
	}
}

func TestEncodeParamsNil(t *testing.T) {
	client := NewClient("test-key")
	query := client.encodeParams(nil)

	if len(query) != 0 {
		t.Errorf("expected empty query for nil params, got %v", query)
	}
}

// ── Helpers ─────────────────────────────────────────────────────

// r is an alias for contains used in one test above.
func r(s, substr string) bool {
	return contains(s, substr)
}

func contains(s, substr string) bool {
	return len(s) >= len(substr) && searchString(s, substr)
}

func searchString(s, substr string) bool {
	for i := 0; i <= len(s)-len(substr); i++ {
		if s[i:i+len(substr)] == substr {
			return true
		}
	}
	return false
}
