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

type FileInfo struct {
	Filename string
	Size     int64
	Changed  time.Time
	MimeType string
}

type PhotoResult struct {
	Previous *Photo
	Current  *Photo
	Next     *Photo
}
