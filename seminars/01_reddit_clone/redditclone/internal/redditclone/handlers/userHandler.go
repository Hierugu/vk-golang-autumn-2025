package handlers

import (
	"redditclone/internal/redditclone/auth"
	"redditclone/internal/redditclone/repository"
)

type UserHandler struct {
	Repo       repository.UserRepository
	jwtManager *auth.JWTManager
}

func NewUserHandler(repo repository.UserRepository, jwtSecret string) *UserHandler {
	return &UserHandler{Repo: repo, jwtManager: auth.NewJWTManager([]byte(jwtSecret))}
}
