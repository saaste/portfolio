package handlers

import (
	"html/template"
	"net/http"
	"strconv"

	"github.com/saaste/portfolio/pkg/settings"
)

func (h *Handler) HandleSettings(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	data := h.getTemplateData(ctx)

	if r.URL.Query().Get("saved") == "true" {
		data.Notification = "Changes saved"
	}

	if r.Method == http.MethodPost {
		err := r.ParseForm()
		if err != nil {
			h.internalServerError(w, err, "failed to parse settings form")
			return
		}

		data.Title = r.Form.Get("title")
		if data.Title == "" {
			data.Errors.Title = "required"
		}

		data.Description = r.Form.Get("description")
		if data.Description == "" {
			data.Errors.Description = "required"
		}

		data.Author = r.Form.Get("author")
		if data.Author == "" {
			data.Errors.Author = "required"
		}

		data.AboutTitle = r.Form.Get("about_title")
		about := r.Form.Get("about")
		data.About = template.HTML(about)
		if data.AboutTitle != "" && about == "" {
			data.Errors.About = "required when about title is set"
		}

		smallSizeValue := r.Form.Get("small_size")
		smallSize, err := strconv.ParseInt(smallSizeValue, 10, 32)
		if err != nil {
			data.Errors.SmallSize = "invalid value"
		} else if smallSize < 100 || smallSize > 1500 {
			data.Errors.SmallSize = "must be between 100 and 1500"
		}
		data.SmallSize = int(smallSize)

		mediumSizeValue := r.Form.Get("medium_size")
		mediumSize, err := strconv.ParseInt(mediumSizeValue, 10, 32)
		if err != nil {
			data.Errors.MediumSize = "invalid value"
		} else if mediumSize < 400 || mediumSize > 2500 {
			data.Errors.MediumSize = "must be between 400 and 2500"
		} else if mediumSize < smallSize {
			data.Errors.MediumSize = "must be larger than small image size"
		}
		data.MediumSize = int(mediumSize)

		if !data.Errors.HasErrors() {
			h.appSettings.Title = data.Title
			h.appSettings.Description = data.Description
			h.appSettings.Author = data.Author
			h.appSettings.AboutTitle = data.AboutTitle
			h.appSettings.About = about
			h.appSettings.SmallSize = int(smallSize)
			h.appSettings.MediumSize = int(mediumSize)

			err = settings.SaveSetting(h.appSettings)
			if err != nil {
				h.internalServerError(w, err, "failed to save settings")
				return
			}

			http.Redirect(w, r, "/admin/settings?saved=true", http.StatusFound)
			return
		}
	}

	h.parseTemplate(w, "settings.html", data)
}
