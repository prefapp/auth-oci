package main

import (
	"slices"
)

func multilogin() {

	registries := parseRegistriesFromDir(CLI.Login.RegistriesDir)

	snapshotsRegistry := findRegistryByUrl(CLI.Login.SnapshotsRegistry, registries)

	releasesRegistry := findRegistryByUrl(CLI.Login.ReleasesRegistry, registries)

	creds := setCreds()

	for registryType, registry := range map[string]Registry{
		"snapshots": snapshotsRegistry,
		"releases":  releasesRegistry,
	} {

		var auth RegistryAuth

		switch registry.AuthStrategy {

		case "aws_oidc":
			auth = loginAWS()

		case "azure_oidc":
			auth = loginAzure(registry.RegistryHost)

		case "dockerhub":
			auth = RegistryAuth{
				Username: creds[registryType+"_user"],
				Password: creds[registryType+"_pass"],
				Registry: "https://index.docker.io/v1/",
			}

		case "ghcr", "generic":
			auth = RegistryAuth{
				Username: creds[registryType+"_user"],
				Password: creds[registryType+"_pass"],
				Registry: registry.RegistryHost,
			}

		default:
			panic("Unknown Auth Strategy: " + registry.AuthStrategy)
		}

		if slices.Contains(CLI.Login.Types, "helm") {

			helmLogin(auth)

		}

		if slices.Contains(CLI.Login.Types, "docker") {

			panic("Not implemented yet")

		}
	}

}

func setCreds() map[string]string {
	creds := map[string]string{
		"releases_user":  CLI.Login.ReleasesRegistryUsername,
		"releases_pass":  CLI.Login.ReleasesRegistryPassword,
		"snapshots_user": CLI.Login.SnapshotsRegistryUsername,
		"snapshots_pass": CLI.Login.SnapshotsRegistryPassword,
	}
	return creds
}
