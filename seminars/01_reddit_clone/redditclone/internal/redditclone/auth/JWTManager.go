package auth

import (
	"redditclone/internal/redditclone/models"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

type JWTManager struct {
	Secret []byte
}

type JWTResponse struct {
	Token string `json:"token"`
}

func NewJWTManager(secret []byte) *JWTManager {
	return &JWTManager{Secret: secret}
}

func (j *JWTManager) Generate(u *models.User) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"user": map[string]string{
			"id":       u.Id,
			"username": u.Username,
		},
		"exp": time.Now().Add(24 * time.Hour).Unix(),
	})
	return token.SignedString(j.Secret)
}
