package plugin

import (
	"embed"
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	"strings"

	"github.com/BurntSushi/toml"
	"github.com/arsenal/internal/config"
)

// 特定ツールのインストールと管理方法を定義する
type Plugin struct {
	Name        string `toml:"name"`
	DisplayName string `toml:"display_name"`
	Description string `toml:"description"`

	// バージョン検出とダウンロード用 URL
	ListURL     string `toml:"list_url"`
	ListFormat  string `toml:"list_format"` // "json", "html", "github"
	DownloadURL string `toml:"download_url"`

	// 展開されたアーカイブ内でバイナリが配置されているパス
	BinPath string `toml:"bin_path"`

	// 展開方法: "tar.gz", "tar.xz", "zip"
	ArchiveType string `toml:"archive_type"`

	// リストからのバージョン抽出
	VersionPrefix string `toml:"version_prefix"` // 例: "v" を削除
	VersionRegex  string `toml:"version_regex"`

	// OS/Arch マッピング
	OSMap   map[string]string `toml:"os_map"`
	ArchMap map[string]string `toml:"arch_map"`

	// インストール後コマンド
	PostInstall []string `toml:"post_install"`

	// 設定する環境変数
	EnvVars map[string]string `toml:"env_vars"`
}

//go:embed builtin
var builtinPlugins embed.FS

// 利用可能なプラグインを管理する
type Registry struct {
	plugins map[string]*Plugin
}

// 新しいプラグインレジストリを作成し、組み込みプラグインを読み込む
func NewRegistry(paths *config.Paths) (*Registry, error) {
	r := &Registry{
		plugins: make(map[string]*Plugin),
	}

	// 埋め込みファイルから組み込みプラグインを読み込み
	if err := r.loadBuiltinPlugins(); err != nil {
		return nil, fmt.Errorf("組み込みプラグイン読み込みエラー: %w", err)
	}

	// ユーザープラグインを読み込み（組み込みを上書き）
	if err := r.loadUserPlugins(paths.Plugins); err != nil {
		return nil, fmt.Errorf("ユーザープラグイン読み込みエラー: %w", err)
	}

	return r, nil
}

func (r *Registry) loadBuiltinPlugins() error {
	entries, err := builtinPlugins.ReadDir("builtin")
	if err != nil {
		return err
	}

	for _, entry := range entries {
		if entry.IsDir() || !strings.HasSuffix(entry.Name(), ".toml") {
			continue
		}

		data, err := builtinPlugins.ReadFile(filepath.Join("builtin", entry.Name()))
		if err != nil {
			return fmt.Errorf("reading %s: %w", entry.Name(), err)
		}

		var p Plugin
		if err := toml.Unmarshal(data, &p); err != nil {
			return fmt.Errorf("parsing %s: %w", entry.Name(), err)
		}

		r.plugins[p.Name] = &p
	}

	return nil
}

func (r *Registry) loadUserPlugins(dir string) error {
	entries, err := os.ReadDir(dir)
	if err != nil {
		if os.IsNotExist(err) {
			return nil
		}
		return err
	}

	for _, entry := range entries {
		if entry.IsDir() || !strings.HasSuffix(entry.Name(), ".toml") {
			continue
		}

		var p Plugin
		if _, err := toml.DecodeFile(filepath.Join(dir, entry.Name()), &p); err != nil {
			return fmt.Errorf("parsing %s: %w", entry.Name(), err)
		}

		r.plugins[p.Name] = &p
	}

	return nil
}

// 名前でプラグインを返す
func (r *Registry) Get(name string) (*Plugin, error) {
	p, ok := r.plugins[name]
	if !ok {
		return nil, fmt.Errorf("不明なツール: %s ('arsenal plugin list' で利用可能なツールを確認)", name)
	}
	return p, nil
}

// 利用可能な全プラグイン名を返す
func (r *Registry) List() []string {
	names := make([]string, 0, len(r.plugins))
	for name := range r.plugins {
		names = append(names, name)
	}
	return names
}

// 全プラグインを返す
func (r *Registry) All() map[string]*Plugin {
	return r.plugins
}

// ダウンロード URL 内のテンプレート変数を置換する
func (p *Plugin) ResolveDownloadURL(version string) string {
	url := p.DownloadURL

	osName := runtime.GOOS
	archName := runtime.GOARCH

	// OS マッピングを適用
	if mapped, ok := p.OSMap[osName]; ok {
		osName = mapped
	}

	// Arch マッピングを適用
	if mapped, ok := p.ArchMap[archName]; ok {
		archName = mapped
	}

	replacer := strings.NewReplacer(
		"{{version}}", version,
		"{{os}}", osName,
		"{{arch}}", archName,
	)

	return replacer.Replace(url)
}

// 現在のプラットフォーム用のアーカイブタイプを返す
func (p *Plugin) ResolveArchiveType() string {
	if p.ArchiveType != "" {
		return p.ArchiveType
	}
	if runtime.GOOS == "windows" {
		return "zip"
	}
	return "tar.gz"
}
