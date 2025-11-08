package db

import (
	"redditclone/internal/redditclone/models"
	"sync"
)

type MemoPostRepo struct {
	mu    sync.RWMutex
	posts map[string]*models.Post
}

func CreateMemoPostRepo() *MemoPostRepo {
	return &MemoPostRepo{sync.RWMutex{}, make(map[string]*models.Post)}
}
