package main

var (
	CREDS = map[string]string{
		"releases_user":  CLI.Login.ReleasesRegistryUsername,
		"releases_pass":  CLI.Login.ReleasesRegistryPassword,
		"snapshots_user": CLI.Login.SnapshotsRegistryUsername,
		"snapshots_pass": CLI.Login.SnapshotsRegistryPassword,
	}
)
