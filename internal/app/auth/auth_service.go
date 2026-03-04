package auth

import (
	"encoding/base64"
	"errors"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

var ErrInvalidCredentials = errors.New("tên đăng nhập hoặc mật khẩu không đúng")

type AuthService interface {
	Login(req LoginRequest) (*LoginResponse, error)
}

type authService struct {
	adminUsername     string
	adminPasswordHash string // stored as base64 to avoid $ expansion in .env
	jwtSecret         []byte
}

func NewAuthService(username, passwordHashB64, jwtSecret string) AuthService {
	return &authService{
		adminUsername:     username,
		adminPasswordHash: passwordHashB64,
		jwtSecret:         []byte(jwtSecret),
	}
}

func (s *authService) Login(req LoginRequest) (*LoginResponse, error) {
	if req.Username != s.adminUsername {
		return nil, ErrInvalidCredentials
	}

	// Decode base64 → raw bcrypt hash (workaround for Viper $ expansion)
	hashBytes, err := base64.StdEncoding.DecodeString(s.adminPasswordHash)
	if err != nil {
		return nil, ErrInvalidCredentials
	}

	if err := bcrypt.CompareHashAndPassword(hashBytes, []byte(req.Password)); err != nil {
		return nil, ErrInvalidCredentials
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"sub": req.Username,
		"exp": time.Now().Add(24 * time.Hour).Unix(),
		"iat": time.Now().Unix(),
	})

	tokenStr, err := token.SignedString(s.jwtSecret)
	if err != nil {
		return nil, err
	}

	return &LoginResponse{Token: tokenStr}, nil
}
