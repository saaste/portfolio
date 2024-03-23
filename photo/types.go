package photo

import (
	"time"
)

type Photo struct {
	FullFileName   string
	MediumFileName string
	SmallFileName  string
	Changed        time.Time
	PhotoInfo      PhotoInfo
	Size           int64
	MimeType       string
}

type PhotoInfo struct {
	Title       string `yaml:"title,omitempty"`
	Description string `yaml:"description,omitempty"`
	AltText     string `yaml:"altText,omitempty"`
}
