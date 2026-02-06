package cli

import (
	"fmt"

	"github.com/spf13/cobra"
)

func newVersionCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "version",
		Short: "バージョン情報を表示",
		Long:  "Arsenal のバージョン情報を表示します。",
		RunE: func(cmd *cobra.Command, args []string) error {
			fmt.Fprintf(cmd.OutOrStdout(), "Arsenal %s\n", versionInfo.Version)
			fmt.Fprintf(cmd.OutOrStdout(), "Commit: %s\n", versionInfo.Commit)
			fmt.Fprintf(cmd.OutOrStdout(), "Built: %s\n", versionInfo.BuildDate)
			return nil
		},
	}
}
