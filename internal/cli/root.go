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
	}

	// TODO: コマンドファイルを実装
	root.AddCommand(
		// newInstallCmd(),
		// newUseCmd(),
		// newUninstallCmd(),
		// newListCmd(),
		// newCurrentCmd(),
		// newSyncCmd(),
		// newDoctorCmd(),
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
