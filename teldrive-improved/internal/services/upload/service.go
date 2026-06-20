package upload

import (
	"context"
	"fmt"
	"io"

	"github.com/tgdrive/teldrive/internal/app"
	"go.uber.org/zap"
)

type Service struct {
	ctx *app.Context
	log *zap.Logger
}

type StartRequest struct {
	Name     string `json:"name" validate:"required"`
	Size     int64  `json:"size" validate:"required"`
	ParentID string `json:"parent_id,omitempty"`
	MimeType string `json:"mime_type"`
}

type StartResponse struct {
	UploadID    string `json:"upload_id"`
	ChunkSize   int64  `json:"chunk_size"`
	TotalChunks int    `json:"total_chunks"`
}

type FileInfo struct {
	ID   string `json:"id"`
	Name string `json:"name"`
	Size int64  `json:"size"`
}

func NewService(ctx *app.Context) *Service {
	return &Service{
		ctx: ctx,
		log: ctx.Logger.Named("upload"),
	}
}

func (s *Service) Start(ctx context.Context, req *StartRequest) (*StartResponse, error) {
	return &StartResponse{
		UploadID:    "upload-123",
		ChunkSize:   1024 * 1024,
		TotalChunks: int(req.Size / (1024 * 1024)),
	}, nil
}

func (s *Service) UploadChunk(ctx context.Context, uploadID string, data io.Reader) error {
	return nil
}

func (s *Service) Complete(ctx context.Context, uploadID string) (*FileInfo, error) {
	return nil, fmt.Errorf("not implemented")
}
