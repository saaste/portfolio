package handlers

import "net/http"

func (h *Handler) HandleAbout(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	h.parseTemplate(w, "about.html", h.getTemplateData(ctx))
}
