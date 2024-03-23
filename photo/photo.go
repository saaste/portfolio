package photo

import (
	"errors"
	"fmt"
	"log"
	"os"
	"path"
	"sort"
	"strings"

	"github.com/djherbis/times"
	"github.com/saaste/portfolio/settings"
	"gopkg.in/yaml.v3"
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

		entryInfo, err := entry.Info()
		if err != nil {
			if !errors.Is(err, os.ErrNotExist) {
				return photos, fmt.Errorf("failed to get file info: %w", err)
			}
		}

		fullPath := fmt.Sprintf("files/%s", entry.Name())
		fileInfo, err := times.Stat(fullPath)
		if err != nil {
			return photos, fmt.Errorf("failed to get file info for %s", fullPath)
		}

		mimeType := ""
		fileExtension := path.Ext(strings.ToLower(entry.Name()))
		switch fileExtension {
		case ".jpg", ".jpeg":
			mimeType = "image/jpeg"
		case ".png":
			mimeType = "image/png"
		}

		switch fileExtension {
		case ".jpg", ".jpeg", ".png":

			hasThumbnals := hasThumbnails(entry.Name())
			if !hasThumbnals || forceThumbnails {
				generateThumbnails(entry.Name(), settings)
			}
			photoInfo := getPhotoInfo(fullPath)
			photos = append(photos, Photo{
				FullFileName:   entry.Name(),
				MediumFileName: getMediumThumbnailFilename(entry.Name()),
				SmallFileName:  getSmallThumbnailFilename(entry.Name()),
				Changed:        fileInfo.ChangeTime(),
				Size:           entryInfo.Size(),
				PhotoInfo:      photoInfo,
				MimeType:       mimeType,
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

func getPhotoInfo(photoFullPath string) PhotoInfo {
	photoInfo := PhotoInfo{}

	photoExt := path.Ext(photoFullPath)
	infoFile := strings.Replace(photoFullPath, photoExt, ".yaml", -1)

	data, err := os.ReadFile(infoFile)
	if err != nil {
		log.Printf("WARNING: failed to read info file %s: %v", infoFile, err)
		return photoInfo
	}

	if err := yaml.Unmarshal(data, &photoInfo); err != nil {
		log.Printf("ERROR: failed to unmarshal info file %s: %v", infoFile, err)
		return photoInfo
	}

	return photoInfo

}
