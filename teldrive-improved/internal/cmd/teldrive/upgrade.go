package teldrive

import (
	"fmt"
	"runtime"

	"github.com/Masterminds/semver/v3"
	"github.com/spf13/cobra"

	"github.com/tgdrive/teldrive/internal/version"
	"github.com/tgdrive/teldrive/internal/utils/updater"
)

func newUpgradeCommand() *cobra.Command {
	return &cobra.Command{
		Use:   "upgrade",
		Short: "Upgrade TelDrive to the latest version",
		RunE: func(cmd *cobra.Command, args []string) error {
			v := version.Get()
			current, err := semver.NewVersion(v.Version)
			if err != nil {
				return fmt.Errorf("invalid current version: %w", err)
			}

			up := updater.New(runtime.GOOS, runtime.GOARCH)
			latest, err := up.GetLatestVersion(cmd.Context())
			if err != nil {
				return fmt.Errorf("failed to check for updates: %w", err)
			}

			if !latest.GreaterThan(current) {
				fmt.Printf("Already up to date (v%s)\n", current)
				return nil
			}

			fmt.Printf("Upgrading from v%s to v%s...\n", current, latest)
			if err := up.DownloadAndInstall(cmd.Context(), latest.String()); err != nil {
				return fmt.Errorf("upgrade failed: %w", err)
			}

			fmt.Println("Upgrade complete! Please restart TelDrive.")
			return nil
		},
	}
}
