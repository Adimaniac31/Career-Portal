package keycloak

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"

	"iiitn-career-portal/internal/config"
)

func CreateUser(cfg config.Config, email, password, name string) (string, error) {
	token, err := getAdminToken(cfg)
	if err != nil {
		return "", err
	}

	payload := map[string]interface{}{
		"username":        email,
		"email":           email,
		"enabled":         true,
		"emailVerified":   true,
		"firstName":       name,
		"requiredActions": []string{},
		"credentials": []map[string]interface{}{
			{
				"type":      "password",
				"value":     password,
				"temporary": false,
			},
		},
	}

	body, _ := json.Marshal(payload)

	req, _ := http.NewRequest(
		"POST",
		cfg.BaseURL+"/admin/realms/"+cfg.Realm+"/users",
		bytes.NewBuffer(body),
	)

	req.Header.Set("Authorization", "Bearer "+token)
	req.Header.Set("Content-Type", "application/json")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	// if resp.StatusCode != 201 {
	// 	return "", errors.New("failed to create keycloak user")
	// }
	if resp.StatusCode != 201 {
		body, _ := io.ReadAll(resp.Body)
		return "", fmt.Errorf(
			"keycloak user creation failed: status=%d body=%s",
			resp.StatusCode,
			string(body),
		)
	}

	location := resp.Header.Get("Location")
	parts := strings.Split(location, "/")
	return parts[len(parts)-1], nil
}

func DeleteUser(cfg config.Config, userID string) {
	token, err := getAdminToken(cfg)
	if err != nil {
		return
	}

	req, _ := http.NewRequest(
		"DELETE",
		cfg.BaseURL+"/admin/realms/"+cfg.Realm+"/users/"+userID,
		nil,
	)
	req.Header.Set("Authorization", "Bearer "+token)

	http.DefaultClient.Do(req)
}
