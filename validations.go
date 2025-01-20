package main

import (
	"fmt"
	"log"
	"os"

	"github.com/xeipuuv/gojsonschema"
	"gopkg.in/yaml.v3"
)

const SCHEMA = `{
	"type": "object",
	"properties": {
	  "name": { "type": "string" },
	  "registry": { "type": "string" },
	  "image_types": {
		"type": "array",
		"items": { "type": "string", "enum": ["snapshots", "releases"] }
	  },
	  "default": { "type": "boolean" },
	  "auth_strategy": {
		"type": "string",
		"enum": ["aws_oidc", "azure_oidc", "generic", "ghcr", "dockerhub"]
	  },
	  "base_paths": {
		"type": "object",
		"properties": {
		  "services": { "type": "string" },
		  "charts": { "type": "string" }
		},
		"required": ["services", "charts"]
	  }
	},
	"required": ["name", "registry", "image_types", "default", "auth_strategy", "base_paths"]
  }`

func validate() {

	if CLI.Validate.All && CLI.Validate.File != "" {

		panic("You can't use --all and --file at the same time")

	}

	if CLI.Validate.All && CLI.Validate.RegistriesDir == "" {

		panic("--all requires --registries-dir, please provide it")

	}

	if CLI.Validate.File != "" && CLI.Validate.RegistriesDir != "" {

		panic("--file and --registries-dir are mutually exclusive, please provide only one")

	}

	if CLI.Validate.All {

		log.Printf("Validating all registries in %s", CLI.Validate.RegistriesDir)

		registries := parseRegistriesFromDir(CLI.Validate.RegistriesDir)

		log.Printf("%d registries found", len(registries))

		log.Printf("All registries are valid")

	} else if CLI.Validate.File != "" {

		validateRegistryFileSchema(CLI.Validate.File)

	} else {

		panic("You need to provide a registry --file or the --all flag")

	}
}

func validateRegistryAgainstSchema(registryData Registry) error {

	schemaLoader := gojsonschema.NewStringLoader(SCHEMA)

	documentLoader := gojsonschema.NewGoLoader(registryData)

	result, err := gojsonschema.Validate(schemaLoader, documentLoader)

	if err != nil {

		return fmt.Errorf("validation: %v", err)

	}

	if !result.Valid() {

		for _, desc := range result.Errors() {

			fmt.Printf("Error: %s\n", desc)

		}

		return fmt.Errorf("document is not valid %s", result.Errors())

	}
	return nil
}

func validateRegistryFileSchema(registryFile string) (Registry, error) {

	var registry Registry

	log.Printf("Reading registry file %s", registryFile)

	fileData, err := os.ReadFile(registryFile)

	if err != nil {

		return Registry{}, err
	}

	log.Printf("Parsing registry file %s", registryFile)

	err = yaml.Unmarshal(fileData, &registry)

	if err != nil {

		return Registry{}, err

	}

	log.Printf("Validating registry file %s", registryFile)

	err = validateRegistryAgainstSchema(registry)

	if err != nil {

		log.Printf("Registry file %s is invalid", registryFile)

		return Registry{}, err

	}

	log.Printf("Registry file %s is valid", registryFile)

	return registry, nil
}
