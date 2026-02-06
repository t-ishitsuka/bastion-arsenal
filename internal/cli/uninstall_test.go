package cli

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/arsenal/internal/config"
	"github.com/arsenal/internal/plugin"
	"github.com/arsenal/internal/version"
)

// newUninstallCmd が正しく作成されるかテストする
func TestNewUninstallCmd(t *testing.T) {
	cmd := newUninstallCmd()

	if cmd.Use != "uninstall <tool> <version>" {
		t.Errorf("Use = %q, want %q", cmd.Use, "uninstall <tool> <version>")
	}
}

// runUninstall が正しく動作するかテストする（成功）
func TestRunUninstallSuccess(t *testing.T) {
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
	versions := []string{"18.19.0", "20.10.0"}
	for _, v := range versions {
		versionDir := filepath.Join(paths.Versions, "node", v)
		if err := os.MkdirAll(versionDir, 0755); err != nil {
			t.Fatalf("バージョンディレクトリ作成エラー: %v", err)
		}
	}

	var err error
	registry, err = plugin.NewRegistry(paths)
	if err != nil {
		t.Fatalf("レジストリ作成エラー: %v", err)
	}

	manager = version.NewManager(paths, registry)

	// runUninstall を実行
	err = runUninstall("node", "18.19.0")
	if err != nil {
		t.Errorf("runUninstall() エラー: %v", err)
	}

	// バージョンディレクトリが削除されたか確認
	versionDir := filepath.Join(paths.Versions, "node", "18.19.0")
	if _, err := os.Stat(versionDir); !os.IsNotExist(err) {
		t.Error("バージョンディレクトリが削除されませんでした")
	}

	// 他のバージョンが残っているか確認
	otherVersionDir := filepath.Join(paths.Versions, "node", "20.10.0")
	if _, err := os.Stat(otherVersionDir); os.IsNotExist(err) {
		t.Error("他のバージョンが削除されてしまいました")
	}
}

// runUninstall が現在アクティブなバージョンを削除し、他のバージョンに自動切り替えするかテストする
func TestRunUninstallCurrentVersionWithAutoSwitch(t *testing.T) {
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

	// ダミーのバージョンディレクトリを複数作成
	versions := []string{"18.19.0", "20.10.0", "21.0.0"}
	for _, v := range versions {
		versionDir := filepath.Join(paths.Versions, "node", v)
		if err := os.MkdirAll(versionDir, 0755); err != nil {
			t.Fatalf("バージョンディレクトリ作成エラー: %v", err)
		}
	}

	// 現在のバージョンを 21.0.0 に設定（symlink）
	currentVersionDir := filepath.Join(paths.Versions, "node", "21.0.0")
	symlinkPath := filepath.Join(paths.Current, "node")
	if err := os.Symlink(currentVersionDir, symlinkPath); err != nil {
		t.Fatalf("symlink 作成エラー: %v", err)
	}

	var err error
	registry, err = plugin.NewRegistry(paths)
	if err != nil {
		t.Fatalf("レジストリ作成エラー: %v", err)
	}

	manager = version.NewManager(paths, registry)

	// runUninstall を実行（現在のバージョン 21.0.0 を削除）
	err = runUninstall("node", "21.0.0")
	if err != nil {
		t.Errorf("runUninstall() エラー: %v", err)
	}

	// 21.0.0 のバージョンディレクトリが削除されたか確認
	if _, err := os.Stat(currentVersionDir); !os.IsNotExist(err) {
		t.Error("バージョンディレクトリが削除されませんでした")
	}

	// symlink が 20.10.0 に自動切り替えされたか確認
	link, err := os.Readlink(symlinkPath)
	if err != nil {
		t.Fatalf("symlink 読み込みエラー: %v", err)
	}

	expectedLink := filepath.Join(paths.Versions, "node", "20.10.0")
	if link != expectedLink {
		t.Errorf("symlink が正しく切り替わりませんでした\ngot:  %s\nwant: %s", link, expectedLink)
	}
}

// runUninstall が最後のバージョンを削除した時に symlink も削除するかテストする
func TestRunUninstallLastVersion(t *testing.T) {
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

	// ダミーのバージョンディレクトリを作成（1つだけ）
	versionDir := filepath.Join(paths.Versions, "node", "20.10.0")
	if err := os.MkdirAll(versionDir, 0755); err != nil {
		t.Fatalf("バージョンディレクトリ作成エラー: %v", err)
	}

	// 現在のバージョンを設定（symlink）
	symlinkPath := filepath.Join(paths.Current, "node")
	if err := os.Symlink(versionDir, symlinkPath); err != nil {
		t.Fatalf("symlink 作成エラー: %v", err)
	}

	var err error
	registry, err = plugin.NewRegistry(paths)
	if err != nil {
		t.Fatalf("レジストリ作成エラー: %v", err)
	}

	manager = version.NewManager(paths, registry)

	// runUninstall を実行（最後のバージョン）
	err = runUninstall("node", "20.10.0")
	if err != nil {
		t.Errorf("runUninstall() エラー: %v", err)
	}

	// バージョンディレクトリが削除されたか確認
	if _, err := os.Stat(versionDir); !os.IsNotExist(err) {
		t.Error("バージョンディレクトリが削除されませんでした")
	}

	// symlink も削除されたか確認
	if _, err := os.Lstat(symlinkPath); !os.IsNotExist(err) {
		t.Error("symlink が削除されませんでした")
	}
}

// runUninstall が未インストールバージョンでエラーを返すかテストする
func TestRunUninstallNotInstalled(t *testing.T) {
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

	// runUninstall を実行（未インストールバージョン）
	err = runUninstall("node", "99.99.99")
	if err == nil {
		t.Error("未インストールバージョンなのにエラーが返されませんでした")
	}
}

// runUninstall が不明なツールでエラーを返すかテストする
func TestRunUninstallUnknownTool(t *testing.T) {
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

	// runUninstall を実行（存在しないツール）
	err = runUninstall("nonexistent", "1.0.0")
	if err == nil {
		t.Error("存在しないツールでエラーが返されませんでした")
	}
}
