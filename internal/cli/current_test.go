package cli

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/arsenal/internal/config"
	"github.com/arsenal/internal/plugin"
	"github.com/arsenal/internal/version"
)

// newCurrentCmd が正しく作成されるかテストする
func TestNewCurrentCmd(t *testing.T) {
	cmd := newCurrentCmd()

	if cmd.Use != "current" {
		t.Errorf("Use = %q, want %q", cmd.Use, "current")
	}
}

// runCurrent が正しく動作するかテストする（アクティブなツールがない場合）
func TestRunCurrentEmpty(t *testing.T) {
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

	// runCurrent を実行
	err = runCurrent()
	if err != nil {
		t.Errorf("runCurrent() エラー: %v", err)
	}
}

// runCurrent が正しく動作するかテストする（アクティブなツールがある場合）
func TestRunCurrentWithActive(t *testing.T) {
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
	nodeVersionDir := filepath.Join(paths.Versions, "node", "20.10.0")
	if err := os.MkdirAll(nodeVersionDir, 0755); err != nil {
		t.Fatalf("バージョンディレクトリ作成エラー: %v", err)
	}

	// symlink を作成
	symlinkPath := filepath.Join(paths.Current, "node")
	if err := os.Symlink(nodeVersionDir, symlinkPath); err != nil {
		t.Fatalf("symlink 作成エラー: %v", err)
	}

	var err error
	registry, err = plugin.NewRegistry(paths)
	if err != nil {
		t.Fatalf("レジストリ作成エラー: %v", err)
	}

	manager = version.NewManager(paths, registry)

	// runCurrent を実行
	err = runCurrent()
	if err != nil {
		t.Errorf("runCurrent() エラー: %v", err)
	}
}
