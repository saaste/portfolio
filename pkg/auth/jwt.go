package auth

import (
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/saaste/portfolio/pkg/settings"
)

type JwtParser struct {
	appSettings *settings.AppSettings
}

func NewJwtParser(appSettings *settings.AppSettings) *JwtParser {
	return &JwtParser{
		appSettings: appSettings,
	}
}

func (j *JwtParser) CreateJWT() (string, error) {
	now := time.Now().UTC()
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"iat": now.Unix(),
		"nbf": now.Unix(),
		"exp": now.Add(24 * time.Hour).Unix(),
	})

	tokenString, err := token.SignedString([]byte(j.appSettings.JWTSecret))
	return tokenString, err
}

func (j *JwtParser) ParseJWT(tokenString string) error {
	_, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}

		return []byte(j.appSettings.JWTSecret), nil
	})
	if err != nil {
		return fmt.Errorf("invalid token: %v", err)
	}

	return nil
}
