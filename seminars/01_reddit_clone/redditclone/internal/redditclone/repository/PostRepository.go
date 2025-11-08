package repository

import "redditclone/internal/redditclone/models"

type PostRepository interface {
	GetAll() ([]models.Post, error)
	Create(post models.Post) error
}
