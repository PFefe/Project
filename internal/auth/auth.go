package auth

import (
	"errors"
	"net/http"
	"strings"
)

// GetAPIKey extracts the API key from the request header
// Example:
// Authorization: ApiKey{insert-api-key-here}
func GetAPIKey(headers *http.Request) (string, error) {
	val := headers.Header.Get("Authorization")
	if val == "" {
		return "", errors.New("no API key provided")
	}
	vals := strings.Split(
		val,
		" ",
	)
	if len(vals) != 2 {
		return "", errors.New("invalid API key format")
	}
	if vals[0] != "ApiKey" {
		return "", errors.New("invalid API key type")
	}
	return vals[1], nil
}
