package main

type Registry struct {
	Name         string            `yaml:"name" json:"name"`
	RegistryHost string            `yaml:"registry" json:"registry"`
	AuthStrategy string            `yaml:"auth_strategy" json:"auth_strategy"`
	Default      bool              `yaml:"default" json:"default"`
	ImageTypes   []string          `yaml:"image_types" json:"image_types"`
	BasePaths    map[string]string `yaml:"base_paths" json:"base_paths"`
}
