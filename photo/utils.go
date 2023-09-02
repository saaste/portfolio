package photo

import (
	"fmt"
	"log"
	"os"
	"path"
	"strings"

	"github.com/disintegration/imaging"
	"github.com/saaste/portfolio/settings"
)

func getSmallThumbnailFilename(fileName string) string {
	fileExt := path.Ext(fileName)
	fileNameWithoutExt := strings.TrimSuffix(fileName, fileExt)
	return fmt.Sprintf("%s_small%s", fileNameWithoutExt, fileExt)
}

func getMediumThumbnailFilename(fileName string) string {
	fileExt := path.Ext(fileName)
	fileNameWithoutExt := strings.TrimSuffix(fileName, fileExt)
	return fmt.Sprintf("%s_medium%s", fileNameWithoutExt, fileExt)
}

func isThumbnail(fileName string) bool {
	if strings.Contains(fileName, "_small.") || strings.Contains(fileName, "_medium.") {
		return true
	}
	return false
}

func hasThumbnails(fileName string) bool {
	if isThumbnail(fileName) {
		return true
	}

	smallThumbnail := fmt.Sprintf("files/%s", getSmallThumbnailFilename(fileName))
	mediumThumbnail := fmt.Sprintf("files/%s", getMediumThumbnailFilename(fileName))

	if _, err := os.Stat(smallThumbnail); err != nil {
		return false
	}

	if _, err := os.Stat(mediumThumbnail); err != nil {
		return false
	}

	return true
}

func generateThumbnails(fileName string, settings *settings.AppSettings) error {
	log.Printf("Generating thumbnails for %s\n", fileName)
	smallThumbnail := getSmallThumbnailFilename(fileName)
	mediumThumbnail := getMediumThumbnailFilename(fileName)
	fileFullPath := fmt.Sprintf("files/%s", fileName)

	source, err := imaging.Open(fileFullPath, imaging.AutoOrientation(true))
	if err != nil {
		return fmt.Errorf("failed to open the source file: %s: %w", fileFullPath, err)
	}

	smallTn := imaging.Fit(source, settings.SmallSize, settings.SmallSize, imaging.Lanczos)
	mediumTn := imaging.Fit(source, settings.MediumSize, settings.MediumSize, imaging.Lanczos)

	err = imaging.Save(smallTn, fmt.Sprintf("files/%s", smallThumbnail))
	if err != nil {
		return fmt.Errorf("failed to save small thumbnail: %s, %w", smallThumbnail, err)
	}
	err = imaging.Save(mediumTn, fmt.Sprintf("files/%s", mediumThumbnail))
	if err != nil {
		return fmt.Errorf("failed to save medium thumbnail: %s, %w", smallThumbnail, err)
	}

	return nil
}
