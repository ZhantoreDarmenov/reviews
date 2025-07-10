package services

import (
	_ "bytes"
	"context"
	_ "encoding/json"
	"errors"
	"fmt"
	"github.com/dgrijalva/jwt-go"
	"github.com/google/uuid"
	_ "github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
	_ "io"
	_ "io/ioutil"
	"log"
	_ "math/rand"
	_ "net/http"
	_ "net/url"
	_ "os"
	"reviews/internal/models"
	"reviews/internal/repositories"
	"reviews/utils"
	"strconv"
	_ "strings"
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

func (s *UserService) SignIn(ctx context.Context, email, password string) (models.Tokens, error) {
	user, err := s.UserRepo.GetUserByEmail(ctx, email)
	if err != nil {
		log.Printf("User not found: %s", email)
		return models.Tokens{}, errors.New("user not found")
	}

	// Compare the provided password with the hashed password
	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password)); err != nil {
		log.Printf("Invalid password for user: %s", email)
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

func (s *UserService) CreateUser(ctx context.Context, user models.User) (models.User, error) {
	return s.UserRepo.CreateUser(ctx, user)
}

func (s *UserService) SignUp(ctx context.Context, user models.User, inputCode string) (models.SignUpResponse, error) {
	// 1. Получаем ожидаемый код из базы
	codeFromDB, err := s.UserRepo.GetVerificationCodeByEmail(ctx, user.Email)
	if err != nil {
		return models.SignUpResponse{}, err
	}

	// 2. Сравниваем коды
	if inputCode != codeFromDB {
		return models.SignUpResponse{}, models.ErrInvalidVerificationCode
	}

	// 3. Хешируем пароль
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		return models.SignUpResponse{}, err
	}
	user.Password = string(hashedPassword)
	user.Role = "client"

	// 4. Сохраняем пользователя
	newUser, err := s.UserRepo.CreateUser(ctx, user)
	if err != nil {
		return models.SignUpResponse{}, err
	}

	// 5. Можно очистить использованный код, если хочешь
	_ = s.UserRepo.ClearVerificationCode(ctx, user.Email)

	return models.SignUpResponse{User: newUser}, nil
}

func (s *UserService) UpgradeToTrainer(ctx context.Context, userID int) error {
	return s.UserRepo.UpdateUserRole(ctx, userID, "trainer")
}

// UpdateProfile updates user's profile. Changing email or password requires verification.
func (s *UserService) UpdateProfile(ctx context.Context, userID int, req models.UserUpdateRequest) (models.User, error) {
	user, err := s.UserRepo.GetUserByID(ctx, userID)
	if err != nil {
		return models.User{}, err
	}

	if req.Name != "" {
		user.Name = req.Name
	}
	if req.Phone != "" {
		user.Phone = req.Phone
	}

	if req.Email != "" || req.Password != "" {
		if req.VerificationCode == "" {
			return models.User{}, models.ErrInvalidVerificationCode
		}
		emailToCheck := req.Email
		if emailToCheck == "" {
			emailToCheck = user.Email
		}
		code, err := s.UserRepo.GetVerificationCodeByEmail(ctx, emailToCheck)
		if err != nil {
			return models.User{}, err
		}
		if code != req.VerificationCode {
			return models.User{}, models.ErrInvalidVerificationCode
		}
		_ = s.UserRepo.ClearVerificationCode(ctx, emailToCheck)
		if req.Email != "" {
			user.Email = req.Email
		}
		if req.Password != "" {
			hashed, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
			if err != nil {
				return models.User{}, err
			}
			user.Password = string(hashed)
		}
	}

	return s.UserRepo.UpdateUser(ctx, user)
}
