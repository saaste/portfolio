package main

import (
	"errors"
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/saaste/portfolio/handlers"
	"github.com/saaste/portfolio/photo"
	"github.com/saaste/portfolio/settings"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

func pollChanges(handler *handlers.Handler, appSettings *settings.AppSettings) {
	tick := time.Tick(appSettings.RefreshInterval)
	for range tick {
		newSettings, err := settings.ReadSettings()
		if err != nil {
			log.Fatalf("failed to read app settings: %s", err)
		}

		photos, err := photo.GetPhotos(newSettings, false)
		if err != nil {
			log.Fatalf("failed to get photos: %s", err)
		}

		log.Println("Handler updated")
		handler.Update(newSettings, photos)
	}

}

func main() {
	generateThumbnails := flag.Bool("generate-thumbnails", false, "Generate all thumbnails")
	port := flag.String("port", "8000", "Port")
	flag.Parse()

	appSettings, err := settings.ReadSettings()
	if err != nil {
		log.Fatalf("failed to read app settings: %s", err)
	}

	if *generateThumbnails {
		photo.GetPhotos(appSettings, true)
		os.Exit(0)
	}

	log.Printf("Photo polling initialized with %s", appSettings.RefreshInterval)

	photos, err := photo.GetPhotos(appSettings, false)
	if err != nil {
		log.Fatalf("failed to get photos: %s", err)
	}

	handler := handlers.NewHandler(appSettings, photos)

	r := chi.NewRouter()
	r.Use(middleware.Logger)
	r.Get("/", handler.HandleRoot)
	r.Get("/about", handler.HandleAbout)
	r.Get("/feed", handler.HandleFeed)
	r.Get("/rss.xml", handler.HandleFeed)
	r.Get("/atom.xml", handler.HandleFeed)
	r.Get("/feed.json", handler.HandleFeed)

	handler.FileServer(r, "/photo", http.Dir("files"))
	handler.FileServer(r, "/static", http.Dir("ui/static"))

	go pollChanges(handler, appSettings)

	err = http.ListenAndServe(fmt.Sprintf(":%s", *port), r)
	if errors.Is(err, http.ErrServerClosed) {
		log.Println("Server closed")
	} else if err != nil {
		log.Fatalf("Error starting server: %v\n", err)
	}
}
