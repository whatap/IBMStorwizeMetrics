package IBMStorwizeMetrics

import (
	"crypto/tls"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"strings"
	"sync"
	"time"
)

// AuthResponse represents the structure of the authentication response
type AuthResponse struct {
	Token   string `json:"token"`
	Expires int64  `json:"expires"` // Assuming UNIX timestamp for simplicity
}

// AuthCache is used to cache the authentication token
type AuthCache struct {
	sync.Mutex
	Token              string
	Expires            time.Time
	URL                string
	Username           string
	Password           string
	InsecureSkipVerify bool
}

// NewAuthCache creates a new AuthCache instance
func NewAuthCache(url, username, password string, insecureSkipVerify bool) *AuthCache {
	return &AuthCache{
		URL:                url,
		Username:           username,
		Password:           password,
		InsecureSkipVerify: insecureSkipVerify,
	}
}

// GetToken fetches or returns the cached authentication token
func (authCache *AuthCache) GetToken() (string, error) {
	authCache.Lock()
	defer authCache.Unlock()

	// If the token is still valid, return it
	if time.Now().Before(authCache.Expires) {
		return authCache.Token, nil
	}

	// Otherwise, fetch a new token
	token, expires, err := authCache.fetchAuthToken()
	if err != nil {
		return "", err
	}

	authCache.Token = token
	authCache.Expires = expires

	return authCache.Token, nil
}

func (authCache *AuthCache) fetchAuthToken() (string, time.Time, error) {
	// Concatenate /auth to the endpoint URL
	authURL := authCache.URL + "/auth"

	// Prepare the request to your authentication endpoint
	req, err := http.NewRequest("POST", authURL, nil)
	if err != nil {
		return "Error while authenticating", time.Time{}, err
	}

	// Set custom headers for authentication
	req.Header.Set("X-Auth-Username", authCache.Username)
	req.Header.Set("X-Auth-Password", authCache.Password)

	// Create a custom HTTP client with TLS/SSL certificate verification disabled
	httpClient := &http.Client{
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{
				InsecureSkipVerify: authCache.InsecureSkipVerify, // Note: Set to true for development purposes only
			},
		},
	}

	// Use the custom HTTP client to make the request
	resp, err := httpClient.Do(req)
	if err != nil {
		return "", time.Time{}, err
	}
	defer resp.Body.Close()

	// Check for non-successful response status
	if resp.StatusCode != http.StatusOK {
		return "", time.Time{}, errors.New("authentication failed: status code " + resp.Status)
	}

	// Read the response body
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", time.Time{}, fmt.Errorf("reading response body: %v", err)
	}

	// Unmarshal the JSON to get the token
	var authResp AuthResponse
	if err := json.Unmarshal(body, &authResp); err != nil {
		return "", time.Time{}, fmt.Errorf("unmarshaling response: %v", err)
	}

	// Extract the payload from the JWT token
	parts := strings.Split(authResp.Token, ".")
	if len(parts) != 3 {
		return "", time.Time{}, errors.New("invalid JWT token received")
	}
	payload, err := base64.RawURLEncoding.DecodeString(parts[1])
	if err != nil {
		return "", time.Time{}, err
	}

	// Parse the payload to extract the exp claim
	var claims struct {
		Exp int64 `json:"exp"`
	}
	if err := json.Unmarshal(payload, &claims); err != nil {
		return "", time.Time{}, err
	}

	// Convert the exp claim to time.Time
	expires := time.Unix(claims.Exp, 0)

	return authResp.Token, expires, nil
}
