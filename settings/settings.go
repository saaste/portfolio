package settings

import (
	"fmt"
	"os"
	"strings"

	"gopkg.in/yaml.v3"
)

func ReadSettings() (*AppSettings, error) {
	data, err := os.ReadFile("settings.yaml")
	if err != nil {
		return nil, fmt.Errorf("failed to read settings.yaml: %v", err)
	}

	var appSettings AppSettings
	if err := yaml.Unmarshal(data, &appSettings); err != nil {
		return nil, fmt.Errorf("failed to unmarshal settings.yaml: %v", err)
	}

	appSettings.BaseURL = strings.TrimSuffix(appSettings.BaseURL, "/")

	return &appSettings, nil
}
