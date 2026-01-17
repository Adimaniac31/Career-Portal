package keycloak

import (
	"encoding/json"
	"errors"
	"net/http"
	"net/url"

	"iiitn-career-portal/internal/config"
)

func getAdminToken(cfg config.Config) (string, error) {
	data := url.Values{}
	data.Set("grant_type", "client_credentials")
	data.Set("client_id", cfg.ClientID)
	data.Set("client_secret", cfg.ClientSecret)

	tokenURL := cfg.BaseURL +
		"/realms/" + cfg.Realm +
		"/protocol/openid-connect/token"

	resp, err := http.PostForm(tokenURL, data)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		return "", errors.New("failed to get keycloak admin token")
	}

	var res struct {
		AccessToken string `json:"access_token"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&res); err != nil {
		return "", err
	}

	return res.AccessToken, nil
}
