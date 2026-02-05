package config

import (
	"os"
	"path/filepath"
	"testing"
)

// GetPaths 関数が正しくパスを返すかテストする
func TestGetPaths(t *testing.T) {
	paths, err := GetPaths()
	if err != nil {
		t.Fatalf("GetPaths() エラー: %v", err)
	}

	home, err := os.UserHomeDir()
	if err != nil {
		t.Fatalf("ホームディレクトリ取得エラー: %v", err)
	}

	expectedRoot := filepath.Join(home, ".arsenal")

	if paths.Root != expectedRoot {
		t.Errorf("Root = %q, want %q", paths.Root, expectedRoot)
	}

	if paths.Versions != filepath.Join(expectedRoot, "versions") {
		t.Errorf("Versions パスが正しくありません")
	}

	if paths.Current != filepath.Join(expectedRoot, "current") {
		t.Errorf("Current パスが正しくありません")
	}

	if paths.Plugins != filepath.Join(expectedRoot, "plugins") {
		t.Errorf("Plugins パスが正しくありません")
	}

	if paths.Config != filepath.Join(expectedRoot, ConfigFile) {
		t.Errorf("Config パスが正しくありません")
	}
}

// ディレクトリ作成が正しく動作するかテストする
func TestEnsureDirs(t *testing.T) {
	// 一時ディレクトリを作成
	tmpDir := t.TempDir()

	paths := &Paths{
		Root:     filepath.Join(tmpDir, "arsenal"),
		Versions: filepath.Join(tmpDir, "arsenal", "versions"),
		Current:  filepath.Join(tmpDir, "arsenal", "current"),
		Plugins:  filepath.Join(tmpDir, "arsenal", "plugins"),
		Config:   filepath.Join(tmpDir, "arsenal", "config.toml"),
	}

	// ディレクトリを作成
	if err := paths.EnsureDirs(); err != nil {
		t.Fatalf("EnsureDirs() エラー: %v", err)
	}

	// 各ディレクトリが存在するか確認
	dirs := []string{paths.Root, paths.Versions, paths.Current, paths.Plugins}
	for _, dir := range dirs {
		if _, err := os.Stat(dir); os.IsNotExist(err) {
			t.Errorf("ディレクトリが作成されていません: %s", dir)
		}
	}
}

// 正しいパスを返すかテストする
func TestToolVersionPath(t *testing.T) {
	paths := &Paths{
		Versions: "/home/user/.arsenal/versions",
	}

	tests := []struct {
		tool    string
		version string
		want    string
	}{
		{"node", "20.10.0", "/home/user/.arsenal/versions/node/20.10.0"},
		{"go", "1.22.0", "/home/user/.arsenal/versions/go/1.22.0"},
		{"python", "3.12.0", "/home/user/.arsenal/versions/python/3.12.0"},
	}

	for _, tt := range tests {
		got := paths.ToolVersionPath(tt.tool, tt.version)
		if got != tt.want {
			t.Errorf("ToolVersionPath(%q, %q) = %q, want %q", tt.tool, tt.version, got, tt.want)
		}
	}
}

// 正しいパスを返すかテストする
func TestToolCurrentPath(t *testing.T) {
	paths := &Paths{
		Current: "/home/user/.arsenal/current",
	}

	tests := []struct {
		tool string
		want string
	}{
		{"node", "/home/user/.arsenal/current/node"},
		{"go", "/home/user/.arsenal/current/go"},
		{"python", "/home/user/.arsenal/current/python"},
	}

	for _, tt := range tests {
		got := paths.ToolCurrentPath(tt.tool)
		if got != tt.want {
			t.Errorf("ToolCurrentPath(%q) = %q, want %q", tt.tool, got, tt.want)
		}
	}
}

// 正しいパスを返すかテストする
func TestToolBinPath(t *testing.T) {
	paths := &Paths{
		Current: "/home/user/.arsenal/current",
	}

	tests := []struct {
		tool string
		want string
	}{
		{"node", "/home/user/.arsenal/current/node/bin"},
		{"go", "/home/user/.arsenal/current/go/bin"},
		{"python", "/home/user/.arsenal/current/python/bin"},
	}

	for _, tt := range tests {
		got := paths.ToolBinPath(tt.tool)
		if got != tt.want {
			t.Errorf("ToolBinPath(%q) = %q, want %q", tt.tool, got, tt.want)
		}
	}
}
