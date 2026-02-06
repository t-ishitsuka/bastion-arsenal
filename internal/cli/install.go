package cli

import (
	"fmt"

	"github.com/spf13/cobra"
)

func newInstallCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "install <tool> <version>",
		Short: "ツールの指定バージョンをインストール",
		Long: `指定したツールの特定バージョンをダウンロードしてインストールします。

使用例:
  arsenal install node 20.10.0
  arsenal install go 1.22.0`,
		Args: cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			return runInstall(args[0], args[1])
		},
	}
}

func runInstall(toolName, version string) error {
	// プラグイン情報を取得（存在確認）
	p, err := registry.Get(toolName)
	if err != nil {
		return err
	}

	fmt.Printf("%s %s をインストールします\n", p.DisplayName, version)
	fmt.Println()

	// インストール実行
	if err := manager.Install(toolName, version); err != nil {
		return err
	}

	fmt.Println()
	fmt.Println("次のコマンドで有効化できます:")
	fmt.Printf("  arsenal use %s %s\n", toolName, version)

	return nil
}
