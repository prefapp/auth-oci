package main

import (
	"encoding/json"
	"io"
	"log"
	"net/http"
	"os"
)

func isGhRunner() bool {

	return os.Getenv("ACTIONS_ID_TOKEN_REQUEST_URL") != "" && os.Getenv("ACTIONS_ID_TOKEN_REQUEST_TOKEN") != ""

}

type OIDCResponse struct {
	Value string `json:"value"`
}

func getOIDCToken(audience string) string {

	requestURL := os.Getenv("ACTIONS_ID_TOKEN_REQUEST_URL") + "&audience=" + audience

	requestToken := os.Getenv("ACTIONS_ID_TOKEN_REQUEST_TOKEN")

	req, err := http.NewRequest("GET", requestURL, nil)

	if err != nil {

		log.Fatalf("Failed to create HTTP request: %v", err)

	}

	req.Header.Set("Authorization", "bearer "+requestToken)

	client := &http.Client{}

	resp, err := client.Do(req)

	if err != nil {

		log.Fatalf("Failed to make HTTP request: %v", err)

	}

	defer resp.Body.Close()

	if resp.StatusCode != 200 {

		log.Fatalf("Failed to get OIDC token: %v", resp.Status)

	}

	body, err := io.ReadAll(resp.Body)

	if err != nil {

		log.Fatalf("Failed to read response body: %v", err)

	}

	var oidcResponse OIDCResponse

	if err := json.Unmarshal(body, &oidcResponse); err != nil {

		log.Fatalf("Failed to parse response JSON: %v", err)

	}

	return oidcResponse.Value
}
