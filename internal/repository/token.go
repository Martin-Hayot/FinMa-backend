package repository

import (
	"errors"
	"sync"
	"time"
)

var (
	ErrTokenNotFound = errors.New("token not found")
)

type Token struct {
	AccessToken    string    `json:"access_token"`
	RefreshToken   string    `json:"refresh_token"`
	AccessExpires  time.Time `json:"access_expires"`
	RefreshExpires time.Time `json:"refresh_expires"`
	CreatedAt      time.Time `json:"created_at"`
}

type TokenRepository interface {
	GetToken() (*Token, error)
	SaveToken(token *Token) error
	HasToken() bool
	ClearToken()
}

type tokenRepository struct {
	mu    sync.RWMutex
	token *Token
}

func NewTokenRepository() TokenRepository {
	return &tokenRepository{}
}

func (r *tokenRepository) GetToken() (*Token, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	if r.token == nil {
		return nil, ErrTokenNotFound
	}

	// Return a copy to prevent external modification
	tokenCopy := *r.token
	return &tokenCopy, nil
}

func (r *tokenRepository) SaveToken(token *Token) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	// Store a copy with current timestamp
	tokenCopy := *token
	tokenCopy.CreatedAt = time.Now()
	r.token = &tokenCopy

	return nil
}

func (r *tokenRepository) HasToken() bool {
	r.mu.RLock()
	defer r.mu.RUnlock()
	return r.token != nil
}

func (r *tokenRepository) ClearToken() {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.token = nil
}
