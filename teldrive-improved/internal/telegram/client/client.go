package client

import (
	"context"
	"fmt"

	"go.uber.org/zap"

	"github.com/tgdrive/teldrive/internal/config"
)

type TelegramClient struct {
	cfg config.TGConfig
	log *zap.Logger
}

func NewClient(cfg config.TGConfig, log *zap.Logger) *TelegramClient {
	return &TelegramClient{
		cfg: cfg,
		log: log,
	}
}

func (c *TelegramClient) Connect(ctx context.Context) error {
	c.log.Info("connecting to telegram")
	return nil
}

func (c *TelegramClient) Close() error {
	c.log.Info("closing telegram client")
	return nil
}

func (c *TelegramClient) Upload(ctx context.Context, data []byte, filename string) (string, error) {
	return "", fmt.Errorf("not implemented")
}

func (c *TelegramClient) Download(ctx context.Context, fileID string) ([]byte, error) {
	return nil, fmt.Errorf("not implemented")
}

func (c *TelegramClient) Delete(ctx context.Context, fileID string) error {
	return fmt.Errorf("not implemented")
}
