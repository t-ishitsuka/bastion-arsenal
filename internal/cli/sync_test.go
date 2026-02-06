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

// newSyncCmd が正しく作成されるかテストする
func TestNewSyncCmd(t *testing.T) {
	cmd := newSyncCmd()

	if cmd.Use != "sync" {
		t.Errorf("Use = %q, want %q", cmd.Use, "sync")
	}
}

// runSync が正しく動作するかテストする（成功）
func TestRunSyncSuccess(t *testing.T) {
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

	// .toolversions ファイルを作成
	toolversionsContent := "testnode 20.10.0\n"
	toolversionsPath := filepath.Join(tmpDir, ".toolversions")
	if err := os.WriteFile(toolversionsPath, []byte(toolversionsContent), 0644); err != nil {
		t.Fatalf(".toolversions 作成エラー: %v", err)
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

	// runSync を実行
	err = runSync()
	if err != nil {
		t.Errorf("runSync() エラー: %v", err)
	}

	// インストールされたか確認
	versionDir := filepath.Join(paths.Versions, "testnode", "20.10.0")
	if _, err := os.Stat(versionDir); os.IsNotExist(err) {
		t.Error("バージョンがインストールされませんでした")
	}

	// symlink が作成されたか確認
	symlinkPath := filepath.Join(paths.Current, "testnode")
	if _, err := os.Lstat(symlinkPath); err != nil {
		t.Error("symlink が作成されませんでした")
	}
}

// runSync が .toolversions がない場合にエラーを返すかテストする
func TestRunSyncNoToolversions(t *testing.T) {
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

	// 作業ディレクトリを変更
	originalWd, _ := os.Getwd()
	defer func() { _ = os.Chdir(originalWd) }()
	_ = os.Chdir(tmpDir)

	// runSync を実行（.toolversions がない）
	err = runSync()
	if err == nil {
		t.Error(".toolversions がないのにエラーが返されませんでした")
	}
}

// runSync が既にインストール済みのバージョンをスキップするかテストする
func TestRunSyncAlreadyInstalled(t *testing.T) {
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

	// 既にインストール済みのバージョンディレクトリを作成
	versionDir := filepath.Join(paths.Versions, "node", "20.10.0")
	if err := os.MkdirAll(versionDir, 0755); err != nil {
		t.Fatalf("バージョンディレクトリ作成エラー: %v", err)
	}

	// .toolversions ファイルを作成
	toolversionsContent := "node 20.10.0\n"
	toolversionsPath := filepath.Join(tmpDir, ".toolversions")
	if err := os.WriteFile(toolversionsPath, []byte(toolversionsContent), 0644); err != nil {
		t.Fatalf(".toolversions 作成エラー: %v", err)
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

	// runSync を実行
	err = runSync()
	if err != nil {
		t.Errorf("runSync() エラー: %v", err)
	}

	// symlink が作成されたか確認
	symlinkPath := filepath.Join(paths.Current, "node")
	if _, err := os.Lstat(symlinkPath); err != nil {
		t.Error("symlink が作成されませんでした")
	}
}
