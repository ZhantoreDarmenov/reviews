package service

import (
	"fmt"
	"math/rand"
	"reviews/internal/models"
	"sync"
)

type AuthService struct {
	user     models.User
	sessions map[string]string
	mu       sync.Mutex
}

func NewAuthService(username, password string) *AuthService {
	return &AuthService{
		user:     models.User{Username: username, Password: password},
		sessions: make(map[string]string),
	}
}

func (s *AuthService) Login(username, password string) (string, bool) {
	if username != s.user.Username || password != s.user.Password {
		return "", false
	}
	token := fmt.Sprintf("%x", rand.Int63())
	s.mu.Lock()
	s.sessions[token] = s.user.Username
	s.mu.Unlock()
	return token, true
}

func (s *AuthService) Authenticate(token string) bool {
	s.mu.Lock()
	defer s.mu.Unlock()
	_, ok := s.sessions[token]
	return ok
}
