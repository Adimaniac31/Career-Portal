package keycloak

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"maps"
	"net/http"

	"iiitn-career-portal/internal/config"
)

func AssignRealmRole(cfg config.Config, userID, roleName string) error {
	log.Println("[KC] AssignRealmRole start",
		"userID=", userID,
		"role=", roleName,
	)

	// 0️⃣ Get admin token
	token, err := getAdminToken(cfg)
	if err != nil {
		log.Println("[KC] failed to get admin token:", err)
		return err
	}
	log.Println("[KC] admin token acquired")

	// 1️⃣ Fetch role representation
	roleURL := cfg.BaseURL +
		"/admin/realms/" + cfg.Realm +
		"/roles/" + roleName

	log.Println("[KC] fetching role from:", roleURL)

	roleReq, _ := http.NewRequest("GET", roleURL, nil)
	roleReq.Header.Set("Authorization", "Bearer "+token)

	roleResp, err := http.DefaultClient.Do(roleReq)
	if err != nil {
		log.Println("[KC] role fetch request failed:", err)
		return err
	}
	defer roleResp.Body.Close()

	roleBody, _ := io.ReadAll(roleResp.Body)

	log.Println("[KC] role fetch status:", roleResp.StatusCode)
	log.Println("[KC] role fetch body:", string(roleBody))

	if roleResp.StatusCode != 200 {
		return fmt.Errorf(
			"failed to fetch role: status=%d body=%s",
			roleResp.StatusCode,
			string(roleBody),
		)
	}

	var role map[string]interface{}
	if err := json.Unmarshal(roleBody, &role); err != nil {
		log.Println("[KC] failed to decode role JSON:", err)
		return err
	}

	log.Println("[KC] role object keys:", maps.Keys(role))
	log.Println("[KC] role id:", role["id"], "name:", role["name"])

	// 2️⃣ Assign role to user
	payload, _ := json.Marshal([]map[string]interface{}{role})
	log.Println("[KC] role assignment payload:", string(payload))

	assignURL := cfg.BaseURL +
		"/admin/realms/" + cfg.Realm +
		"/users/" + userID +
		"/role-mappings/realm"

	log.Println("[KC] assigning role via:", assignURL)

	assignReq, _ := http.NewRequest(
		"POST",
		assignURL,
		bytes.NewBuffer(payload),
	)
	assignReq.Header.Set("Authorization", "Bearer "+token)
	assignReq.Header.Set("Content-Type", "application/json")

	assignResp, err := http.DefaultClient.Do(assignReq)
	if err != nil {
		log.Println("[KC] role assignment request failed:", err)
		return err
	}
	defer assignResp.Body.Close()

	assignBody, _ := io.ReadAll(assignResp.Body)

	log.Println("[KC] role assignment status:", assignResp.StatusCode)
	log.Println("[KC] role assignment body:", string(assignBody))

	if assignResp.StatusCode != 204 {
		return fmt.Errorf(
			"role assignment failed: status=%d body=%s",
			assignResp.StatusCode,
			string(assignBody),
		)
	}

	log.Println("[KC] role assignment successful")
	return nil
}
