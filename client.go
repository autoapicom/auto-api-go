// Package autoapi provides a client for the auto-api.com car listings API.
package autoapi

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"reflect"
	"strconv"
	"strings"
	"time"
)

// Client is the auto-api.com API client.
type Client struct {
	apiKey     string
	baseURL    string
	apiVersion string
	httpClient *http.Client
}

// NewClient creates a new API client with the given API key.
func NewClient(apiKey string) *Client {
	return &Client{
		apiKey:     apiKey,
		baseURL:    "https://auto-api.com",
		apiVersion: "v2",
		httpClient: &http.Client{Timeout: 30 * time.Second},
	}
}

// SetBaseURL overrides the default base URL.
func (c *Client) SetBaseURL(baseURL string) {
	c.baseURL = strings.TrimRight(baseURL, "/")
}

// SetAPIVersion overrides the default API version (default: "v2").
func (c *Client) SetAPIVersion(version string) {
	c.apiVersion = version
}

// GetFilters returns available filters for a source (brands, models, body types, etc.)
func (c *Client) GetFilters(ctx context.Context, source string) (map[string]interface{}, error) {
	var result map[string]interface{}
	err := c.get(ctx, fmt.Sprintf("api/%s/%s/filters", c.apiVersion, source), nil, &result)
	return result, err
}

// GetOffers returns a paginated list of offers with optional filters.
func (c *Client) GetOffers(ctx context.Context, source string, params *OffersParams) (*OffersResponse, error) {
	query := c.encodeParams(params)
	var result OffersResponse
	err := c.get(ctx, fmt.Sprintf("api/%s/%s/offers", c.apiVersion, source), query, &result)
	return &result, err
}

// GetOffer returns a single offer by inner_id.
func (c *Client) GetOffer(ctx context.Context, source string, innerID string) (*OffersResponse, error) {
	query := url.Values{"inner_id": {innerID}}
	var result OffersResponse
	err := c.get(ctx, fmt.Sprintf("api/%s/%s/offer", c.apiVersion, source), query, &result)
	return &result, err
}

// GetChangeID returns a change_id for the given date (format: yyyy-mm-dd).
func (c *Client) GetChangeID(ctx context.Context, source string, date string) (int, error) {
	query := url.Values{"date": {date}}
	var result ChangeIDResponse
	err := c.get(ctx, fmt.Sprintf("api/%s/%s/change_id", c.apiVersion, source), query, &result)
	return result.ChangeID, err
}

// GetChanges returns a changes feed (added/changed/removed) starting from change_id.
func (c *Client) GetChanges(ctx context.Context, source string, changeID int) (*ChangesResponse, error) {
	query := url.Values{"change_id": {strconv.Itoa(changeID)}}
	var result ChangesResponse
	err := c.get(ctx, fmt.Sprintf("api/%s/%s/changes", c.apiVersion, source), query, &result)
	return &result, err
}

// GetOfferByURL returns offer data by its URL on the marketplace.
func (c *Client) GetOfferByURL(ctx context.Context, offerURL string) (map[string]interface{}, error) {
	var result map[string]interface{}
	err := c.post(ctx, "api/v1/offer/info", map[string]string{"url": offerURL}, &result)
	return result, err
}

func (c *Client) get(ctx context.Context, endpoint string, query url.Values, target interface{}) error {
	if query == nil {
		query = url.Values{}
	}
	query.Set("api_key", c.apiKey)

	reqURL := fmt.Sprintf("%s/%s?%s", c.baseURL, endpoint, query.Encode())
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, reqURL, nil)
	if err != nil {
		return err
	}

	return c.do(req, target)
}

func (c *Client) post(ctx context.Context, endpoint string, data interface{}, target interface{}) error {
	body, err := json.Marshal(data)
	if err != nil {
		return err
	}

	reqURL := fmt.Sprintf("%s/%s", c.baseURL, endpoint)
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, reqURL, bytes.NewReader(body))
	if err != nil {
		return err
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("x-api-key", c.apiKey)

	return c.do(req, target)
}

func (c *Client) do(req *http.Request, target interface{}) error {
	resp, err := c.httpClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return c.handleError(resp.StatusCode, respBody)
	}

	if err := json.Unmarshal(respBody, target); err != nil {
		return &ApiError{
			StatusCode: resp.StatusCode,
			Message:    fmt.Sprintf("invalid JSON response: %s", string(respBody[:min(len(respBody), 200)])),
		}
	}

	return nil
}

func (c *Client) handleError(statusCode int, body []byte) error {
	message := fmt.Sprintf("API error: %d", statusCode)

	var parsed map[string]interface{}
	if json.Unmarshal(body, &parsed) == nil {
		if msg, ok := parsed["message"].(string); ok {
			message = msg
		}
	}

	if statusCode == 401 || statusCode == 403 {
		return &AuthError{StatusCode: statusCode, Message: message}
	}

	return &ApiError{StatusCode: statusCode, Message: message, Body: body}
}

func (c *Client) encodeParams(params *OffersParams) url.Values {
	query := url.Values{}
	if params == nil {
		return query
	}

	v := reflect.ValueOf(*params)
	t := v.Type()

	for i := 0; i < v.NumField(); i++ {
		field := v.Field(i)
		tag := t.Field(i).Tag.Get("url")
		if tag == "" {
			continue
		}

		parts := strings.SplitN(tag, ",", 2)
		name := parts[0]
		omitempty := len(parts) > 1 && parts[1] == "omitempty"

		switch field.Kind() {
		case reflect.String:
			if s := field.String(); s != "" || !omitempty {
				query.Set(name, s)
			}
		case reflect.Int:
			if n := field.Int(); n != 0 || !omitempty {
				query.Set(name, strconv.FormatInt(n, 10))
			}
		}
	}

	return query
}
