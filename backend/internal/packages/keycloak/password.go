package keycloak

import (
	"encoding/json"
	"errors"
	"io"
	"log"
	"net/http"
	"net/url"

	"iiitn-career-portal/internal/config"
)

type TokenResponse struct {
	AccessToken string `json:"access_token"`
	ExpiresIn   int    `json:"expires_in"`
	TokenType   string `json:"token_type"`
}

func PasswordGrant(cfg config.Config, email, password string) (*TokenResponse, error) {
	data := url.Values{}
	data.Set("grant_type", "password")
	data.Set("client_id", cfg.ClientID)
	data.Set("client_secret", cfg.ClientSecret)
	data.Set("username", email)
	data.Set("password", password)

	tokenURL := cfg.BaseURL +
		"/realms/" + cfg.Realm +
		"/protocol/openid-connect/token"

	resp, err := http.PostForm(tokenURL, data)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		body, _ := io.ReadAll(resp.Body)
		log.Println("[KC] password grant failed:", resp.StatusCode, string(body))
		return nil, errors.New("invalid credentials")
	}

	var token TokenResponse
	if err := json.NewDecoder(resp.Body).Decode(&token); err != nil {
		return nil, err
	}

	return &token, nil
}
