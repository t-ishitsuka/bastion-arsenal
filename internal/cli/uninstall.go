package cli

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/arsenal/internal/terminal"
)

func newUninstallCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "uninstall <tool> <version>",
		Short: "ツールの指定バージョンをアンインストール",
		Long: `指定したツールの特定バージョンをアンインストールします。

現在アクティブなバージョンをアンインストールする場合、
他にインストール済みバージョンがあれば自動的に最新版に切り替わります。
最後の1つをアンインストールする場合は symlink が削除されます。

使用例:
  arsenal uninstall node 18.0.0
  arsenal uninstall go 1.21.0`,
		Args: cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			return runUninstall(args[0], args[1])
		},
	}
}

func runUninstall(toolName, version string) error {
	// プラグイン情報を取得（存在確認）
	p, err := registry.Get(toolName)
	if err != nil {
		return err
	}

	// 現在のバージョンか確認
	current, _ := manager.Current(toolName)
	isCurrentVersion := current == version

	if isCurrentVersion {
		terminal.PrintWarning("%s %s は現在アクティブなバージョンです", p.DisplayName, version)
	}

	terminal.PrintfBlue("🗑️  %s %s をアンインストールします\n", p.DisplayName, version)

	// アンインストール前に他のバージョン一覧を取得
	versions, _ := manager.List(toolName)
	var remainingVersions []string
	for _, v := range versions {
		if v != version {
			remainingVersions = append(remainingVersions, v)
		}
	}

	// アンインストール実行
	if err := manager.Uninstall(toolName, version); err != nil {
		return err
	}

	// 現在のバージョンを削除した場合、他のバージョンがあれば自動切り替え
	if isCurrentVersion && len(remainingVersions) > 0 {
		// 最新のバージョン（リストの最後）に切り替え
		latestVersion := remainingVersions[len(remainingVersions)-1]
		if err := manager.Use(toolName, latestVersion); err != nil {
			terminal.PrintWarning("%s に自動切り替えできませんでした: %v", latestVersion, err)
		} else {
			terminal.PrintfCyan("🔄 自動的に %s %s に切り替えました\n", p.DisplayName, latestVersion)
		}
	}

	// 他にインストール済みバージョンがあるか確認
	if len(remainingVersions) > 0 {
		fmt.Println()
		terminal.PrintlnBlue("他のインストール済みバージョン:")
		for _, v := range remainingVersions {
			if isCurrentVersion && v == remainingVersions[len(remainingVersions)-1] {
				fmt.Printf("  * %s %s\n", terminal.Green(v), terminal.Yellow("(現在使用中)"))
			} else {
				fmt.Printf("    %s\n", v)
			}
		}
	} else {
		fmt.Println()
		terminal.PrintfYellow("%s のインストール済みバージョンがなくなりました\n", p.DisplayName)
	}

	return nil
}
