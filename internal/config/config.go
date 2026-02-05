package config

import (
	"os"
	"path/filepath"
)

const (
	AppName         = "arsenal"
	ToolVersionFile = ".toolversions"
	ConfigFile      = "config.toml"
)

// Arsenal の重要なディレクトリパスを保持する
type Paths struct {
	Root     string // ~/.arsenal
	Versions string // ~/.arsenal/versions
	Current  string // ~/.arsenal/current
	Plugins  string // ~/.arsenal/plugins
	Config   string // ~/.arsenal/config.toml
}

// グローバル設定を保持する
type Config struct {
	DefaultShell string `toml:"default_shell"`
	AutoSync     bool   `toml:"auto_sync"`
}

// 標準的なディレクトリレイアウトを返す
func GetPaths() (*Paths, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return nil, err
	}

	root := filepath.Join(home, ".arsenal")
	return &Paths{
		Root:     root,
		Versions: filepath.Join(root, "versions"),
		Current:  filepath.Join(root, "current"),
		Plugins:  filepath.Join(root, "plugins"),
		Config:   filepath.Join(root, ConfigFile),
	}, nil
}

// 必要な全ディレクトリを作成する
func (p *Paths) EnsureDirs() error {
	dirs := []string{p.Root, p.Versions, p.Current, p.Plugins}
	for _, d := range dirs {
		if err := os.MkdirAll(d, 0755); err != nil {
			return err
		}
	}
	return nil
}

// ツールのバージョンディレクトリのパスを返す
// 例: ~/.arsenal/versions/node/20.10.0
func (p *Paths) ToolVersionPath(tool, version string) string {
	return filepath.Join(p.Versions, tool, version)
}

// ツールの symlink パスを返す
// 例: ~/.arsenal/current/node
func (p *Paths) ToolCurrentPath(tool string) string {
	return filepath.Join(p.Current, tool)
}

// ツールバージョンの bin ディレクトリを返す
// 例: ~/.arsenal/current/node/bin
func (p *Paths) ToolBinPath(tool string) string {
	return filepath.Join(p.Current, tool, "bin")
}
