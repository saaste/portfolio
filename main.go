package main

import (
	"errors"
	"flag"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/saaste/portfolio/photo"
	"github.com/saaste/portfolio/settings"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"

	"html/template"

	"github.com/gorilla/feeds"
)

var photos []photo.Photo = make([]photo.Photo, 0)
var appSettings *settings.AppSettings

type TemplateData struct {
	Photos     []photo.Photo
	Year       int
	Title      string
	Author     string
	About      template.HTML
	SmallSize  int
	MediumSize int
}

func getTemplateData() TemplateData {
	return TemplateData{
		Photos:     photos,
		Year:       time.Now().Year(),
		Title:      appSettings.Title,
		Author:     appSettings.Author,
		About:      template.HTML(appSettings.About),
		SmallSize:  appSettings.SmallSize,
		MediumSize: appSettings.MediumSize,
	}
}

func handleRoot(w http.ResponseWriter, r *http.Request) {
	t, err := template.ParseFiles("ui/templates/base.html", "ui/templates/album.html")
	if err != nil {
		log.Printf("Failed to parse album templates: %s", err)
		http.Error(w, "Internal Server Error", 500)
	}
	err = t.ExecuteTemplate(w, "base", getTemplateData())
	if err != nil {
		log.Printf("Failed to execute album template: %s", err)
		http.Error(w, "Internal Server Error", 500)

	}
}

func handleAbout(w http.ResponseWriter, r *http.Request) {
	t, err := template.ParseFiles("ui/templates/base.html", "ui/templates/about.html")
	if err != nil {
		log.Printf("Failed to parse about templates: %s", err)
		http.Error(w, "Internal Server Error", 500)
	}
	err = t.ExecuteTemplate(w, "base", getTemplateData())
	if err != nil {
		log.Printf("Failed to execute about template: %s", err)
		http.Error(w, "Internal Server Error", 500)

	}
}

func handleFeed(w http.ResponseWriter, r *http.Request) {
	baseUrl := fmt.Sprintf("https://%s", r.Host)
	now := time.Now()
	feed := &feeds.Feed{
		Title:       appSettings.Title,
		Link:        &feeds.Link{Href: fmt.Sprintf("%s/feed", baseUrl)},
		Description: appSettings.Title,
		Author:      &feeds.Author{Name: appSettings.Author},
		Created:     now,
	}

	for _, photo := range photos {
		feed.Items = append(feed.Items,
			&feeds.Item{
				Title:   "Photo",
				Link:    &feeds.Link{Href: fmt.Sprintf("%s/photo/%s", baseUrl, photo.FullFileName)},
				Author:  &feeds.Author{Name: appSettings.Author},
				Created: photo.Changed,
			})
	}

	atom, err := feed.ToAtom()
	if err != nil {
		log.Printf("Failed to create feed: %s", err)
		http.Error(w, "Internal Server Error", 500)
	}

	w.Header().Set("Content-Type", "application/xml")
	io.WriteString(w, atom)
}

func fileServer(r chi.Router, path string, root http.FileSystem) {
	if strings.ContainsAny(path, "{}*") {
		panic("FileServer does not permit any URL parameters.")
	}

	if path != "/" && path[len(path)-1] != '/' {
		r.Get(path, http.RedirectHandler(path+"/", http.StatusMovedPermanently).ServeHTTP)
		path += "/"
	}
	path += "*"

	r.Get(path, func(w http.ResponseWriter, r *http.Request) {
		rctx := chi.RouteContext(r.Context())
		pathPrefix := strings.TrimSuffix(rctx.RoutePattern(), "/*")
		fs := http.StripPrefix(pathPrefix, http.FileServer(root))
		fs.ServeHTTP(w, r)
	})
}

func pollPhotos(appSettings *settings.AppSettings) {
	tick := time.Tick(appSettings.RefreshInterval)
	for range tick {
		p, err := photo.GetPhotos(appSettings, false)

		if err != nil {
			log.Fatalf("failed to get photos: %s", err)
		}

		photos = p
	}

}

func main() {
	generateThumbnails := flag.Bool("generate-thumbnails", false, "Generate all thumbnails")
	flag.Parse()

	s, err := settings.ReadSettings()
	if err != nil {
		log.Fatalf("failed to read app settings: %s", err)
	}
	appSettings = s

	if *generateThumbnails {
		photo.GetPhotos(appSettings, true)
		os.Exit(0)
	}

	log.Printf("Photo polling initialized with %s", appSettings.RefreshInterval)

	p, err := photo.GetPhotos(appSettings, false)
	if err != nil {
		log.Fatalf("failed to get photos: %s", err)
	}
	photos = p

	r := chi.NewRouter()
	r.Use(middleware.Logger)
	r.Get("/", handleRoot)
	r.Get("/about", handleAbout)
	r.Get("/feed", handleFeed)

	fileServer(r, "/photo", http.Dir("files"))
	fileServer(r, "/static", http.Dir("ui/static"))

	go pollPhotos(appSettings)

	err = http.ListenAndServe(":8000", r)
	if errors.Is(err, http.ErrServerClosed) {
		log.Println("Server closed")
	} else if err != nil {
		log.Fatalf("Error starting server: %v\n", err)
	}
}
