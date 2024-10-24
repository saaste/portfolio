package photo

import (
	"fmt"
	"io/fs"
	"log"
	"os"
	"path"
	"sort"
	"strings"

	"github.com/disintegration/imaging"
	"github.com/djherbis/times"
	"github.com/saaste/portfolio/pkg/settings"
	"gopkg.in/yaml.v3"
)

func (p *PhotoRepo) getSmallThumbnailFilename(fileName string) string {
	fileExt := path.Ext(fileName)
	fileNameWithoutExt := strings.TrimSuffix(fileName, fileExt)
	return fmt.Sprintf("%s_small%s", fileNameWithoutExt, fileExt)
}

func (p *PhotoRepo) getMediumThumbnailFilename(fileName string) string {
	fileExt := path.Ext(fileName)
	fileNameWithoutExt := strings.TrimSuffix(fileName, fileExt)
	return fmt.Sprintf("%s_medium%s", fileNameWithoutExt, fileExt)
}

func (p *PhotoRepo) isThumbnail(fileName string) bool {
	if strings.Contains(fileName, "_small.") || strings.Contains(fileName, "_medium.") {
		return true
	}
	return false
}

func (p *PhotoRepo) isOriginalPhoto(filename string) bool {
	filenameLower := strings.ToLower(filename)

	// Check that it is an image
	if !strings.HasSuffix(filenameLower, ".jpg") && !strings.HasSuffix(filenameLower, ".jpeg") && !strings.HasSuffix(filenameLower, ".png") {
		return false
	}

	if strings.Contains(filenameLower, "_small.") || strings.Contains(filenameLower, "_medium.") {
		return false
	}

	return true
}

func (p *PhotoRepo) getOriginalPhotos(entries []fs.DirEntry) ([]FileInfo, error) {
	result := make([]FileInfo, 0)

	// Filter original photos
	for _, entry := range entries {
		if !entry.IsDir() && p.isOriginalPhoto(entry.Name()) {
			fullPath := fmt.Sprintf("files/%s", entry.Name())
			fileInfo, err := times.Stat(fullPath)
			if err != nil {
				return nil, fmt.Errorf("failed to fetch stat: %w", err)
			}

			entryInfo, err := entry.Info()
			if err != nil {
				return nil, fmt.Errorf("failed to fetch info: %w", err)
			}

			mimeType := ""
			fileExtension := path.Ext(strings.ToLower(entry.Name()))
			switch fileExtension {
			case ".jpg", ".jpeg":
				mimeType = "image/jpeg"
			case ".png":
				mimeType = "image/png"
			}

			result = append(result, FileInfo{
				Filename: entry.Name(),
				Size:     entryInfo.Size(),
				Changed:  fileInfo.ChangeTime(),
				MimeType: mimeType,
			})
		}
	}

	// Sort by change date
	sort.Slice(result, func(i, j int) bool {
		return result[i].Changed.After(result[j].Changed)
	})

	return result, nil
}

func (p *PhotoRepo) hasThumbnails(fileName string) bool {
	if p.isThumbnail(fileName) {
		return true
	}

	smallThumbnail := fmt.Sprintf("files/%s", p.getSmallThumbnailFilename(fileName))
	mediumThumbnail := fmt.Sprintf("files/%s", p.getMediumThumbnailFilename(fileName))

	if _, err := os.Stat(smallThumbnail); err != nil {
		return false
	}

	if _, err := os.Stat(mediumThumbnail); err != nil {
		return false
	}

	return true
}

func (p *PhotoRepo) generateThumbnails(fileName string, settings *settings.AppSettings) error {
	log.Printf("INFO: Generating thumbnails for %s\n", fileName)
	smallThumbnail := p.getSmallThumbnailFilename(fileName)
	mediumThumbnail := p.getMediumThumbnailFilename(fileName)
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

func (p *PhotoRepo) getPhotoInfo(photoFullPath string) PhotoInfo {
	photoInfo := PhotoInfo{}

	photoExt := path.Ext(photoFullPath)
	infoFile := strings.Replace(photoFullPath, photoExt, ".yaml", -1)

	data, err := os.ReadFile(infoFile)
	if err != nil {
		log.Printf("WARNING: failed to read info file: %v", err)
		photoInfo.Title = "Untitled"
		return photoInfo
	}

	if err := yaml.Unmarshal(data, &photoInfo); err != nil {
		log.Printf("ERROR: failed to unmarshal info file %s: %v", infoFile, err)
		return photoInfo
	}

	if photoInfo.Title == "" {
		photoInfo.Title = "Untitled"
	}

	return photoInfo
}
