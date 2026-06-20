package share

import (
	"context"

	"github.com/tgdrive/teldrive/internal/app"
	"go.uber.org/zap"
)

type Service struct {
	ctx *app.Context
	log *zap.Logger
}

type ShareInfo struct {
	ID        string `json:"id"`
	FileID    string `json:"file_id"`
	Token     string `json:"token"`
	ExpiresAt string `json:"expires_at,omitempty"`
}

func NewService(ctx *app.Context) *Service {
	return &Service{
		ctx: ctx,
		log: ctx.Logger.Named("share"),
	}
}

func (s *Service) Create(ctx context.Context, fileID string, expiresIn int) (*ShareInfo, error) {
	return &ShareInfo{}, nil
}

func (s *Service) GetByToken(ctx context.Context, token string) (*ShareInfo, error) {
	return nil, nil
}
