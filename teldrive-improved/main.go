package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/tgdrive/teldrive/cmd/teldrive"
	"github.com/tgdrive/teldrive/internal/utils/logger"
)

var (
	Version   = "dev"
	Commit    = "unknown"
	BuildTime = "unknown"
)

func main() {
	ctx, cancel := signal.NotifyContext(
		context.Background(),
		os.Interrupt,
		syscall.SIGTERM,
		syscall.SIGQUIT,
	)
	defer cancel()

	log := logger.New(logger.Config{Level: "info", Format: "json"})
	log.Info("starting teldrive",
		"version", Version,
		"commit", Commit,
		"build_time", BuildTime,
	)

	rootCmd := teldrive.NewRootCommand(Version, Commit, BuildTime)
	if err := rootCmd.ExecuteContext(ctx); err != nil {
		log.Error("application error", "error", err)
		os.Exit(1)
	}

	shutdownCtx, shutdownCancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer shutdownCancel()

	select {
	case <-shutdownCtx.Done():
		fmt.Fprintln(os.Stderr, "shutdown timed out")
		os.Exit(1)
	default:
		log.Info("shutdown complete")
	}
}
