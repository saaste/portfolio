package handlers

import (
	"fmt"
	"html/template"
	"log"
	"net/http"
	"time"

	"github.com/saaste/portfolio/photo"
	"github.com/saaste/portfolio/settings"
)

type TemplateData struct {
	Photos      []photo.Photo
	Year        int
	Title       string
	Description string
	Author      string
	About       template.HTML
	SmallSize   int
	MediumSize  int
	BaseURL     string
}

type Handler struct {
	appSettings *settings.AppSettings
	photos      []photo.Photo
}

func NewHandler(appSettings *settings.AppSettings, photos []photo.Photo) *Handler {
	return &Handler{
		appSettings: appSettings,
		photos:      photos,
	}
}

func (h *Handler) Update(appSettings *settings.AppSettings, photos []photo.Photo) {
	h.appSettings = appSettings
	h.photos = photos
}

func (h *Handler) parseTemplate(w http.ResponseWriter, templateFile string) {
	baseTemplate := fmt.Sprintf("ui/%s/templates/base.html", h.appSettings.Theme)
	targetTemplate := fmt.Sprintf("ui/%s/templates/%s", h.appSettings.Theme, templateFile)
	t, err := template.ParseFiles(baseTemplate, targetTemplate)
	if err != nil {
		log.Printf("ERROR: Failed to parse about templates: %s", err)
		http.Error(w, "Internal Server Error", 500)
	}
	err = t.ExecuteTemplate(w, "base", h.getTemplateData())
	if err != nil {
		log.Printf("ERROR: Failed to execute about template: %s", err)
		http.Error(w, "Internal Server Error", 500)

	}
}

func (h *Handler) getTemplateData() *TemplateData {
	return &TemplateData{
		Photos:      h.photos,
		Year:        time.Now().Year(),
		Title:       h.appSettings.Title,
		Description: h.appSettings.Description,
		Author:      h.appSettings.Author,
		About:       template.HTML(h.appSettings.About),
		SmallSize:   h.appSettings.SmallSize,
		MediumSize:  h.appSettings.MediumSize,
		BaseURL:     h.appSettings.BaseURL,
	}
}
