package cli

import (
	"fmt"

	"github.com/spf13/cobra"
)

func newInitShellCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "init-shell [bash|zsh|fish]",
		Short: "シェル設定スクリプトを生成",
		Long: `指定したシェル用の初期化スクリプトを生成します。

このコマンドの出力をシェルの設定ファイルに追加してください。

使用例:
  # Bash の場合 (~/.bashrc に追加)
  bastion-arsenal init-shell bash >> ~/.bashrc

  # Zsh の場合 (~/.zshrc に追加)
  bastion-arsenal init-shell zsh >> ~/.zshrc

  # Fish の場合 (~/.config/fish/config.fish に追加)
  bastion-arsenal init-shell fish >> ~/.config/fish/config.fish`,
		Args:      cobra.ExactArgs(1),
		ValidArgs: []string{"bash", "zsh", "fish"},
		RunE: func(cmd *cobra.Command, args []string) error {
			return runInitShell(args[0])
		},
	}
}

func runInitShell(shell string) error {
	arsenalDir := paths.Root

	switch shell {
	case "bash", "zsh":
		fmt.Println("# Arsenal の初期化")
		fmt.Printf("export PATH=\"%s/current/*/bin:$PATH\"\n", arsenalDir)
		fmt.Println()
		fmt.Println("# Arsenal の補完を有効化")
		fmt.Printf("eval \"$(bastion-arsenal completion %s)\"\n", shell)

	case "fish":
		fmt.Println("# Arsenal の初期化")
		fmt.Printf("set -gx PATH %s/current/*/bin $PATH\n", arsenalDir)
		fmt.Println()
		fmt.Println("# Arsenal の補完を有効化")
		fmt.Println("bastion-arsenal completion fish | source")

	default:
		return fmt.Errorf("サポートされていないシェル: %s (bash, zsh, fish のみ対応)", shell)
	}

	return nil
}
