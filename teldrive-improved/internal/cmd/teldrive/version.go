package teldrive

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/tgdrive/teldrive/internal/version"
)

func newVersionCommand() *cobra.Command {
	return &cobra.Command{
		Use:   "version",
		Short: "Print version information",
		Run: func(cmd *cobra.Command, args []string) {
			v := version.Get()
			fmt.Printf("TelDrive %s (commit: %s, built: %s)\n", v.Version, v.Commit, v.BuildTime)
		},
	}
}
