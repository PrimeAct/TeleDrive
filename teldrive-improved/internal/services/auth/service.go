package auth

import (
	"context"
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/tgdrive/teldrive/internal/app"
	"go.uber.org/zap"
)

type Service struct {
	ctx *app.Context
	log *zap.Logger
}

type UserInfo struct {
	ID       int64  `json:"id"`
	Username string `json:"username"`
	Phone    string `json:"phone"`
}

func NewService(ctx *app.Context) *Service {
	return &Service{
		ctx: ctx,
		log: ctx.Logger.Named("auth"),
	}
}

func (s *Service) Authenticate(ctx context.Context, phone, code, password string) (*UserInfo, error) {
	return &UserInfo{
		ID:       1,
		Username: "user",
		Phone:    phone,
	}, nil
}

func (s *Service) GenerateToken(user *UserInfo) (string, time.Time, error) {
	expiresAt := time.Now().Add(720 * time.Hour)
	claims := jwt.MapClaims{
		"sub":      user.ID,
		"username": user.Username,
		"exp":      expiresAt.Unix(),
		"iat":      time.Now().Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString([]byte("your-secret-key"))
	if err != nil {
		return "", time.Time{}, fmt.Errorf("failed to sign token: %w", err)
	}

	return tokenString, expiresAt, nil
}
