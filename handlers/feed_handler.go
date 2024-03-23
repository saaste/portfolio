package handlers

import (
	"encoding/json"
	"fmt"
	"html"
	"io"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/saaste/portfolio/photo"
)

const (
	generatorString = "Portfolio (https://github.com/saaste/portfolio)"
)

type jsonFeed struct {
	Version     string       `json:"version"`
	Title       string       `json:"title"`
	Description string       `json:"description"`
	HomePageUrl string       `json:"home_page_url"`
	FeedUrl     string       `json:"feed_url,omitempty"`
	Authors     []jsonAuthor `json:"authors,omitempty"`
	Items       []jsonItem   `json:"items"`
}

type jsonAuthor struct {
	Name string `json:"name"`
}

type jsonItem struct {
	Id            string   `json:"id"`
	Title         string   `json:"title"`
	ContentHtml   string   `json:"content_html"`
	Url           string   `json:"url"`
	DatePublished string   `json:"date_published"`
	Tags          []string `json:"tags,omitempty"`
}

func (h *Handler) HandleFeed(w http.ResponseWriter, r *http.Request) {
	switch {
	case strings.HasSuffix(r.RequestURI, "rss.xml"):
		h.toRSS(w)
	case strings.HasSuffix(r.RequestURI, "atom.xml"):
		h.toAtom(w)
	case strings.HasSuffix(r.RequestURI, "feed.json"):
		h.toJSON(w)
	}
}

func (h *Handler) toRSS(w http.ResponseWriter) {
	pubDate := time.Now()
	if len(h.photos) > 0 {
		pubDate = h.photos[0].Changed
	}

	builder := make([]string, 0)
	builder = append(builder, `<?xml version="1.0"?>`)
	builder = append(builder, `<rss version="2.0" xmlns:atom="http://www.w3.org/2005/Atom">`)
	builder = append(builder, "\t<channel>")
	builder = append(builder, fmt.Sprintf("\t\t<title>%s</title>", html.EscapeString(h.appSettings.Title)))
	builder = append(builder, fmt.Sprintf("\t\t<link>%s/</link>", html.EscapeString(h.appSettings.BaseURL)))
	builder = append(builder, fmt.Sprintf("\t\t<description>%s</description>", html.EscapeString(h.appSettings.Description)))
	builder = append(builder, fmt.Sprintf("\t\t<pubDate>%s</pubDate>", pubDate.Format(time.RFC1123Z)))
	builder = append(builder, fmt.Sprintf("\t\t<lastBuildDate>%s</lastBuildDate>", pubDate.Format(time.RFC1123Z)))
	builder = append(builder, fmt.Sprintf("\t\t<generator>%s</generator>", generatorString))
	builder = append(builder, fmt.Sprintf("\t\t<copyright>Copyright © %s %d</copyright>", html.EscapeString(h.appSettings.Author), time.Now().Year()))
	builder = append(builder, fmt.Sprintf(`%s<atom:link href="%s/rss.xml" rel="self" type="application/rss+xml"></atom:link>`, "\t\t", h.appSettings.BaseURL))

	for _, photo := range h.photos {
		builder = append(builder, "\t\t<item>")

		if photo.PhotoInfo.Title != "" {
			builder = append(builder, fmt.Sprintf("\t\t\t<title>%s</title>", photo.PhotoInfo.Title))
		} else {
			builder = append(builder, "\t\t\t<title>Photo without title</title>")
		}

		builder = append(builder, fmt.Sprintf("\t\t\t<link>%s/photo/%s</link>", h.appSettings.BaseURL, photo.FullFileName))

		builder = append(builder, fmt.Sprintf("\t\t\t<description>%s</description>", h.buildDescription(photo)))
		builder = append(builder, fmt.Sprintf("\t\t\t<pubDate>%s</pubDate>", photo.Changed.Format(time.RFC1123Z)))
		builder = append(builder, fmt.Sprintf("\t\t\t<guid>%s/photo/%s</guid>", h.appSettings.BaseURL, photo.FullFileName))

		builder = append(builder, "\t\t</item>")
	}

	builder = append(builder, "\t</channel>")
	builder = append(builder, "</rss>")

	content := strings.Join(builder, "\n")
	w.Header().Set("Content-Type", "application/xml")
	io.WriteString(w, content)
}

func (h *Handler) toAtom(w http.ResponseWriter) {
	pubDate := time.Now()
	if len(h.photos) > 0 {
		pubDate = h.photos[0].Changed
	}

	builder := make([]string, 0)

	builder = append(builder, `<?xml version="1.0" encoding="utf-8"?>`)
	builder = append(builder, `<feed xmlns="http://www.w3.org/2005/Atom">`)
	builder = append(builder, fmt.Sprintf("\t<title>%s</title>", html.EscapeString(h.appSettings.Title)))
	builder = append(builder, fmt.Sprintf("\t<subtitle>%s</subtitle>", html.EscapeString(h.appSettings.Description)))
	builder = append(builder, fmt.Sprintf(`%s<link href="%s" />`, "\t", h.appSettings.BaseURL))
	builder = append(builder, fmt.Sprintf(`%s<link href="%s/atom.xml" rel="self" />`, "\t", h.appSettings.BaseURL))
	builder = append(builder, fmt.Sprintf("\t<updated>%s</updated>", pubDate.Format(time.RFC3339)))

	builder = append(builder, "\t<author>")
	builder = append(builder, fmt.Sprintf("\t\t<name>%s</name>", html.EscapeString(h.appSettings.Author)))
	builder = append(builder, "\t</author>")
	builder = append(builder, fmt.Sprintf("\t<id>%s/atom.xml</id>", h.appSettings.BaseURL))
	builder = append(builder, fmt.Sprintf("\t<generator>%s</generator>", generatorString))

	for _, photo := range h.photos {
		builder = append(builder, "\t<entry>")
		builder = append(builder, fmt.Sprintf("\t\t<id>%s/photos/%s</id>", h.appSettings.BaseURL, photo.FullFileName))
		if photo.PhotoInfo.Title != "" {
			builder = append(builder, fmt.Sprintf("\t\t<title>%s</title>", html.EscapeString(photo.PhotoInfo.Title)))
		} else {
			builder = append(builder, "\t\t<title>Photo without title</title>")
		}

		builder = append(builder, fmt.Sprintf("\t\t<updated>%s</updated>", photo.Changed.Format(time.RFC3339)))
		builder = append(builder, fmt.Sprintf("\t\t<content type=\"html\">%s</content>", h.buildDescription(photo)))
		builder = append(builder, "\t</entry>")
	}

	builder = append(builder, "</feed>")
	content := strings.Join(builder, "\n")

	w.Header().Set("Content-Type", "application/xml")
	io.WriteString(w, content)
}

func (h *Handler) toJSON(w http.ResponseWriter) {
	feed := jsonFeed{
		Version:     "https://jsonfeed.org/version/1.1",
		Title:       h.appSettings.Title,
		Description: h.appSettings.Description,
		HomePageUrl: h.appSettings.BaseURL,
		FeedUrl:     fmt.Sprintf("%s/feed.json", h.appSettings.BaseURL),
	}

	if h.appSettings.Author != "" {
		feed.Authors = append(feed.Authors, jsonAuthor{Name: h.appSettings.Author})
	}

	for _, photo := range h.photos {
		title := "Photo without title"
		if photo.PhotoInfo.Title != "" {
			title = photo.PhotoInfo.Title
		}
		item := jsonItem{
			Id:            fmt.Sprintf("%s/%s", h.appSettings.BaseURL, photo.FullFileName),
			Title:         title,
			ContentHtml:   h.buildDescription(photo),
			Url:           fmt.Sprintf("%s/photos/%s", h.appSettings.BaseURL, photo.FullFileName),
			DatePublished: photo.Changed.Format(time.RFC3339),
		}
		feed.Items = append(feed.Items, item)
	}

	byt, err := json.MarshalIndent(feed, "", "  ")
	if err != nil {
		log.Printf("ERROR: Failed to marshal JSON feed: %v", err)
		http.Error(w, "Internal Server Error", 500)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	io.WriteString(w, string(byt))
}

func (h *Handler) buildDescription(photo photo.Photo) string {
	descBuilder := make([]string, 0)
	if photo.PhotoInfo.Title != "" {
		descBuilder = append(descBuilder, fmt.Sprintf("<p>%s</p>", photo.PhotoInfo.Title))
	}

	if photo.PhotoInfo.AltText != "" {
		descBuilder = append(descBuilder, fmt.Sprintf("<img src=\"%s/photos/%s\" alt=\"%s\" />", h.appSettings.BaseURL, photo.SmallFileName, photo.PhotoInfo.AltText))
	} else {
		descBuilder = append(descBuilder, fmt.Sprintf("<img src=\"%s/photos/%s\" />", h.appSettings.BaseURL, photo.SmallFileName))
	}

	if photo.PhotoInfo.Description != "" {
		descBuilder = append(descBuilder, fmt.Sprintf("<p>%s</p>", photo.PhotoInfo.Description))
	}

	return html.EscapeString(strings.Join(descBuilder, ""))
}
