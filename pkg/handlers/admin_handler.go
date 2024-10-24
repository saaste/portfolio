package handlers

import "net/http"

func (h *Handler) HandleAdmin(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	h.parseTemplate(w, "admin.html", h.getTemplateData(ctx))
}
