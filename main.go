package main

import (
	"errors"
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/saaste/portfolio/pkg/auth"
	"github.com/saaste/portfolio/pkg/handlers"
	"github.com/saaste/portfolio/pkg/photo"
	"github.com/saaste/portfolio/pkg/settings"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

func pollChanges(handler *handlers.Handler, appSettings *settings.AppSettings, photoRepo *photo.PhotoRepo) {
	tick := time.Tick(appSettings.RefreshInterval)
	for range tick {
		newSettings, err := settings.ReadSettings()
		if err != nil {
			log.Fatalf("failed to read app settings: %s", err)
		}

		photos, err := photoRepo.GetPhotos()
		if err != nil {
			log.Fatalf("failed to get photos: %s", err)
		}

		log.Println("INFO: Handler updated")
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

	photoRepo := photo.NewPhotoRepo(appSettings)

	if *generateThumbnails {
		photoRepo.GenerateThumbnail()
		os.Exit(0)
	}

	photos, err := photoRepo.GetPhotos()
	if err != nil {
		log.Fatalf("failed to get photos: %s", err)
	}

	handler := handlers.NewHandler(appSettings, photoRepo, photos)
	authMiddleware := auth.NewAuthMiddleware(appSettings)

	r := chi.NewRouter()
	r.Use(middleware.Logger)
	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)
	r.Use(middleware.Logger)
	r.Use(middleware.StripSlashes)
	r.Use(authMiddleware.Authenticate)

	r.Group(func(r chi.Router) {
		r.Use(authMiddleware.RequiresAuthentication)
		r.Get("/admin", handler.HandleAdmin)
		r.Get("/admin/settings", handler.HandleSettings)
		r.Post("/admin/settings", handler.HandleSettings)
	})

	r.Group(func(r chi.Router) {
		r.Get("/", handler.HandleRoot)
		r.Get("/about", handler.HandleAbout)
		r.Get("/feed", handler.HandleFeed)
		r.Get("/rss.xml", handler.HandleFeed)
		r.Get("/atom.xml", handler.HandleFeed)
		r.Get("/feed.json", handler.HandleFeed)
		r.Get("/photos/{filename}", handler.HandlePhoto)
		r.Get("/login", handler.HandleLogin)
		r.Post("/login", handler.HandleLogin)
		r.Get("/logout", handler.LogoutHandler)
	})

	handler.FileServer(r, "/photo", http.Dir("files"))
	handler.FileServer(r, "/static", http.Dir(fmt.Sprintf("ui/%s/static", appSettings.Theme)))

	go pollChanges(handler, appSettings, photoRepo)
	log.Printf("INFO: Photo polling initialized with %s", appSettings.RefreshInterval)

	err = http.ListenAndServe(fmt.Sprintf(":%s", *port), r)
	if errors.Is(err, http.ErrServerClosed) {
		log.Println("INFO: Server closed")
	} else if err != nil {
		log.Fatalf("Error starting server: %v\n", err)
	}
}
