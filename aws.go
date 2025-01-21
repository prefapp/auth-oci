package main

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/ecr"
	"github.com/aws/aws-sdk-go-v2/service/sts"
)

type OIDCResponse struct {
	Value string `json:"value"`
}

func loginAWS() RegistryAuth {

	audience := "sts.amazonaws.com"
	tokenFile := "/tmp/awsjwt"

	token := getOIDCToken(audience)

	saveTokenToFile(tokenFile, token)

	os.Setenv("AWS_WEB_IDENTITY_TOKEN_FILE", tokenFile)

	roleArn := os.Getenv("AWS_ROLE_ARN")

	if roleArn == "" {
		log.Fatalf("AWS_ROLE_ARN is not set")
	}

	os.Setenv("AWS_ROLE_ARN", roleArn)

	cfg, err := config.LoadDefaultConfig(context.TODO())

	if err != nil {

		log.Fatalf("Failed to load AWS configuration: %v", err)

	}

	ecrClient := ecr.NewFromConfig(cfg)

	resp, err := ecrClient.GetAuthorizationToken(context.TODO(), &ecr.GetAuthorizationTokenInput{})

	if err != nil {
		panic(fmt.Sprintf("Failed to get ECR authorization token: %v", err))
	}

	if len(resp.AuthorizationData) == 0 {
		panic("No authorization data returned by ECR")
	}

	authData := resp.AuthorizationData[0]

	decodedToken, err := base64.StdEncoding.DecodeString(*authData.AuthorizationToken)

	if err != nil {
		panic(fmt.Sprintf("Failed to decode authorization token: %v", err))
	}

	parts := strings.SplitN(string(decodedToken), ":", 2)
	if len(parts) != 2 {
		panic("Invalid authorization token format")
	}

	registryURL := *authData.ProxyEndpoint

	return RegistryAuth{
		Username: parts[0],
		Password: parts[1],
		Registry: registryURL,
	}
}

func getOIDCToken(audience string) string {

	requestURL := os.Getenv("ACTIONS_ID_TOKEN_REQUEST_URL") + "&audience=" + audience
	requestToken := os.Getenv("ACTIONS_ID_TOKEN_REQUEST_TOKEN")

	if requestURL == "" || requestToken == "" {

		log.Fatalf("ACTIONS_ID_TOKEN_REQUEST_URL or ACTIONS_ID_TOKEN_REQUEST_TOKEN is not set")

	}

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

func saveTokenToFile(filename, token string) {
	err := os.WriteFile(filename, []byte(token), 0644)
	if err != nil {
		log.Fatalf("Failed to write token to file %s: %v", filename, err)
	}
}

func getCallerIdentity(cfg aws.Config) {

	stsClient := sts.NewFromConfig(cfg)

	resp, err := stsClient.GetCallerIdentity(context.TODO(), &sts.GetCallerIdentityInput{})

	if err != nil {

		log.Fatalf("Failed to get caller identity: %v", err)

	}

	fmt.Printf("Account: %s\n", *resp.Account)
	fmt.Printf("UserID: %s\n", *resp.UserId)
	fmt.Printf("ARN: %s\n", *resp.Arn)
}
