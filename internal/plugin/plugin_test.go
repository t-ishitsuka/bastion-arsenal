package plugin

import (
	"testing"

	"github.com/arsenal/internal/config"
)

// レジストリが正しく作成されるかテストする
func TestNewRegistry(t *testing.T) {
	tmpDir := t.TempDir()
	paths := &config.Paths{
		Plugins: tmpDir,
	}

	registry, err := NewRegistry(paths)
	if err != nil {
		t.Fatalf("NewRegistry() エラー: %v", err)
	}

	if registry == nil {
		t.Fatal("registry が nil です")
	}

	// 組み込みプラグイン（node.toml）が読み込まれているか確認
	if _, err := registry.Get("node"); err != nil {
		t.Errorf("node プラグインが読み込まれていません: %v", err)
	}
}

// ダウンロード URL が正しく解決されるかテストする
func TestPluginResolveDownloadURL(t *testing.T) {
	plugin := &Plugin{
		DownloadURL: "https://example.com/{{tool}}-{{version}}-{{os}}-{{arch}}.tar.gz",
		OSMap: map[string]string{
			"linux":  "linux",
			"darwin": "macos",
		},
		ArchMap: map[string]string{
			"amd64": "x64",
			"arm64": "arm64",
		},
	}

	// テンプレート変数が置換されることを確認
	url := plugin.ResolveDownloadURL("1.0.0")
	if url == "" {
		t.Error("ResolveDownloadURL が空文字列を返しました")
	}

	// {{version}} が置換されているか確認
	if !contains(url, "1.0.0") {
		t.Errorf("URL にバージョンが含まれていません: %s", url)
	}

	// テンプレート変数が残っていないか確認
	if contains(url, "{{version}}") || contains(url, "{{os}}") || contains(url, "{{arch}}") {
		t.Errorf("テンプレート変数が置換されていません: %s", url)
	}
}

// アーカイブタイプが正しく解決されるかテストする
func TestPluginResolveArchiveType(t *testing.T) {
	tests := []struct {
		name        string
		plugin      *Plugin
		wantContain string // 期待される文字列の一部
	}{
		{
			name:        "指定されたタイプを使用",
			plugin:      &Plugin{ArchiveType: "zip"},
			wantContain: "zip",
		},
		{
			name:        "tar.gz を使用",
			plugin:      &Plugin{ArchiveType: "tar.gz"},
			wantContain: "tar.gz",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.plugin.ResolveArchiveType()
			if got != tt.wantContain {
				t.Errorf("ResolveArchiveType() = %q, want %q", got, tt.wantContain)
			}
		})
	}
}

// プラグインリストが正しく返されるかテストする
func TestRegistryList(t *testing.T) {
	tmpDir := t.TempDir()
	paths := &config.Paths{
		Plugins: tmpDir,
	}

	registry, err := NewRegistry(paths)
	if err != nil {
		t.Fatalf("NewRegistry() エラー: %v", err)
	}

	list := registry.List()
	if len(list) == 0 {
		t.Error("プラグインリストが空です")
	}

	// node プラグインが含まれているか確認
	found := false
	for _, name := range list {
		if name == "node" {
			found = true
			break
		}
	}
	if !found {
		t.Error("node プラグインがリストに含まれていません")
	}
}

// 正しいプラグインが取得できるかテストする
func TestRegistryGet(t *testing.T) {
	tmpDir := t.TempDir()
	paths := &config.Paths{
		Plugins: tmpDir,
	}

	registry, err := NewRegistry(paths)
	if err != nil {
		t.Fatalf("NewRegistry() エラー: %v", err)
	}

	// 存在するプラグイン
	plugin, err := registry.Get("node")
	if err != nil {
		t.Errorf("Get(node) エラー: %v", err)
		return
	}
	if plugin == nil {
		t.Error("plugin が nil です")
		return
	}
	if plugin.Name != "node" {
		t.Errorf("plugin.Name = %q, want %q", plugin.Name, "node")
	}

	// 存在しないプラグイン
	_, err = registry.Get("nonexistent")
	if err == nil {
		t.Error("存在しないプラグインでエラーが返されませんでした")
	}
}

// 文字列に部分文字列が含まれるか確認するヘルパー関数
func contains(s, substr string) bool {
	return len(s) >= len(substr) && (s == substr || len(s) > len(substr) && findSubstring(s, substr))
}

func findSubstring(s, substr string) bool {
	for i := 0; i <= len(s)-len(substr); i++ {
		if s[i:i+len(substr)] == substr {
			return true
		}
	}
	return false
}
