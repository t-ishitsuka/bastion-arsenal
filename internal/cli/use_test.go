package cli

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/arsenal/internal/config"
	"github.com/arsenal/internal/plugin"
	"github.com/arsenal/internal/version"
)

// newUseCmd が正しく作成されるかテストする
func TestNewUseCmd(t *testing.T) {
	cmd := newUseCmd()

	if cmd.Use != "use <tool> <version>" {
		t.Errorf("Use = %q, want %q", cmd.Use, "use <tool> <version>")
	}
}

// runUse が正しく動作するかテストする（成功）
func TestRunUseSuccess(t *testing.T) {
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
	versionDir := filepath.Join(paths.Versions, "node", "20.10.0")
	if err := os.MkdirAll(versionDir, 0755); err != nil {
		t.Fatalf("バージョンディレクトリ作成エラー: %v", err)
	}

	var err error
	registry, err = plugin.NewRegistry(paths)
	if err != nil {
		t.Fatalf("レジストリ作成エラー: %v", err)
	}

	manager = version.NewManager(paths, registry)

	// runUse を実行
	err = runUse("node", "20.10.0", false)
	if err != nil {
		t.Errorf("runUse() エラー: %v", err)
	}

	// symlink が作成されたか確認
	symlinkPath := filepath.Join(paths.Current, "node")
	if _, err := os.Lstat(symlinkPath); err != nil {
		t.Error("symlink が作成されませんでした")
	}
}

// runUse が未インストールバージョンでエラーを返すかテストする
func TestRunUseNotInstalled(t *testing.T) {
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

	// runUse を実行（未インストールバージョン）
	err = runUse("node", "99.99.99", false)
	if err == nil {
		t.Error("未インストールバージョンなのにエラーが返されませんでした")
	}
}

// runUse が --local フラグで .toolversions を更新するかテストする
func TestRunUseWithLocal(t *testing.T) {
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
	versionDir := filepath.Join(paths.Versions, "node", "20.10.0")
	if err := os.MkdirAll(versionDir, 0755); err != nil {
		t.Fatalf("バージョンディレクトリ作成エラー: %v", err)
	}

	var err error
	registry, err = plugin.NewRegistry(paths)
	if err != nil {
		t.Fatalf("レジストリ作成エラー: %v", err)
	}

	manager = version.NewManager(paths, registry)

	// 作業ディレクトリを変更
	originalWd, _ := os.Getwd()
	defer func() { _ = os.Chdir(originalWd) }()
	_ = os.Chdir(tmpDir)

	// runUse を実行（--local フラグ付き）
	err = runUse("node", "20.10.0", true)
	if err != nil {
		t.Errorf("runUse() エラー: %v", err)
	}

	// .toolversions が作成されたか確認
	toolversionsPath := filepath.Join(tmpDir, ".toolversions")
	if _, err := os.Stat(toolversionsPath); os.IsNotExist(err) {
		t.Error(".toolversions が作成されませんでした")
	}

	// .toolversions の内容を確認
	content, err := os.ReadFile(toolversionsPath)
	if err != nil {
		t.Fatalf(".toolversions 読み込みエラー: %v", err)
	}

	expected := "node 20.10.0\n"
	if string(content) != expected {
		t.Errorf(".toolversions の内容が正しくありません\ngot:  %q\nwant: %q", string(content), expected)
	}
}

// runUse が不明なツールでエラーを返すかテストする
func TestRunUseUnknownTool(t *testing.T) {
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

	// runUse を実行（存在しないツール）
	err = runUse("nonexistent", "1.0.0", false)
	if err == nil {
		t.Error("存在しないツールでエラーが返されませんでした")
	}
}
