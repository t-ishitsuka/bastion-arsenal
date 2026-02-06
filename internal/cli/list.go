package cli

import (
	"fmt"

	"github.com/arsenal/internal/terminal"
	"github.com/spf13/cobra"
)

func newListCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "ls <tool>",
		Short: "インストール済みバージョン一覧を表示",
		Long: `指定したツールのインストール済みバージョン一覧を表示します。

現在アクティブなバージョンには * マークが付きます。

使用例:
  arsenal ls node
  arsenal ls go`,
		Args: cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			return runList(args[0])
		},
	}
}

func runList(toolName string) error {
	// プラグイン情報を取得
	p, err := registry.Get(toolName)
	if err != nil {
		return err
	}

	// インストール済みバージョンを取得
	versions, err := manager.List(toolName)
	if err != nil {
		return fmt.Errorf("バージョン一覧取得エラー: %w", err)
	}

	if len(versions) == 0 {
		terminal.PrintfYellow("%s のインストール済みバージョンがありません\n", p.DisplayName)
		fmt.Println()
		terminal.PrintlnCyan("インストールするには:")
		fmt.Printf("  arsenal install %s <version>\n", toolName)
		return nil
	}

	// 現在のバージョンを取得
	current, err := manager.Current(toolName)
	if err != nil {
		return fmt.Errorf("現在のバージョン取得エラー: %w", err)
	}

	// 表示
	terminal.PrintfBlue("%s のインストール済みバージョン:\n", p.DisplayName)
	fmt.Println()

	for _, version := range versions {
		if version == current {
			fmt.Printf("  * %s %s\n", terminal.Green(version), terminal.Yellow("(現在使用中)"))
		} else {
			fmt.Printf("    %s\n", version)
		}
	}

	return nil
}
