package file

import (
	"context"
	"fmt"

	"github.com/tgdrive/teldrive/internal/app"
	"go.uber.org/zap"
)

type Service struct {
	ctx *app.Context
	log *zap.Logger
}

type FileInfo struct {
	ID        string `json:"id"`
	Name      string `json:"name"`
	Size      int64  `json:"size"`
	MimeType  string `json:"mime_type"`
	ParentID  string `json:"parent_id,omitempty"`
	ChannelID int64  `json:"channel_id"`
	MessageID int    `json:"message_id"`
	CreatedAt string `json:"created_at"`
	UpdatedAt string `json:"updated_at"`
}

func NewService(ctx *app.Context) *Service {
	return &Service{
		ctx: ctx,
		log: ctx.Logger.Named("file"),
	}
}

func (s *Service) List(ctx context.Context, page, pageSize int) ([]*FileInfo, int64, error) {
	return []*FileInfo{}, 0, nil
}

func (s *Service) Get(ctx context.Context, id string) (*FileInfo, error) {
	return nil, fmt.Errorf("file not found")
}

func (s *Service) Create(ctx context.Context, info *FileInfo) error {
	return nil
}

func (s *Service) Update(ctx context.Context, id string, info *FileInfo) error {
	return nil
}

func (s *Service) Delete(ctx context.Context, id string) error {
	return nil
}
