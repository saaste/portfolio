package photo

import (
	"fmt"
	"log"
	"os"
	"path"
	"sort"
	"strings"

	"github.com/djherbis/times"
	"github.com/saaste/portfolio/settings"
)

func GetPhotos(settings *settings.AppSettings, forceThumbnails bool) ([]Photo, error) {
	photos := make([]Photo, 0)

	entries, err := os.ReadDir("files/")
	if err != nil {
		return photos, fmt.Errorf("failed to read files directory: %w", err)
	}

	for _, entry := range entries {
		if entry.IsDir() {
			continue
		}
		if isThumbnail(entry.Name()) {
			continue
		}

		fullPath := fmt.Sprintf("files/%s", entry.Name())
		fileInfo, err := times.Stat(fullPath)
		if err != nil {
			return photos, fmt.Errorf("failed to get file info for %s", fullPath)
		}

		switch path.Ext(strings.ToLower(entry.Name())) {
		case ".jpg", ".jpeg", ".png":
			hasThumbnals := hasThumbnails(entry.Name())
			if !hasThumbnals || forceThumbnails {
				generateThumbnails(entry.Name(), settings)
			}
			altText := getAltText(fullPath)
			photos = append(photos, Photo{
				FullFileName:   entry.Name(),
				MediumFileName: getMediumThumbnailFilename(entry.Name()),
				SmallFileName:  getSmallThumbnailFilename(entry.Name()),
				Changed:        fileInfo.ChangeTime(),
				AltText:        altText,
			})
		case ".txt":
			continue
		default:
			log.Printf("WARNING: unsupported file type: %s\n", entry.Name())
		}
	}

	if !forceThumbnails {
		log.Printf("Photo list refreshed: %d photos found\n", len(photos))
	}

	sort.Slice(photos, func(i, j int) bool {
		return photos[i].Changed.After(photos[j].Changed)
	})

	return photos, nil
}

func getAltText(photoFullPath string) string {
	photoExt := path.Ext(photoFullPath)
	altFile := strings.Replace(photoFullPath, photoExt, ".txt", -1)

	_, err := os.Stat(altFile)
	if err != nil {
		log.Printf("WARNING: No alt text for %s", photoFullPath)
		return ""
	}

	content, err := os.ReadFile(altFile)
	if err != nil {
		log.Printf("ERROR: failed to read alt text file %s: %v", altFile, err)
		return ""
	}

	return string(content)

}
