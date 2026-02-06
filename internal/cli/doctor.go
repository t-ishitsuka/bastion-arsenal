package cli

import (
	"fmt"

	"github.com/arsenal/internal/version"
	"github.com/spf13/cobra"
)

func newDoctorCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "doctor",
		Short: "環境ヘルスチェック",
		Long: `Arsenal の環境設定をチェックします。

以下の項目を確認します:
  - 必要なディレクトリの存在確認
  - PATH 環境変数の設定確認
  - インストール済みツールの確認`,
		RunE: func(cmd *cobra.Command, args []string) error {
			return runDoctor()
		},
	}
}

func runDoctor() error {
	fmt.Println("Arsenal 環境をチェック中...")
	fmt.Println()

	// 診断を実行
	results := manager.Doctor()

	// 結果を表示
	hasWarnings := false
	hasErrors := false

	for _, result := range results {
		var icon string
		switch result.Status {
		case version.StatusOK:
			icon = "✓"
		case version.StatusWarn:
			icon = "⚠"
			hasWarnings = true
		case version.StatusError:
			icon = "✗"
			hasErrors = true
		}

		fmt.Printf("%s %s: %s\n", icon, result.Name, result.Message)
	}

	fmt.Println()

	// サマリーを表示
	if hasErrors {
		fmt.Println("エラーが検出されました。上記のエラーを修正してください。")
		return fmt.Errorf("環境チェックでエラーが検出されました")
	} else if hasWarnings {
		fmt.Println("警告があります。必要に応じて対応してください。")
	} else {
		fmt.Println("✓ 全てのチェックに合格しました")
	}

	return nil
}
