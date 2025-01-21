package main

import (
	"context"
	"encoding/base64"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/ecr"
)

func loginAWS() RegistryAuth {

	if isGhRunner() {

		log.Printf("GitHub runner detected, using OIDC auth")

		prepareOidcEnvAws()

	} else {

		log.Printf("Skipping OIDC auth, not running in GitHub Runner")
	}

	cfg, err := config.LoadDefaultConfig(context.TODO())

	if err != nil {

		log.Fatalf("Failed to load AWS configuration: %v", err)

	}

	ecrClient := ecr.NewFromConfig(cfg)

	resp, err := ecrClient.GetAuthorizationToken(context.TODO(), &ecr.GetAuthorizationTokenInput{})

	if err != nil {
		panic(
			fmt.Sprintf(
				"Failed to get ECR authorization token: %v",
				err,
			),
		)
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

func prepareOidcEnvAws() {

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
}
