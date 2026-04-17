package main

import (
	"testing"
)

func TestValidateRegistryAgainstSchema(t *testing.T) {
	validRegistry := Registry{
		Name:         "example",
		RegistryHost: "https://example.com",
		ImageTypes:   []string{"snapshots", "releases"},
		Default:      true,
		AuthStrategy: "aws_oidc",
		BasePaths: map[string]string{
			"services": "/services/path",
			"charts":   "/charts/path",
		},
	}

	t.Run("Valid registry", func(t *testing.T) {
		if err := validateRegistryAgainstSchema(validRegistry); err != nil {
			t.Errorf("Expected no error, got %v", err)
		}
	})

}

func TestValidateRegistryAgainstSchemaInvalid(t *testing.T) {
	invalidRegistry := Registry{
		Name:         "example",
		RegistryHost: "https://example.com",
		ImageTypes:   []string{"NOT_VALID", ""},
		Default:      true,
		AuthStrategy: "NOT_VALID_STRATEGY",
		BasePaths: map[string]string{
			"services": "/services/path",
		},
	}

	t.Run("Invalid registry", func(t *testing.T) {
		if err := validateRegistryAgainstSchema(invalidRegistry); err == nil {
			t.Errorf("Expected error, got nil")
		} else {
			t.Log(err)
		}
	})
}
