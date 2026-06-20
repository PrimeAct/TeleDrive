package user

import (
	"context"

	"github.com/tgdrive/teldrive/internal/app"
	"go.uber.org/zap"
)

type Service struct {
	ctx *app.Context
	log *zap.Logger
}

type Stats struct {
	TotalFiles   int64 `json:"total_files"`
	TotalSize    int64 `json:"total_size"`
	StorageUsed  int64 `json:"storage_used"`
	StorageLimit int64 `json:"storage_limit"`
}

func NewService(ctx *app.Context) *Service {
	return &Service{
		ctx: ctx,
		log: ctx.Logger.Named("user"),
	}
}

func (s *Service) GetStats(ctx context.Context) (*Stats, error) {
	return &Stats{
		TotalFiles:   0,
		TotalSize:    0,
		StorageUsed:  0,
		StorageLimit: 1024 * 1024 * 1024 * 1024,
	}, nil
}
