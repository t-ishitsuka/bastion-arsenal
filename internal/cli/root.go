package cli

import (
	"fmt"
	"os"

	"github.com/arsenal/internal/config"
	"github.com/arsenal/internal/plugin"
	"github.com/arsenal/internal/version"
	"github.com/spf13/cobra"
)

var (
	paths    *config.Paths
	registry *plugin.Registry
	manager  *version.Manager
)

func NewRootCmd() *cobra.Command {
	root := &cobra.Command{
		Use:   "arsenal",
		Short: "⚔️  Arsenal - 軽量マルチランタイムバージョンマネージャー",
		Long: `Arsenal は複数のランタイム/SDK バージョンを単一の CLI から管理します。

  arsenal install node 20.10.0    特定バージョンをインストール
  arsenal use node 20.10.0        バージョンを切り替え
  arsenal ls node                 インストール済みバージョン一覧
  arsenal sync                    .toolversions から同期
  arsenal current                 アクティブバージョンを表示
  arsenal doctor                  環境ヘルスチェック`,
		SilenceUsage:  true,
		SilenceErrors: true,
		PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
			return initialize()
		},
		CompletionOptions: cobra.CompletionOptions{
			DisableDefaultCmd: false,
		},
	}

	// completion コマンドの説明を日本語化
	// completion コマンドをカスタマイズ
	completionCmd := &cobra.Command{
		Use:   "completion [bash|zsh|fish|powershell]",
		Short: "シェルの補完スクリプトを生成",
		Long: `指定したシェル用の補完スクリプトを生成します。

使用例:
  # Bash の場合 (~/.bashrc に追加)
  source <(arsenal completion bash)

  # Zsh の場合 (~/.zshrc に追加)
  source <(arsenal completion zsh)

  # Fish の場合
  arsenal completion fish | source`,
		ValidArgs: []string{"bash", "zsh", "fish", "powershell"},
		Args:      cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			switch args[0] {
			case "bash":
				return root.GenBashCompletion(cmd.OutOrStdout())
			case "zsh":
				return root.GenZshCompletion(cmd.OutOrStdout())
			case "fish":
				return root.GenFishCompletion(cmd.OutOrStdout(), true)
			case "powershell":
				return root.GenPowerShellCompletion(cmd.OutOrStdout())
			default:
				return fmt.Errorf("サポートされていないシェル: %s", args[0])
			}
		},
	}
	root.AddCommand(completionCmd)

	// help コマンドの説明を日本語化
	root.SetHelpCommand(&cobra.Command{
		Use:   "help [command]",
		Short: "任意のコマンドのヘルプを表示",
		Long: `任意のコマンドの詳細なヘルプ情報を表示します。

使用例:
  arsenal help
  arsenal help plugin
  arsenal help plugin list`,
		Run: func(c *cobra.Command, args []string) {
			cmd, _, e := root.Find(args)
			if e != nil || cmd == nil {
				if len(args) > 0 {
					c.Printf("不明なコマンド: %q\n", args)
				}
				_ = root.Usage()
				return
			}
			cmd.InitDefaultHelpFlag()
			_ = cmd.Help()
		},
	})

	// TODO: コマンドファイルを実装
	root.AddCommand(
		newInstallCmd(),
		newUseCmd(),
		newUninstallCmd(),
		newListCmd(),
		newLsRemoteCmd(),
		newCurrentCmd(),
		// newSyncCmd(),
		newDoctorCmd(),
		newPluginCmd(),
		// newInitShellCmd(),
	)

	return root
}

// Arsenal の初期化処理を行う
func initialize() error {
	var err error

	paths, err = config.GetPaths()
	if err != nil {
		return fmt.Errorf("パス取得エラー: %w", err)
	}

	if err := paths.EnsureDirs(); err != nil {
		return fmt.Errorf("ディレクトリ作成エラー: %w", err)
	}

	registry, err = plugin.NewRegistry(paths)
	if err != nil {
		return fmt.Errorf("プラグイン読み込みエラー: %w", err)
	}

	manager = version.NewManager(paths, registry)
	_ = manager // 将来の CLI コマンドで使用予定
	return nil
}

// ルートコマンドを実行する
func Execute() {
	if err := NewRootCmd().Execute(); err != nil {
		fmt.Fprintf(os.Stderr, "エラー: %v\n", err)
		os.Exit(1)
	}
}
