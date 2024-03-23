package handlers

import (
	"fmt"
	"html"
	"io"
	"log"
	"net/http"
	"time"

	"github.com/gorilla/feeds"
)

func (h *Handler) HandleFeed(w http.ResponseWriter, r *http.Request) {
	baseUrl := fmt.Sprintf("https://%s", r.Host)
	now := time.Now()
	feed := &feeds.Feed{
		Title:       h.appSettings.Title,
		Link:        &feeds.Link{Href: fmt.Sprintf("%s/feed", baseUrl), Rel: "self"},
		Description: h.appSettings.Title,
		Author:      &feeds.Author{Name: h.appSettings.Author},
		Created:     now,
	}

	templateData := h.getTemplateData()
	for _, photo := range templateData.Photos {
		feed.Items = append(feed.Items,
			&feeds.Item{
				Title:       "Photo",
				Link:        &feeds.Link{Href: fmt.Sprintf("%s/photo/%s", baseUrl, photo.FullFileName)},
				Author:      &feeds.Author{Name: h.appSettings.Author},
				Created:     photo.Changed,
				Description: html.EscapeString(fmt.Sprintf("<img src=\"%s/photo/%s\" />", baseUrl, photo.SmallFileName)),
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
