package main

import "encoding/json"

type Registry struct {
	Name         string            `yaml:"name" json:"name"`
	Url          string            `yaml:"registry,omitempty" json:"registry,omitempty"`
	AuthStrategy string            `yaml:"auth_strategy" json:"auth_strategy,omitempty"`
	Default      bool              `yaml:"default" json:"default"`
	ImageTypes   []string          `yaml:"image_types" json:"image_types"`
	BasePaths    map[string]string `yaml:"base_paths" json:"base_paths"`
}

// Custom YAML unmarshaling to support both 'registry' and 'url' mapping to RegistryHost
func (r *Registry) UnmarshalYAML(unmarshal func(interface{}) error) error {
	var aux map[string]interface{}
	if err := unmarshal(&aux); err != nil {
		return err
	}
	if name, ok := aux["name"].(string); ok {
		r.Name = name
	}
	if reg, ok := aux["registry"].(string); ok {
		r.Url = reg
	} else if url, ok := aux["url"].(string); ok {
		r.Url = url
	}
	if auth, ok := aux["auth_strategy"].(string); ok {
		r.AuthStrategy = auth
	}
	if def, ok := aux["default"].(bool); ok {
		r.Default = def
	}
	if img, ok := aux["image_types"].([]interface{}); ok {
		r.ImageTypes = make([]string, len(img))
		for i, v := range img {
			if s, ok := v.(string); ok {
				r.ImageTypes[i] = s
			}
		}
	}
	if bp, ok := aux["base_paths"].(map[string]interface{}); ok {
		r.BasePaths = make(map[string]string)
		for k, v := range bp {
			if s, ok := v.(string); ok {
				r.BasePaths[k] = s
			}
		}
	}
	return nil
}

// Custom JSON unmarshaling to support both 'registry' and 'url' mapping to RegistryHost
func (r *Registry) UnmarshalJSON(data []byte) error {
	type Alias Registry
	var aux map[string]interface{}
	if err := json.Unmarshal(data, &aux); err != nil {
		return err
	}
	if name, ok := aux["name"].(string); ok {
		r.Name = name
	}
	if reg, ok := aux["registry"].(string); ok {
		r.Url = reg
	} else if url, ok := aux["url"].(string); ok {
		r.Url = url
	}
	if auth, ok := aux["auth_strategy"].(string); ok {
		r.AuthStrategy = auth
	}
	if def, ok := aux["default"].(bool); ok {
		r.Default = def
	}
	if img, ok := aux["image_types"].([]interface{}); ok {
		r.ImageTypes = make([]string, len(img))
		for i, v := range img {
			if s, ok := v.(string); ok {
				r.ImageTypes[i] = s
			}
		}
	}
	if bp, ok := aux["base_paths"].(map[string]interface{}); ok {
		r.BasePaths = make(map[string]string)
		for k, v := range bp {
			if s, ok := v.(string); ok {
				r.BasePaths[k] = s
			}
		}
	}
	return nil
}
