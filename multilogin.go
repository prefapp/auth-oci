package main

import "slices"

func multilogin() {

	registries := parseRegistriesFromDir(CLI.Login.RegistriesDir)

	snapshotsRegistry := findRegistryByUrl(CLI.Login.SnapshotsRegistry, registries)

	releasesRegistry := findRegistryByUrl(CLI.Login.ReleasesRegistry, registries)

	for registryType, registry := range map[string]Registry{"snapshots": snapshotsRegistry, "releases": releasesRegistry} {

		var auth RegistryAuth

		switch registry.AuthStrategy {

		case "aws_oidc":
			auth = loginAWS()

		case "azure_oidc":
			auth = loginAzure(registry.RegistryHost)

		case "dockerhub":
			auth = RegistryAuth{
				Username: CREDS[registryType+"_user"],
				Password: CREDS[registryType+"_pass"],
				Registry: "https://index.docker.io/v1/",
			}

		case "ghcr":
			auth = RegistryAuth{
				Username: CREDS[registryType+"_user"],
				Password: CREDS[registryType+"_pass"],
				Registry: "https://ghcr.io/v2/",
			}

		case "generic":
			auth = RegistryAuth{
				Username: CREDS[registryType+"_user"],
				Password: CREDS[registryType+"_pass"],
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
