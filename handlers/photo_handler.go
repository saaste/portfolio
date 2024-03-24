package handlers

import (
	"fmt"
	"net/http"

	"github.com/go-chi/chi/v5"
)

func (h *Handler) HandlePhoto(w http.ResponseWriter, r *http.Request) {
	filenameParam := chi.URLParam(r, "filename")
	fmt.Println(filenameParam)
	h.parseTemplate(w, "photo.html")
}
