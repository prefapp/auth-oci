package main

import (
	"context"
	"fmt"
	"log"
	"net/url"
	"os"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore/policy"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/to"
	"github.com/Azure/azure-sdk-for-go/sdk/azidentity"
	"github.com/Azure/azure-sdk-for-go/sdk/containers/azcontainerregistry"
)

type RespGhAz struct {
	Value string `json:"value"`
}

func loginAzure(registry string) RegistryAuth {

	if isGhRunner() {

		log.Printf("Requesting federated token")

		preparedOidcEnvAz()

	} else {

		log.Printf("Skipping OIDC auth, not running in GitHub Runner")

	}

	cred, err := azidentity.NewDefaultAzureCredential(nil)

	if err != nil {

		panic(fmt.Sprintf("Failed to authenticate: %v", err))

	}

	ctx := context.Background()

	aadToken, err := cred.GetToken(
		ctx, policy.TokenRequestOptions{
			Scopes: []string{"https://containerregistry.azure.net/.default"},
		},
	)

	if err != nil {
		panic(fmt.Sprintf("Failed to get AAD token: %v", err))
	}

	registryHost, err := getRegistryHostname(registry)

	if err != nil {

		panic(fmt.Sprintf("Failed to get registry hostname: %v", err))

	}

	ac, err := azcontainerregistry.NewAuthenticationClient(
		fmt.Sprintf("https://%s", registryHost),
		&azcontainerregistry.AuthenticationClientOptions{},
	)

	if err != nil {

		panic(fmt.Sprintf("Failed to create authentication client: %v", err))

	}

	rt, err := ac.ExchangeAADAccessTokenForACRRefreshToken(
		ctx, "access_token",
		registryHost,
		&azcontainerregistry.AuthenticationClientExchangeAADAccessTokenForACRRefreshTokenOptions{
			AccessToken: to.Ptr(aadToken.Token),
			Tenant:      to.Ptr(os.Getenv("AZURE_TENANT_ID")),
		},
	)

	if err != nil {

		panic(
			fmt.Sprintf(
				"Failed to exchange AAD access token for ACR refresh token: %v",
				err,
			),
		)

	}

	return RegistryAuth{
		Username: "00000000-0000-0000-0000-000000000000",
		Password: *rt.ACRRefreshToken.RefreshToken,
		Registry: registry,
	}
}

func preparedOidcEnvAz() {

	audience := "api://AzureADTokenExchange"

	tokenFile := "/tmp/ghaz_token"

	token := getOIDCToken(audience)

	saveTokenToFile(tokenFile, token)

	os.Setenv("AZURE_FEDERATED_TOKEN_FILE", tokenFile)

}

func getRegistryHostname(templatePath string) (string, error) {
	registryURL, err := url.Parse("https://" + templatePath)
	if err != nil {
		return "", err
	}
	return registryURL.Host, nil
}
