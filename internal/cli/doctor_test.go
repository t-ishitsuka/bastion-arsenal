package cli

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/arsenal/internal/config"
	"github.com/arsenal/internal/plugin"
	"github.com/arsenal/internal/version"
)

// newDoctorCmd が正しく作成されるかテストする
func TestNewDoctorCmd(t *testing.T) {
	cmd := newDoctorCmd()

	if cmd.Use != "doctor" {
		t.Errorf("Use = %q, want %q", cmd.Use, "doctor")
	}
}

// runDoctor が正しく動作するかテストする（正常な環境）
func TestRunDoctorHealthy(t *testing.T) {
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

	// PATH に current ディレクトリを追加（警告を避けるため）
	oldPath := os.Getenv("PATH")
	defer func() { _ = os.Setenv("PATH", oldPath) }()
	_ = os.Setenv("PATH", paths.Current+":"+oldPath)

	// runDoctor を実行
	err = runDoctor()
	if err != nil {
		t.Errorf("runDoctor() エラー: %v", err)
	}
}

// runDoctor が正しく動作するかテストする（ディレクトリがない環境）
func TestRunDoctorMissingDirs(t *testing.T) {
	// テスト用の環境をセットアップ（ディレクトリを作成しない）
	tmpDir := t.TempDir()
	paths := &config.Paths{
		Root:     filepath.Join(tmpDir, "arsenal"),
		Versions: filepath.Join(tmpDir, "arsenal", "versions"),
		Current:  filepath.Join(tmpDir, "arsenal", "current"),
		Plugins:  filepath.Join(tmpDir, "arsenal", "plugins"),
	}

	// ディレクトリを作成しない（エラーが発生するはず）
	var err error
	registry, err = plugin.NewRegistry(paths)
	if err != nil {
		t.Fatalf("レジストリ作成エラー: %v", err)
	}

	manager = version.NewManager(paths, registry)

	// runDoctor を実行（エラーが返るはず）
	err = runDoctor()
	if err == nil {
		t.Error("ディレクトリが存在しないのにエラーが返されませんでした")
	}
}

// runDoctor が PATH 警告を表示するかテストする
func TestRunDoctorPathWarning(t *testing.T) {
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

	// PATH から current ディレクトリを削除（警告が出るはず）
	oldPath := os.Getenv("PATH")
	defer func() { _ = os.Setenv("PATH", oldPath) }()
	_ = os.Setenv("PATH", "/usr/bin:/bin")

	// runDoctor を実行（警告が出るが成功するはず）
	err = runDoctor()
	if err != nil {
		t.Errorf("runDoctor() エラー: %v", err)
	}
}
