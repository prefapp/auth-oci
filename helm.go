package main

import (
	"log"

	"helm.sh/helm/v3/pkg/action"
	"helm.sh/helm/v3/pkg/cli"
	"helm.sh/helm/v3/pkg/registry"
)

type RegistryAuth struct {
	Username string `help:"Username for the registry."`
	Password string `help:"Password for the registry."`
	Registry string `help:"Registry URL."`
}

func helmLogin(auth RegistryAuth) {

	settings := cli.New()

	actionConfig := new(action.Configuration)

	err := actionConfig.Init(settings.RESTClientGetter(), settings.Namespace(), "", log.Printf)

	if err != nil {

		log.Fatalf("Failed to initialize Helm action configuration: %v", err)

	}

	regClient, err := registry.NewClient()

	if err != nil {

		log.Fatalf("Failed to create Helm registry client: %v", err)

	}

	err = regClient.Login(
		auth.Registry,
		registry.LoginOptBasicAuth(auth.Username, auth.Password),
	)

	if err != nil {

		log.Fatalf("Failed to log in to registry: %v", err)

	} else {

		log.Printf("Successfully logged in to registry %s", auth.Registry)

	}
}
