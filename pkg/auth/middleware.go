package auth

import (
	"context"
	"fmt"
	"net/http"
	"net/url"

	"github.com/saaste/portfolio/pkg/settings"
)

type AuthContextKey int

const (
	AuthContextKeyIsAuthenticated AuthContextKey = 0
)

const (
	AuthCookieKey string = "auth"
)

type AuthMiddleware struct {
	appSettings *settings.AppSettings
	jwtParser   *JwtParser
}

func NewAuthMiddleware(appConfig *settings.AppSettings) *AuthMiddleware {
	return &AuthMiddleware{
		appSettings: appConfig,
		jwtParser:   NewJwtParser(appConfig),
	}
}

func (a *AuthMiddleware) Authenticate(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		authCookie, err := r.Cookie(AuthCookieKey)
		if err != nil {
			next.ServeHTTP(w, r)
			return
		}

		err = a.jwtParser.ParseJWT(authCookie.Value)
		if err != nil {
			next.ServeHTTP(w, r)
			return
		}

		ctx := r.Context()
		ctx = context.WithValue(ctx, AuthContextKeyIsAuthenticated, true)

		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func (a *AuthMiddleware) RequiresAuthentication(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		isAuthenticated := ctx.Value(AuthContextKeyIsAuthenticated)
		if isAuthenticated == nil || isAuthenticated == false {
			uri := fmt.Sprintf("/login/?return=%s", url.QueryEscape(r.RequestURI))
			http.Redirect(w, r, uri, http.StatusFound)
			return
		}

		updatedToken, err := a.jwtParser.CreateJWT()
		if err != nil {
			http.Error(w, http.StatusText(http.StatusForbidden), http.StatusForbidden)
			return
		}

		http.SetCookie(w, &http.Cookie{
			Name:     AuthCookieKey,
			Value:    updatedToken,
			Path:     "/",
			HttpOnly: true,
			Secure:   false,
			SameSite: http.SameSiteLaxMode,
			MaxAge:   60 * 60 * 24 * 7,
		})

		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
