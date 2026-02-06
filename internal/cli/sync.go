package cli

import (
	"os"

	"github.com/spf13/cobra"
)

func newSyncCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "sync",
		Short: ".toolversions からバージョンを同期",
		Long: `.toolversions ファイルに記載された全ツールのバージョンを
インストールして切り替えます。

.toolversions ファイルは現在のディレクトリから上位ディレクトリへと
遡って検索されます。

使用例:
  arsenal sync`,
		RunE: func(cmd *cobra.Command, args []string) error {
			return runSync()
		},
	}
}

func runSync() error {
	// カレントディレクトリを取得
	cwd, err := os.Getwd()
	if err != nil {
		return err
	}

	// sync 実行
	return manager.Sync(cwd)
}
