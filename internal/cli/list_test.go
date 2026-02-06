package cli

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/arsenal/internal/config"
	"github.com/arsenal/internal/plugin"
	"github.com/arsenal/internal/version"
)

// newListCmd が正しく作成されるかテストする
func TestNewListCmd(t *testing.T) {
	cmd := newListCmd()

	if cmd.Use != "ls <tool>" {
		t.Errorf("Use = %q, want %q", cmd.Use, "ls <tool>")
	}
}

// runList が正しく動作するかテストする（インストール済みバージョンがない場合）
func TestRunListEmpty(t *testing.T) {
	// テスト用の環境をセットアップ
	tmpDir := t.TempDir()
	paths := &config.Paths{
		Root:     filepath.Join(tmpDir, "arsenal"),
		Versions: filepath.Join(tmpDir, "arsenal", "versions"),
		Current:  filepath.Join(tmpDir, "arsenal", "current"),
		Plugins:  filepath.Join(tmpDir, "arsenal", "plugins"),
	}

	if err := paths.EnsureDirs(); err != nil {
		t.Fatalf("ディレクトリ作成エラー: %v", err)
	}

	var err error
	registry, err = plugin.NewRegistry(paths)
	if err != nil {
		t.Fatalf("レジストリ作成エラー: %v", err)
	}

	manager = version.NewManager(paths, registry)

	// runList を実行
	err = runList("node")
	if err != nil {
		t.Errorf("runList() エラー: %v", err)
	}
}

// runList が正しく動作するかテストする（インストール済みバージョンがある場合）
func TestRunListWithVersions(t *testing.T) {
	// テスト用の環境をセットアップ
	tmpDir := t.TempDir()
	paths := &config.Paths{
		Root:     filepath.Join(tmpDir, "arsenal"),
		Versions: filepath.Join(tmpDir, "arsenal", "versions"),
		Current:  filepath.Join(tmpDir, "arsenal", "current"),
		Plugins:  filepath.Join(tmpDir, "arsenal", "plugins"),
	}

	if err := paths.EnsureDirs(); err != nil {
		t.Fatalf("ディレクトリ作成エラー: %v", err)
	}

	// ダミーのバージョンディレクトリを作成
	versions := []string{"18.19.0", "20.10.0", "21.0.0"}
	for _, v := range versions {
		versionDir := filepath.Join(paths.Versions, "node", v)
		if err := os.MkdirAll(versionDir, 0755); err != nil {
			t.Fatalf("バージョンディレクトリ作成エラー: %v", err)
		}
	}

	// 現在のバージョンを設定（symlink）
	currentVersion := filepath.Join(paths.Versions, "node", "20.10.0")
	symlinkPath := filepath.Join(paths.Current, "node")
	if err := os.Symlink(currentVersion, symlinkPath); err != nil {
		t.Fatalf("symlink 作成エラー: %v", err)
	}

	var err error
	registry, err = plugin.NewRegistry(paths)
	if err != nil {
		t.Fatalf("レジストリ作成エラー: %v", err)
	}

	manager = version.NewManager(paths, registry)

	// runList を実行
	err = runList("node")
	if err != nil {
		t.Errorf("runList() エラー: %v", err)
	}
}

// runList が不明なツールでエラーを返すかテストする
func TestRunListUnknownTool(t *testing.T) {
	// テスト用の環境をセットアップ
	tmpDir := t.TempDir()
	paths := &config.Paths{
		Root:     filepath.Join(tmpDir, "arsenal"),
		Versions: filepath.Join(tmpDir, "arsenal", "versions"),
		Current:  filepath.Join(tmpDir, "arsenal", "current"),
		Plugins:  filepath.Join(tmpDir, "arsenal", "plugins"),
	}

	if err := paths.EnsureDirs(); err != nil {
		t.Fatalf("ディレクトリ作成エラー: %v", err)
	}

	var err error
	registry, err = plugin.NewRegistry(paths)
	if err != nil {
		t.Fatalf("レジストリ作成エラー: %v", err)
	}

	manager = version.NewManager(paths, registry)

	// runList を実行（存在しないツール）
	err = runList("nonexistent")
	if err == nil {
		t.Error("存在しないツールでエラーが返されませんでした")
	}
}
