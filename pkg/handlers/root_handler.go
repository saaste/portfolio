package handlers

import "net/http"

func (h *Handler) HandleRoot(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	h.parseTemplate(w, "album.html", h.getTemplateData(ctx))
}
