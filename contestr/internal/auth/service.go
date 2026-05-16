package auth

import (
	"contestr/internal/configs"
	"contestr/pkg/logger"
	"context"
	"crypto/subtle"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

const ContextUsernameKey = "admin_username"

var (
	ErrAdminDisabled       = errors.New("admin authentication is disabled")
	ErrInvalidCredentials  = errors.New("invalid credentials")
	ErrInvalidToken        = errors.New("invalid token")
)

type Claims struct {
	Username string `json:"sub"`
	jwt.RegisteredClaims
}

type Service struct {
	cfg *configs.Config
}

func NewService(ctx context.Context, cfg *configs.Config) *Service {
	if !cfg.Admin.Enabled() {
		logger.Warn(ctx, "admin authentication is disabled: set APP_ADMIN_JWT_SECRET, APP_ADMIN_USERNAME, and APP_ADMIN_PASSWORD")
	}
	return &Service{cfg: cfg}
}

func (s *Service) Enabled() bool {
	return s.cfg.Admin.Enabled()
}

func (s *Service) ValidateCredentials(username, password string) error {
	if !s.Enabled() {
		return ErrAdminDisabled
	}

	cfg := s.cfg.Admin
	if subtle.ConstantTimeCompare([]byte(username), []byte(cfg.Username)) != 1 {
		return ErrInvalidCredentials
	}
	if subtle.ConstantTimeCompare([]byte(password), []byte(cfg.Password)) != 1 {
		return ErrInvalidCredentials
	}
	return nil
}

func (s *Service) IssueToken(username string) (string, int64, error) {
	if !s.Enabled() {
		return "", 0, ErrAdminDisabled
	}

	ttl := s.cfg.Admin.JWTTTL
	if ttl <= 0 {
		ttl = 24 * time.Hour
	}

	now := time.Now()
	expiresAt := now.Add(ttl)
	claims := Claims{
		Username: username,
		RegisteredClaims: jwt.RegisteredClaims{
			Subject:   username,
			IssuedAt:  jwt.NewNumericDate(now),
			ExpiresAt: jwt.NewNumericDate(expiresAt),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	signed, err := token.SignedString([]byte(s.cfg.Admin.JWTSecret))
	if err != nil {
		return "", 0, fmt.Errorf("sign token: %w", err)
	}

	return signed, int64(ttl.Seconds()), nil
}

func (s *Service) ParseToken(tokenString string) (*Claims, error) {
	if !s.Enabled() {
		return nil, ErrAdminDisabled
	}

	tokenString = strings.TrimPrefix(tokenString, "Bearer ")
	token, err := jwt.ParseWithClaims(tokenString, &Claims{}, func(token *jwt.Token) (interface{}, error) {
		if token.Method != jwt.SigningMethodHS256 {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(s.cfg.Admin.JWTSecret), nil
	})
	if err != nil {
		return nil, ErrInvalidToken
	}

	claims, ok := token.Claims.(*Claims)
	if !ok || !token.Valid {
		return nil, ErrInvalidToken
	}
	if claims.Username == "" {
		claims.Username = claims.Subject
	}
	return claims, nil
}
