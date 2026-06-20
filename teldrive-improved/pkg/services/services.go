package services

import (
	"github.com/tgdrive/teldrive/internal/app"
	"github.com/tgdrive/teldrive/internal/services/auth"
	"github.com/tgdrive/teldrive/internal/services/file"
	"github.com/tgdrive/teldrive/internal/services/share"
	"github.com/tgdrive/teldrive/internal/services/upload"
	"github.com/tgdrive/teldrive/internal/services/user"
	"github.com/tgdrive/teldrive/internal/telegram/client"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

type Services struct {
	Auth   *auth.Service
	File   *file.Service
	Upload *upload.Service
	Share  *share.Service
	User   *user.Service
}

func Initialize(ctx *app.Context) *Services {
	return &Services{
		Auth:   auth.NewService(ctx),
		File:   file.NewService(ctx),
		Upload: upload.NewService(ctx),
		Share:  share.NewService(ctx),
		User:   user.NewService(ctx),
	}
}

type Checker struct {
	db       *gorm.DB
	tgClient client.Client
	log      *zap.Logger
}

type CheckResults struct {
	OrphanedFiles []string `json:"orphaned_files"`
	CleanedCount  int      `json:"cleaned_count"`
}

func NewChecker(db *gorm.DB, tgClient client.Client, log *zap.Logger) *Checker {
	return &Checker{
		db:       db,
		tgClient: tgClient,
		log:      log,
	}
}

func (c *Checker) Run(ctx context.Context, cfg interface{}) (*CheckResults, error) {
	return &CheckResults{
		OrphanedFiles: []string{},
		CleanedCount:  0,
	}, nil
}
