package IBMStorwizeMetrics

import (
	"crypto/tls"
	"encoding/base64"
	"encoding/json"
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
	authURL := authCache.URL + "/auth"

	req, err := http.NewRequest("POST", authURL, nil)
	if err != nil {
		return "", time.Time{}, err
	}

	req.Header.Set("X-Auth-Username", authCache.Username)
	req.Header.Set("X-Auth-Password", authCache.Password)

	httpClient := &http.Client{
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{
				InsecureSkipVerify: authCache.InsecureSkipVerify,
			},
		},
	}

	resp, err := httpClient.Do(req)
	if err != nil {
		return "", time.Time{}, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", time.Time{}, fmt.Errorf("authentication failed: status code %s", resp.Status)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", time.Time{}, fmt.Errorf("reading response body: %v", err)
	}

	// fmt.Println("Auth response body:", string(body)) // 디버깅용

	// 임시 구조체로 응답 파싱
	var authResp struct {
		Token string `json:"token"`
	}
	if err := json.Unmarshal(body, &authResp); err != nil {
		return "", time.Time{}, fmt.Errorf("unmarshaling response: %v", err)
	}

	// JWT 포맷 여부 확인
	parts := strings.Split(authResp.Token, ".")
	if len(parts) == 3 {
		// JWT인 경우 exp 추출
		payload, err := base64.RawURLEncoding.DecodeString(parts[1])
		if err != nil {
			return "", time.Time{}, fmt.Errorf("decoding JWT payload: %v", err)
		}

		var claims struct {
			Exp int64 `json:"exp"`
		}
		if err := json.Unmarshal(payload, &claims); err != nil {
			return "", time.Time{}, fmt.Errorf("unmarshaling JWT payload: %v", err)
		}

		expires := time.Unix(claims.Exp, 0)
		return authResp.Token, expires, nil
	}

	// JWT 형식이 아닌 경우: 임의 만료시간 설정 (예: 1시간 후)
	return authResp.Token, time.Now().Add(1 * time.Hour), nil
}
