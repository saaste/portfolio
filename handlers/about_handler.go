package handlers

import "net/http"

func (h *Handler) HandleAbout(w http.ResponseWriter, r *http.Request) {
	h.parseTemplate(w, "about.html", h.getTemplateData())
}
