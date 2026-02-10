package autoapi

import (
	"encoding/json"
	"fmt"
)

// ApiError represents a non-auth API error response.
type ApiError struct {
	StatusCode int
	Message    string
	Body       json.RawMessage
}

func (e *ApiError) Error() string {
	return fmt.Sprintf("API error %d: %s", e.StatusCode, e.Message)
}

// AuthError represents an authentication error (401/403).
type AuthError struct {
	StatusCode int
	Message    string
}

func (e *AuthError) Error() string {
	return fmt.Sprintf("auth error %d: %s", e.StatusCode, e.Message)
}
