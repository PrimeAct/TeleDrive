package client

import (
	"context"
	"fmt"
	"sync"

	"go.uber.org/zap"

	"github.com/tgdrive/teldrive/internal/config"
)

type Pool struct {
	clients []Client
	mu      sync.RWMutex
	log     *zap.Logger
	config  config.TGConfig
}

type Client interface {
	Connect(ctx context.Context) error
	Close() error
	Upload(ctx context.Context, data []byte, filename string) (string, error)
	Download(ctx context.Context, fileID string) ([]byte, error)
	Delete(ctx context.Context, fileID string) error
}

func NewPool(cfg config.TGConfig, log *zap.Logger) (*Pool, error) {
	pool := &Pool{
		clients: make([]Client, 0, cfg.Workers),
		log:     log,
		config:  cfg,
	}

	for i := 0; i < cfg.Workers; i++ {
		client := NewClient(cfg, log.Named(fmt.Sprintf("tg-client-%d", i)))
		if err := client.Connect(context.Background()); err != nil {
			log.Warn("failed to connect telegram client", "worker", i, "error", err)
			continue
		}
		pool.clients = append(pool.clients, client)
	}

	if len(pool.clients) == 0 {
		return nil, fmt.Errorf("no telegram clients could be initialized")
	}

	log.Info("telegram client pool initialized", "workers", len(pool.clients))
	return pool, nil
}

func (p *Pool) Get() Client {
	p.mu.RLock()
	defer p.mu.RUnlock()
	if len(p.clients) == 0 {
		return nil
	}
	return p.clients[0]
}

func (p *Pool) Close() error {
	p.mu.Lock()
	defer p.mu.Unlock()
	for _, client := range p.clients {
		if err := client.Close(); err != nil {
			p.log.Warn("error closing client", "error", err)
		}
	}
	return nil
}
