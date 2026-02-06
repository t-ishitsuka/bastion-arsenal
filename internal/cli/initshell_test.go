package cli

import (
	"bytes"
	"io"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/arsenal/internal/config"
	"github.com/arsenal/internal/plugin"
	"github.com/arsenal/internal/version"
)

// newInitShellCmd が正しく作成されるかテストする
func TestNewInitShellCmd(t *testing.T) {
	cmd := newInitShellCmd()

	if cmd.Use != "init-shell [bash|zsh|fish]" {
		t.Errorf("Use = %q, want %q", cmd.Use, "init-shell [bash|zsh|fish]")
	}
}

// runInitShell が bash スクリプトを生成するかテストする
func TestRunInitShellBash(t *testing.T) {
	// テスト用の環境をセットアップ
	tmpDir := t.TempDir()
	paths = &config.Paths{
		Root:     filepath.Join(tmpDir, "arsenal"),
		Versions: filepath.Join(tmpDir, "arsenal", "versions"),
		Current:  filepath.Join(tmpDir, "arsenal", "current"),
		Plugins:  filepath.Join(tmpDir, "arsenal", "plugins"),
	}

	// 標準出力をキャプチャ
	oldStdout := os.Stdout
	r, w, _ := os.Pipe()
	os.Stdout = w

	// runInitShell を実行
	err := runInitShell("bash")
	if err != nil {
		t.Errorf("runInitShell() エラー: %v", err)
	}

	// 標準出力を復元
	_ = w.Close()
	os.Stdout = oldStdout

	var buf bytes.Buffer
	_, _ = io.Copy(&buf, r)
	output := buf.String()

	// 出力内容を確認
	if !strings.Contains(output, "export PATH=") {
		t.Error("export PATH が含まれていません")
	}
	if !strings.Contains(output, paths.Root) {
		t.Errorf("Arsenal のパスが含まれていません: %s", paths.Root)
	}
	if !strings.Contains(output, "arsenal completion bash") {
		t.Error("補完スクリプトが含まれていません")
	}
}

// runInitShell が zsh スクリプトを生成するかテストする
func TestRunInitShellZsh(t *testing.T) {
	// テスト用の環境をセットアップ
	tmpDir := t.TempDir()
	paths = &config.Paths{
		Root:     filepath.Join(tmpDir, "arsenal"),
		Versions: filepath.Join(tmpDir, "arsenal", "versions"),
		Current:  filepath.Join(tmpDir, "arsenal", "current"),
		Plugins:  filepath.Join(tmpDir, "arsenal", "plugins"),
	}

	// 標準出力をキャプチャ
	oldStdout := os.Stdout
	r, w, _ := os.Pipe()
	os.Stdout = w

	// runInitShell を実行
	err := runInitShell("zsh")
	if err != nil {
		t.Errorf("runInitShell() エラー: %v", err)
	}

	// 標準出力を復元
	_ = w.Close()
	os.Stdout = oldStdout

	var buf bytes.Buffer
	_, _ = io.Copy(&buf, r)
	output := buf.String()

	// 出力内容を確認
	if !strings.Contains(output, "export PATH=") {
		t.Error("export PATH が含まれていません")
	}
	if !strings.Contains(output, "arsenal completion zsh") {
		t.Error("補完スクリプトが含まれていません")
	}
}

// runInitShell が fish スクリプトを生成するかテストする
func TestRunInitShellFish(t *testing.T) {
	// テスト用の環境をセットアップ
	tmpDir := t.TempDir()
	paths = &config.Paths{
		Root:     filepath.Join(tmpDir, "arsenal"),
		Versions: filepath.Join(tmpDir, "arsenal", "versions"),
		Current:  filepath.Join(tmpDir, "arsenal", "current"),
		Plugins:  filepath.Join(tmpDir, "arsenal", "plugins"),
	}

	// 標準出力をキャプチャ
	oldStdout := os.Stdout
	r, w, _ := os.Pipe()
	os.Stdout = w

	// runInitShell を実行
	err := runInitShell("fish")
	if err != nil {
		t.Errorf("runInitShell() エラー: %v", err)
	}

	// 標準出力を復元
	_ = w.Close()
	os.Stdout = oldStdout

	var buf bytes.Buffer
	_, _ = io.Copy(&buf, r)
	output := buf.String()

	// 出力内容を確認
	if !strings.Contains(output, "set -gx PATH") {
		t.Error("set -gx PATH が含まれていません")
	}
	if !strings.Contains(output, "arsenal completion fish") {
		t.Error("補完スクリプトが含まれていません")
	}
}

// runInitShell が不明なシェルでエラーを返すかテストする
func TestRunInitShellUnknownShell(t *testing.T) {
	// テスト用の環境をセットアップ
	tmpDir := t.TempDir()
	paths = &config.Paths{
		Root:     filepath.Join(tmpDir, "arsenal"),
		Versions: filepath.Join(tmpDir, "arsenal", "versions"),
		Current:  filepath.Join(tmpDir, "arsenal", "current"),
		Plugins:  filepath.Join(tmpDir, "arsenal", "plugins"),
	}

	var err error
	registry, err = plugin.NewRegistry(paths)
	if err != nil {
		t.Fatalf("レジストリ作成エラー: %v", err)
	}

	manager = version.NewManager(paths, registry)

	// runInitShell を実行（存在しないシェル）
	err = runInitShell("unknown")
	if err == nil {
		t.Error("存在しないシェルでエラーが返されませんでした")
	}
}
