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

// newLsRemoteCmd が正しく作成されるかテストする
func TestNewLsRemoteCmd(t *testing.T) {
	cmd := newLsRemoteCmd()

	if cmd.Use != "ls-remote <tool>" {
		t.Errorf("Use = %q, want %q", cmd.Use, "ls-remote <tool>")
	}
}

// runLsRemote が正しく動作するかテストする（成功）
func TestRunLsRemoteSuccess(t *testing.T) {
	// テスト用の HTTP サーバーを起動（JSON を返す）
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`[
			{"version": "v21.5.0"},
			{"version": "v21.4.0"},
			{"version": "v20.10.0"},
			{"version": "v20.9.0"},
			{"version": "v18.19.0"}
		]`))
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
list_url = "` + server.URL + `"
list_format = "json"
version_prefix = "v"
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

	// runLsRemote を実行
	err = runLsRemote("testnode", 3)
	if err != nil {
		t.Errorf("runLsRemote() エラー: %v", err)
	}
}

// runLsRemote が不明なツールでエラーを返すかテストする
func TestRunLsRemoteUnknownTool(t *testing.T) {
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

	// runLsRemote を実行（存在しないツール）
	err = runLsRemote("nonexistent", 20)
	if err == nil {
		t.Error("存在しないツールでエラーが返されませんでした")
	}
}

// runLsRemote が list_url のないツールでエラーを返すかテストする
func TestRunLsRemoteNoListURL(t *testing.T) {
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

	// テスト用プラグイン定義を作成（list_url なし）
	pluginContent := `name = "testtool"
display_name = "Test Tool"
description = "Test tool without list_url"
`
	pluginPath := filepath.Join(paths.Plugins, "testtool.toml")
	if err := os.WriteFile(pluginPath, []byte(pluginContent), 0644); err != nil {
		t.Fatalf("プラグインファイル作成エラー: %v", err)
	}

	var err error
	registry, err = plugin.NewRegistry(paths)
	if err != nil {
		t.Fatalf("レジストリ作成エラー: %v", err)
	}

	manager = version.NewManager(paths, registry)

	// runLsRemote を実行（list_url なし）
	err = runLsRemote("testtool", 20)
	if err == nil {
		t.Error("list_url がないのにエラーが返されませんでした")
	}
}
