package handlers

import (
	"context"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"time"

	"github.com/saaste/portfolio/pkg/auth"
	"github.com/saaste/portfolio/pkg/photo"
	"github.com/saaste/portfolio/pkg/settings"
)

type TemplateData struct {
	Photo           *photo.PhotoResult
	Photos          []photo.Photo
	Year            int
	Title           string
	Description     string
	Author          string
	About           template.HTML
	AboutTitle      string
	SmallSize       int
	MediumSize      int
	BaseURL         string
	RedirectURL     string
	Errors          *Errors
	Notification    string
	IsAuthenticated bool
}

type Errors struct {
	General     string
	Title       string
	Description string
	Author      string
	About       string
	SmallSize   string
	MediumSize  string
}

func (e *Errors) HasErrors() bool {
	return e.General != "" || e.Title != "" || e.Description != "" || e.Author != "" || e.About != "" || e.SmallSize != "" || e.MediumSize != ""
}

type Handler struct {
	appSettings *settings.AppSettings
	photoRepo   *photo.PhotoRepo
	photos      []photo.Photo
	jwtParser   *auth.JwtParser
}

func NewHandler(appSettings *settings.AppSettings, photoRepo *photo.PhotoRepo, photos []photo.Photo) *Handler {
	return &Handler{
		appSettings: appSettings,
		photoRepo:   photoRepo,
		photos:      photos,
		jwtParser:   auth.NewJwtParser(appSettings),
	}
}

func (h *Handler) Update(appSettings *settings.AppSettings, photos []photo.Photo) {
	h.appSettings = appSettings
	h.photos = photos
}

func (h *Handler) parseTemplate(w http.ResponseWriter, templateFile string, data *TemplateData) {
	baseTemplate := fmt.Sprintf("ui/%s/templates/base.html", h.appSettings.Theme)
	targetTemplate := fmt.Sprintf("ui/%s/templates/%s", h.appSettings.Theme, templateFile)
	t, err := template.ParseFiles(baseTemplate, targetTemplate)
	if err != nil {
		log.Printf("ERROR: Failed to parse about templates: %s", err)
		http.Error(w, "Internal Server Error", 500)
	}
	err = t.ExecuteTemplate(w, "base", data)
	if err != nil {
		log.Printf("ERROR: Failed to execute about template: %s", err)
		http.Error(w, "Internal Server Error", 500)

	}
}

func (h *Handler) getTemplateData(ctx context.Context) *TemplateData {
	return &TemplateData{
		Photos:          h.photos,
		Year:            time.Now().Year(),
		Title:           h.appSettings.Title,
		Description:     h.appSettings.Description,
		Author:          h.appSettings.Author,
		About:           template.HTML(h.appSettings.About),
		AboutTitle:      h.appSettings.AboutTitle,
		SmallSize:       h.appSettings.SmallSize,
		MediumSize:      h.appSettings.MediumSize,
		BaseURL:         h.appSettings.BaseURL,
		RedirectURL:     "/",
		Errors:          &Errors{},
		Notification:    "",
		IsAuthenticated: ctx.Value(auth.AuthContextKeyIsAuthenticated) == true,
	}
}

func (h *Handler) internalServerError(w http.ResponseWriter, err error, message string) {
	fmt.Printf("ERROR: %s: %v\n", message, err)
	http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
}

func (h *Handler) notFound(w http.ResponseWriter) {
	http.Error(w, http.StatusText(http.StatusNotFound), http.StatusNotFound)
}
