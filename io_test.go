package main

import (
	"fmt"
	"testing"
)

func TestRegistriesParser(t *testing.T) {

	registries := parseRegistriesFromDir("./test/fixtures/registries")

	fmt.Printf("Registries: %v\n", registries[0].RegistryHost)

	t.Run("Registries number are ok", func(t *testing.T) {
		if len(registries) != 2 {
			t.Errorf("Expected 2 registries, got %v", len(registries))
		}
	})

	t.Run("Registries are parsed correctly", func(t *testing.T) {
		if registries[0].RegistryHost != "acrsnapshots2.azurecr.io" {
			t.Errorf("Expected acrsnapshots2.azurecr.io, got %v", registries[0].RegistryHost)
		}

		if registries[1].RegistryHost != "acrsnapshots.azurecr.io" {
			t.Errorf("Expected acrsnapshots.azurecr.io, got %v", registries[1].RegistryHost)
		}
	})

	t.Run("Registries finder works as expected", func(t *testing.T) {

		registry := findRegistryByUrl("acrsnapshots2.azurecr.io", registries)

		fmt.Printf("Registry: %v\n", registry)

		if registry.RegistryHost != "acrsnapshots2.azurecr.io" {
			t.Errorf("Expected acrsnapshots2.azurecr.io, got %v", registry.RegistryHost)
		}
	})

}
