package teldrive

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"reflect"

	"github.com/spf13/cobra"
	"go.uber.org/zap"

	"github.com/tgdrive/teldrive/internal/config"
	"github.com/tgdrive/teldrive/internal/database"
	"github.com/tgdrive/teldrive/internal/telegram/client"
	"github.com/tgdrive/teldrive/internal/utils/logger"
	"github.com/tgdrive/teldrive/pkg/services"
)

func newCheckCommand() *cobra.Command {
	var cfg config.CheckConfig
	loader := config.NewLoader()

	cmd := &cobra.Command{
		Use:   "check",
		Short: "Check and clean Telegram storage",
		Long: `Check for orphaned files in Telegram channels and optionally clean them up.
This command helps maintain storage efficiency by removing files that are no longer referenced in the database.`,
		RunE: func(cmd *cobra.Command, args []string) error {
			return runCheck(cmd.Context(), &cfg)
		},
		PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
			if err := loader.Load(cmd, &cfg); err != nil {
				return fmt.Errorf("failed to load config: %w", err)
			}
			return nil
		},
	}

	loader.RegisterFlags(cmd.Flags(), reflect.TypeFor[config.CheckConfig]())
	return cmd
}

func runCheck(ctx context.Context, cfg *config.CheckConfig) error {
	log := logger.New(logger.Config{Level: cfg.Log.Level, Format: "console"})
	defer log.Sync()

	log.Info("starting check operation",
		"dry_run", cfg.DryRun,
		"user", cfg.User,
	)

	db, err := database.Initialize(cfg.DB, log)
	if err != nil {
		return fmt.Errorf("failed to initialize database: %w", err)
	}
	defer database.Close(db)

	tgClient, err := client.New(cfg.TG, log)
	if err != nil {
		return fmt.Errorf("failed to initialize telegram client: %w", err)
	}
	defer tgClient.Close()

	checker := services.NewChecker(db, tgClient, log)
	results, err := checker.Run(ctx, cfg)
	if err != nil {
		return fmt.Errorf("check operation failed: %w", err)
	}

	if cfg.ExportFile != "" {
		data, err := json.MarshalIndent(results, "", "  ")
		if err != nil {
			return fmt.Errorf("failed to marshal results: %w", err)
		}
		if err := os.WriteFile(cfg.ExportFile, data, 0644); err != nil {
			return fmt.Errorf("failed to write results: %w", err)
		}
		log.Info("results exported", "file", cfg.ExportFile)
	}

	log.Info("check complete",
		"orphaned_files", len(results.OrphanedFiles),
		"cleaned_files", results.CleanedCount,
	)

	return nil
}
