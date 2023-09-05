package photo

import (
	"time"
)

type Photo struct {
	FullFileName   string
	MediumFileName string
	SmallFileName  string
	Changed        time.Time
	AltText        string
}
