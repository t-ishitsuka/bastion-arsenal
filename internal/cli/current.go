package cli

import (
	"fmt"
	"sort"

	"github.com/spf13/cobra"
	"github.com/arsenal/internal/terminal"
)

func newCurrentCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "current",
		Short: "現在アクティブなツールバージョンを表示",
		Long: `現在アクティブな全ツールのバージョンを表示します。

symlink で設定されているバージョンを確認できます。
~/.arsenal/current/ ディレクトリの内容を表示します。`,
		RunE: func(cmd *cobra.Command, args []string) error {
			return runCurrent()
		},
	}
}

func runCurrent() error {
	// 全ツールの現在のバージョンを取得
	currentAll, err := manager.CurrentAll()
	if err != nil {
		return fmt.Errorf("アクティブバージョン取得エラー: %w", err)
	}

	if len(currentAll) == 0 {
		terminal.PrintlnYellow("アクティブなツールがありません")
		fmt.Println()
		terminal.PrintlnCyan("ツールをインストールして使用するには:")
		fmt.Println("  arsenal install <tool> <version>")
		fmt.Println("  arsenal use <tool> <version>")
		return nil
	}

	terminal.PrintlnBlue("アクティブなツール:")
	fmt.Println()

	// ツール名でソート
	tools := make([]string, 0, len(currentAll))
	for tool := range currentAll {
		tools = append(tools, tool)
	}
	sort.Strings(tools)

	// 表示
	for _, tool := range tools {
		version := currentAll[tool]
		// プラグイン情報を取得して表示名を使う
		p, err := registry.Get(tool)
		if err == nil {
			fmt.Printf("  %s: %s\n", p.DisplayName, terminal.Green(version))
		} else {
			fmt.Printf("  %s: %s\n", tool, terminal.Green(version))
		}
	}

	return nil
}
