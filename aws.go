package main

import (
	"context"
	"encoding/base64"
	"fmt"
	"strings"

	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/ecr"
)

func loginAWS() RegistryAuth {

	cfg, err := config.LoadDefaultConfig(context.TODO())

	if err != nil {
		panic(fmt.Sprintf("Failed to load AWS configuration: %v", err))
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
