package teldrive

import (
	"github.com/spf13/cobra"
	"github.com/tgdrive/teldrive/internal/version"
)

func NewRootCommand(version, commit, buildTime string) *cobra.Command {
	version.Set(version, commit, buildTime)

	cmd := &cobra.Command{
		Use:   "teldrive",
		Short: "TelDrive - Telegram File Storage Server",
		Long: `TelDrive is a powerful utility that enables you to organize
your Telegram files with an intuitive web interface and Rclone compatibility.`,
		SilenceUsage:  true,
		SilenceErrors: true,
		CompletionOptions: cobra.CompletionOptions{
			DisableDefaultCmd: true,
		},
	}

	cmd.AddCommand(
		newRunCommand(),
		newCheckCommand(),
		newUpgradeCommand(),
		newVersionCommand(),
	)

	return cmd
}
