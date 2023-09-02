package settings

import "time"

type AppSettings struct {
	SmallSize       int           `yaml:"smallSize"`
	MediumSize      int           `yaml:"mediumSize"`
	RefreshInterval time.Duration `yaml:"refreshInterval"`
	Title           string        `yaml:"title"`
	Author          string        `yaml:"author"`
	About           string        `yaml:"about"`
}
