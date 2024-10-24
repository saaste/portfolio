package handlers

import (
	"log"
	"net/http"

	"github.com/go-chi/chi/v5"
)

func (h *Handler) HandlePhoto(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	filenameParam := chi.URLParam(r, "filename")
	photoResult, err := h.photoRepo.GetPhoto(filenameParam)
	if err != nil {
		log.Printf("ERROR: Failed to parse about templates: %s", err) // TODO: Error handler
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	} else if photoResult.Current == nil {
		http.Error(w, "Not Found", http.StatusNotFound)
		return
	}

	data := h.getTemplateData(ctx)
	data.Photo = photoResult

	h.parseTemplate(w, "photo.html", data)
}
