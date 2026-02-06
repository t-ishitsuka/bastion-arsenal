package cli

import (
	"fmt"

	"github.com/arsenal/internal/terminal"
	"github.com/arsenal/internal/version"
	"github.com/spf13/cobra"
)

func newLsRemoteCmd() *cobra.Command {
	var limit int
	var all bool
	var ltsOnly bool

	cmd := &cobra.Command{
		Use:   "ls-remote <tool>",
		Short: "リモートの利用可能なバージョン一覧を表示",
		Long: `指定したツールの、リモートから取得可能なバージョン一覧を表示します。

デフォルトでは最新20件を表示します。

使用例:
  arsenal ls-remote node
  arsenal ls-remote node --limit 50
  arsenal ls-remote node --all
  arsenal ls-remote node --lts-only`,
		Args: cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			// --all が指定された場合は limit を 0 に設定（無制限）
			if all {
				limit = 0
			}
			return runLsRemote(args[0], limit, ltsOnly)
		},
	}

	cmd.Flags().IntVarP(&limit, "limit", "n", 20, "表示件数（0で全件表示）")
	cmd.Flags().BoolVar(&all, "all", false, "全バージョンを表示")
	cmd.Flags().BoolVar(&ltsOnly, "lts-only", false, "LTS バージョンのみ表示")

	return cmd
}

func runLsRemote(toolName string, limit int, ltsOnly bool) error {
	// プラグイン情報を取得
	p, err := registry.Get(toolName)
	if err != nil {
		return err
	}

	terminal.PrintfBlue("📡 %s の利用可能なバージョン一覧を取得中...\n", p.DisplayName)
	fmt.Println()

	// リモートからバージョン一覧を取得
	versions, err := manager.ListRemote(toolName, 0) // 0 = 無制限で取得してからフィルタリング
	if err != nil {
		return err
	}

	// LTS のみフィルタリング
	if ltsOnly {
		var ltsVersions []version.RemoteVersion
		for _, v := range versions {
			if v.LTS != "" {
				ltsVersions = append(ltsVersions, v)
			}
		}
		versions = ltsVersions
	}

	// 件数制限を適用
	if limit > 0 && len(versions) > limit {
		versions = versions[:limit]
	}

	if len(versions) == 0 {
		if ltsOnly {
			terminal.PrintfYellow("%s の LTS バージョンが見つかりませんでした\n", p.DisplayName)
		} else {
			terminal.PrintfYellow("%s の利用可能なバージョンが見つかりませんでした\n", p.DisplayName)
		}
		return nil
	}

	// 表示
	header := fmt.Sprintf("%s の利用可能なバージョン", p.DisplayName)
	if ltsOnly {
		header += "（LTS のみ）"
	}
	if limit > 0 && len(versions) == limit {
		header += fmt.Sprintf("（最新 %d 件）", limit)
	}
	terminal.PrintlnBlue(header + ":")
	fmt.Println()

	for _, v := range versions {
		if v.LTS != "" {
			fmt.Printf("  %s %s\n", terminal.Green(v.Version), terminal.Yellow("(LTS: "+v.LTS+")"))
		} else {
			fmt.Printf("  %s\n", v.Version)
		}
	}

	if limit > 0 && len(versions) == limit && !ltsOnly {
		fmt.Println()
		terminal.PrintlnCyan("全バージョンを表示するには --all を使用してください:")
		fmt.Printf("  arsenal ls-remote %s --all\n", toolName)
	}

	return nil
}
