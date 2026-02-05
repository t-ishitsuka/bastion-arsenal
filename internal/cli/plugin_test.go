package cli

import (
	"testing"

	"github.com/arsenal/internal/config"
	"github.com/arsenal/internal/plugin"
)

// newPluginCmd が正しく作成されるかテストする
func TestNewPluginCmd(t *testing.T) {
	cmd := newPluginCmd()

	if cmd.Use != "plugin" {
		t.Errorf("Use = %q, want %q", cmd.Use, "plugin")
	}

	// サブコマンドが存在するか確認
	if len(cmd.Commands()) == 0 {
		t.Error("サブコマンドが登録されていません")
	}

	// list サブコマンドが存在するか確認
	listCmd, _, err := cmd.Find([]string{"list"})
	if err != nil {
		t.Errorf("list サブコマンドが見つかりません: %v", err)
		return
	}
	if listCmd == nil {
		t.Error("list サブコマンドが nil です")
		return
	}
}

// runPluginList が正しく動作するかテストする
func TestRunPluginList(t *testing.T) {
	// テスト用のレジストリをセットアップ
	tmpDir := t.TempDir()
	paths := &config.Paths{
		Plugins: tmpDir,
	}

	var err error
	registry, err = plugin.NewRegistry(paths)
	if err != nil {
		t.Fatalf("レジストリ作成エラー: %v", err)
	}

	// runPluginList を実行
	err = runPluginList()
	if err != nil {
		t.Errorf("runPluginList() エラー: %v", err)
	}
}
