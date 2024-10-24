package photo

import (
	"fmt"
	"log"
	"os"

	"github.com/saaste/portfolio/pkg/settings"
)

type PhotoRepo struct {
	settings *settings.AppSettings
}

func NewPhotoRepo(settings *settings.AppSettings) *PhotoRepo {
	return &PhotoRepo{
		settings: settings,
	}
}

func (p *PhotoRepo) GetPhotos() ([]Photo, error) {
	photos := make([]Photo, 0)

	entries, err := os.ReadDir("files/")
	if err != nil {
		return photos, fmt.Errorf("failed to read files directory: %w", err)
	}

	originalPhotos, err := p.getOriginalPhotos(entries)
	if err != nil {
		return photos, fmt.Errorf("failed to get original photos: %w", err)
	}

	for _, entry := range originalPhotos {
		fullPath := fmt.Sprintf("files/%s", entry.Filename)

		hasThumbnals := p.hasThumbnails(entry.Filename)
		if !hasThumbnals {
			p.generateThumbnails(entry.Filename, p.settings)
		}
		photoInfo := p.getPhotoInfo(fullPath)
		photos = append(photos, Photo{
			FullFileName:   entry.Filename,
			MediumFileName: p.getMediumThumbnailFilename(entry.Filename),
			SmallFileName:  p.getSmallThumbnailFilename(entry.Filename),
			Changed:        entry.Changed,
			Size:           entry.Size,
			PhotoInfo:      photoInfo,
			MimeType:       entry.MimeType,
		})
	}

	return photos, nil
}

func (p *PhotoRepo) GetPhoto(fileName string) (*PhotoResult, error) {
	result := &PhotoResult{}

	entries, err := os.ReadDir("files/")
	if err != nil {
		return nil, fmt.Errorf("failed to read files directory: %w", err)
	}

	originalPhotos, err := p.getOriginalPhotos(entries)
	if err != nil {
		return nil, fmt.Errorf("failed to get original photos: %w", err)
	}

	for _, entry := range originalPhotos {
		fullPath := fmt.Sprintf("files/%s", entry.Filename)

		hasThumbnals := p.hasThumbnails(entry.Filename)
		if !hasThumbnals {
			p.generateThumbnails(entry.Filename, p.settings)
		}
		photoInfo := p.getPhotoInfo(fullPath)
		photo := &Photo{
			FullFileName:   entry.Filename,
			MediumFileName: p.getMediumThumbnailFilename(entry.Filename),
			SmallFileName:  p.getSmallThumbnailFilename(entry.Filename),
			Changed:        entry.Changed,
			Size:           entry.Size,
			PhotoInfo:      photoInfo,
			MimeType:       entry.MimeType,
		}

		if entry.Filename == fileName {
			result.Current = photo
		} else if result.Current != nil {
			result.Next = photo
			break
		} else {
			result.Previous = photo
		}
	}

	return result, nil
}

func (p *PhotoRepo) GenerateThumbnail() error {
	log.Printf("INFO: Starting thumbnail generation\n")
	entries, err := os.ReadDir("files/")
	if err != nil {
		return fmt.Errorf("failed to read files directory: %w", err)
	}

	originalPhotos, err := p.getOriginalPhotos(entries)
	if err != nil {
		return fmt.Errorf("failed to get original photos: %w", err)
	}

	for _, entry := range originalPhotos {
		p.generateThumbnails(entry.Filename, p.settings)
	}

	log.Printf("INFO: Thumbnails generated for %d photos\n", len(originalPhotos))
	return nil
}
