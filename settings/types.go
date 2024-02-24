package settings

import "time"

type AppSettings struct {
	SmallSize       int           `yaml:"smallSize"`
	MediumSize      int           `yaml:"mediumSize"`
	RefreshInterval time.Duration `yaml:"refreshInterval"`
	BaseURL         string        `yaml:"baseUrl" default:""`
	Title           string        `yaml:"title"`
	Description     string        `yaml:"description" default:""`
	Author          string        `yaml:"author"`
	About           string        `yaml:"about"`
}
