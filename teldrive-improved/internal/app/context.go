package app

import (
	"context"

	"go.uber.org/zap"
	"gorm.io/gorm"

	"github.com/tgdrive/teldrive/internal/telegram/client"
	"github.com/tgdrive/teldrive/internal/utils/cache"
)

type Context struct {
	ctx    context.Context
	DB     *gorm.DB
	Cache  cache.Cache
	TGPool *client.Pool
	Logger *zap.Logger
}

func NewContext(ctx context.Context, db *gorm.DB, cache cache.Cache, tgPool *client.Pool, logger *zap.Logger) *Context {
	return &Context{
		ctx:    ctx,
		DB:     db,
		Cache:  cache,
		TGPool: tgPool,
		Logger: logger,
	}
}

func (c *Context) Context() context.Context {
	return c.ctx
}

func (c *Context) Done() <-chan struct{} {
	return c.ctx.Done()
}
