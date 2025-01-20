package main

import (
	"fmt"
	"os"
	"path/filepath"
)

func parseRegistriesFromDir(dir string) []Registry {

	var registries []Registry

	mapRegistriesByHost := make(map[string]bool)

	mapRegistriesByName := make(map[string]bool)

	filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {

		if err != nil {
			return err
		}

		if !info.IsDir() && (filepath.Ext(path) == ".yaml" || filepath.Ext(path) == ".yml") {

			registry, err := validateRegistryFileSchema(path)

			if err != nil {

				panic(err)

			}

			if mapRegistriesByHost[registry.RegistryHost] {

				panic(fmt.Sprintf("Duplicated registry %s", registry.RegistryHost))

			}

			if mapRegistriesByName[registry.Name] {

				panic(fmt.Sprintf("Duplicated registry %s", registry.Name))

			}

			mapRegistriesByHost[registry.RegistryHost] = true
			mapRegistriesByName[registry.Name] = true

			registries = append(registries, registry)
		}

		return nil
	})

	return registries
}

func findRegistryByUrl(url string, registries []Registry) Registry {

	for _, r := range registries {

		if r.RegistryHost == url {

			return r

		}
	}

	panic(fmt.Sprintf("Registry %s not found", url))
}
