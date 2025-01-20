package main

var CLI struct {
	Login struct {
		RegistriesDir             string   `help:"The directory where the registries files are stored. Example: ./registries"`
		ReleasesRegistry          string   `help:"The registry host to use for releases.Example: ghcr.io"`
		SnapshotsRegistry         string   `help:"The registry host to use for snapshots. Example: ghcr.io"`
		ReleasesRegistryUsername  string   `help:"If you use dockerhub, generic, or ghcr, you need to provide a username."`
		ReleasesRegistryPassword  string   `help:"If you use dockerhub, generic, or ghcr, you need to provide a password."`
		SnapshotsRegistryUsername string   `help:"If you use dockerhub, generic, or ghcr, you need to provide a username."`
		SnapshotsRegistryPassword string   `help:"If you use dockerhub, generic, or ghcr, you need to provide a password."`
		Types                     []string `help:"Types of technologies to authenticate to. Options: docker, helm"`
	} `cmd:"" help:"Authenticate to helm registries Azure, AWS, Dockerhub, Ghcr"`
	Validate struct {
		All           bool   `help:"Validate all registries."`
		RegistriesDir string `help:"The directory where the registries files are stored."`
		File          string `help:"Validate a specific registry file."`
	} `cmd:"" help:"Validate the registries files."`
}
