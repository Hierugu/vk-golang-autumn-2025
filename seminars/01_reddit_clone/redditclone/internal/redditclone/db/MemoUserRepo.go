package db

import (
	"fmt"
	"redditclone/internal/redditclone/models"
	"sync"

	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

type MemoryUserRepository struct {
	mu    sync.RWMutex
	users map[string]*models.User
}

func CreateMemoryUserRepository() *MemoryUserRepository {
	return &MemoryUserRepository{sync.RWMutex{}, make(map[string]*models.User)}
}

func (r *MemoryUserRepository) Register(username, password string) (*models.User, error) {
	r.mu.Lock()
	defer r.mu.Unlock()

	if _, exists := r.users[username]; exists {
		return nil, fmt.Errorf("User %s already exists", username)
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}

	u := &models.User{
		Id:           uuid.NewString(),
		Username:     username,
		PasswordHash: string(hashedPassword),
	}
	r.users[username] = u
	return u, nil
}

func (r *MemoryUserRepository) Login(username, password string) (*models.User, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	u, exists := r.users[username]
	if !exists {
		return nil, fmt.Errorf("User %s does not exists", username)
	}

	if err := bcrypt.CompareHashAndPassword([]byte(u.PasswordHash), []byte(password)); err != nil {
		return nil, fmt.Errorf("User %s wrong password", username)
	}

	return u, nil
}

func (r *MemoryUserRepository) GetAll() ([]models.User, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	users := make([]models.User, 0, len(r.users))
	for _, u := range r.users {
		users = append(users, *u)
	}
	return users, nil
}
