package services

import (
	"context"
	"errors"
	"fmt"
	"github.com/dgrijalva/jwt-go"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
	"log"
	"reviews/internal/models"
	"reviews/internal/repositories"
	"reviews/utils"
	"strconv"
	"time"
)

type tokenClaims struct {
	jwt.StandardClaims
	UserID int    `json:"user_id"`
	Role   string `json:"role"`
}
type UserService struct {
	UserRepo     *repositories.UserRepository
	TokenManager *utils.Manager
}

// SignUp creates a new user and returns auth tokens for the newly created account.
func (s *UserService) SignUp(ctx context.Context, user models.User) (models.Tokens, error) {
	hashed, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		return models.Tokens{}, err
	}
	user.Password = string(hashed)
	user.Role = "user"

	created, err := s.UserRepo.CreateUser(ctx, user)
	if err != nil {
		return models.Tokens{}, err
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, &tokenClaims{
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Add(tokenTTL).Unix(),
			IssuedAt:  time.Now().Unix(),
		},
		UserID: created.ID,
		Role:   created.Role,
	})

	accessToken, err := token.SignedString([]byte(signingKey))
	if err != nil {
		return models.Tokens{}, err
	}

	return s.CreateSession(ctx, created, accessToken)
}

func (s *UserService) SignIn(ctx context.Context, login, password string) (models.Tokens, error) {
	user, err := s.UserRepo.GetUserByLogin(ctx, login)
	if err != nil {
		log.Printf("User not found: %s", login)
		return models.Tokens{}, errors.New("user not found")
	}

	// Compare the provided password with the hashed password
	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password)); err != nil {
		log.Printf("Invalid password for user: %s", login)
		return models.Tokens{}, errors.New("invalid password")
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, &tokenClaims{
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Add(tokenTTL).Unix(),
			IssuedAt:  time.Now().Unix(),
		},
		UserID: user.ID,
		Role:   user.Role,
	})

	accessToken, err := token.SignedString([]byte(signingKey))
	if err != nil {
		log.Printf("Error signing token: %v", err)
		return models.Tokens{}, err
	}
	fmt.Println("login token:", accessToken)
	tokens, err := s.CreateSession(ctx, user, accessToken)
	if err != nil {
		log.Printf("Error creating session: %v", err)
		return models.Tokens{}, err
	}

	return tokens, nil
}

const (
	salt       = "sadasdnsadna"
	tokenTTL   = 120 * time.Minute
	signingKey = "asdadsadadaadsasd"
)

func (s *UserService) CreateSession(ctx context.Context, user models.User, accessToken string) (models.Tokens, error) {
	var (
		res models.Tokens
		err error
	)

	userIDStr := strconv.Itoa(user.ID)

	res.AccessToken = accessToken

	// Generate RefreshToken using UUID as a fallback
	res.RefreshToken = uuid.New().String() // Fallback if TokenManager is unavailable
	if s.TokenManager != nil {
		res.RefreshToken, err = s.TokenManager.NewRefreshToken()
		if err != nil {
			return res, err
		}
	}

	// Создание и сохранение сессии с RefreshToken
	session := models.Session{
		RefreshToken: res.RefreshToken,
		ExpiresAt:    time.Now().Add(24 * 30 * 2 * time.Hour),
	}

	err = s.UserRepo.SetSession(ctx, userIDStr, session)
	if err != nil {
		return res, err
	}

	return res, nil
}
