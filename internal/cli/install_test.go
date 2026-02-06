package cli

import (
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"testing"

	"github.com/arsenal/internal/config"
	"github.com/arsenal/internal/plugin"
	"github.com/arsenal/internal/version"
)

// newInstallCmd が正しく作成されるかテストする
func TestNewInstallCmd(t *testing.T) {
	cmd := newInstallCmd()

	if cmd.Use != "install <tool> <version>" {
		t.Errorf("Use = %q, want %q", cmd.Use, "install <tool> <version>")
	}
}

// runInstall が正しく動作するかテストする（インストール成功）
func TestRunInstallSuccess(t *testing.T) {
	// テスト用の HTTP サーバーを起動（ダミーのアーカイブを返す）
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// 空の tar.gz を返す（最小限のヘッダー）
		gzipHeader := []byte{0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0xff}
		tarEnd := []byte{0x03, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00}
		_, _ = w.Write(append(gzipHeader, tarEnd...))
	}))
	defer server.Close()

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

	// テスト用プラグイン定義を作成
	pluginContent := `name = "testnode"
display_name = "Test Node.js"
description = "Test Node.js runtime"
download_url = "` + server.URL + `/node-v{{version}}.tar.gz"
bin_path = "bin"
archive_type = "tar.gz"
`
	pluginPath := filepath.Join(paths.Plugins, "testnode.toml")
	if err := os.WriteFile(pluginPath, []byte(pluginContent), 0644); err != nil {
		t.Fatalf("プラグインファイル作成エラー: %v", err)
	}

	var err error
	registry, err = plugin.NewRegistry(paths)
	if err != nil {
		t.Fatalf("レジストリ作成エラー: %v", err)
	}

	manager = version.NewManager(paths, registry)

	// runInstall を実行
	err = runInstall("testnode", "20.10.0")
	if err != nil {
		t.Errorf("runInstall() エラー: %v", err)
	}

	// インストールディレクトリが作成されたか確認
	installDir := filepath.Join(paths.Versions, "testnode", "20.10.0")
	if _, err := os.Stat(installDir); os.IsNotExist(err) {
		t.Error("インストールディレクトリが作成されませんでした")
	}
}

// runInstall が既にインストール済みの場合にエラーを返すかテストする
func TestRunInstallAlreadyInstalled(t *testing.T) {
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

	// 既にインストール済みのディレクトリを作成
	installDir := filepath.Join(paths.Versions, "node", "20.10.0")
	if err := os.MkdirAll(installDir, 0755); err != nil {
		t.Fatalf("インストールディレクトリ作成エラー: %v", err)
	}

	var err error
	registry, err = plugin.NewRegistry(paths)
	if err != nil {
		t.Fatalf("レジストリ作成エラー: %v", err)
	}

	manager = version.NewManager(paths, registry)

	// runInstall を実行（エラーが返るはず）
	err = runInstall("node", "20.10.0")
	if err == nil {
		t.Error("既にインストール済みなのにエラーが返されませんでした")
	}
}

// runInstall が不明なツールでエラーを返すかテストする
func TestRunInstallUnknownTool(t *testing.T) {
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

	// runInstall を実行（存在しないツール）
	err = runInstall("nonexistent", "1.0.0")
	if err == nil {
		t.Error("存在しないツールでエラーが返されませんでした")
	}
}
