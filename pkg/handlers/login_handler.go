package handlers

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/saaste/portfolio/pkg/auth"
)

func (h *Handler) HandleLogin(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	redirectUrl := cleanRedirectURL(r.URL.Query().Get("return"))

	data := h.getTemplateData(ctx)
	data.RedirectURL = redirectUrl

	if r.Method == http.MethodPost {
		err := r.ParseForm()
		if err != nil {
			h.internalServerError(w, err, "failed to parse login form")
			return
		}

		password := r.Form.Get("password")
		data.RedirectURL = cleanRedirectURL(r.Form.Get("redirect_url"))

		if password == h.appSettings.Password {
			jwtToken, err := h.jwtParser.CreateJWT()
			if err != nil {
				h.internalServerError(w, err, "failed to create JWT token")
				return
			}
			http.SetCookie(w, &http.Cookie{
				Name:     auth.AuthCookieKey,
				Value:    jwtToken,
				Path:     "/",
				HttpOnly: true,
				Secure:   false,
				SameSite: http.SameSiteLaxMode,
				MaxAge:   60 * 60 * 24 * 7,
			})
			fmt.Printf("Redirecting to url: %s\n", data.RedirectURL)
			http.Redirect(w, r, data.RedirectURL, http.StatusFound)
			return
		}
		data.Errors.General = "Invalid password"
	}

	h.parseTemplate(w, "login.html", data)
}

func cleanRedirectURL(url string) string {
	if url == "" || !strings.HasPrefix(url, "/") {
		return "/"
	}
	return url
}
