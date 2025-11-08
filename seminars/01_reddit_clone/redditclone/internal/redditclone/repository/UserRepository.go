package repository

import "redditclone/internal/redditclone/models"

type UserRepository interface {
	GetAll() ([]models.User, error)
	Register(username, password string) (*models.User, error)
	Login(username, password string) (*models.User, error)
}
