package version

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/arsenal/internal/config"
)

// .toolversions ファイルの内容を表す
type ToolVersions struct {
	Tools map[string]string // tool -> version
}

// 指定ディレクトリから .toolversions ファイルを読み込むか、
// 上位ディレクトリを辿って探す
func ReadToolVersions(dir string) (*ToolVersions, string, error) {
	path, err := findToolVersionsFile(dir)
	if err != nil {
		return nil, "", err
	}

	tv, err := parseToolVersionsFile(path)
	if err != nil {
		return nil, "", err
	}

	return tv, path, nil
}

// .toolversions ファイルを書き込む
func WriteToolVersions(dir string, tv *ToolVersions) error {
	path := filepath.Join(dir, config.ToolVersionFile)
	f, err := os.Create(path)
	if err != nil {
		return err
	}
	defer f.Close()

	for tool, version := range tv.Tools {
		fmt.Fprintf(f, "%s %s\n", tool, version)
	}

	return nil
}

// .toolversions で指定された全バージョンをインストールして切り替える
func (m *Manager) Sync(dir string) error {
	tv, path, err := ReadToolVersions(dir)
	if err != nil {
		return fmt.Errorf(".toolversions 読み込みエラー: %w", err)
	}

	fmt.Printf("📋 %s から同期中\n", path)

	for tool, version := range tv.Tools {
		fmt.Printf("\n── %s %s ──\n", tool, version)

		// インストール済みか確認
		versionDir := m.paths.ToolVersionPath(tool, version)
		if _, err := os.Stat(versionDir); os.IsNotExist(err) {
			// インストール
			if err := m.Install(tool, version); err != nil {
				fmt.Printf("⚠️  %s %s のインストールに失敗: %v\n", tool, version, err)
				continue
			}
		} else {
			fmt.Printf("   既にインストール済み\n")
		}

		// このバージョンに切り替え
		if err := m.Use(tool, version); err != nil {
			fmt.Printf("⚠️  %s を %s に切り替えるのに失敗: %v\n", tool, version, err)
			continue
		}
	}

	fmt.Printf("\n✅ 同期完了\n")
	return nil
}

// ディレクトリツリーを上に辿って .toolversions を探す
func findToolVersionsFile(dir string) (string, error) {
	absDir, err := filepath.Abs(dir)
	if err != nil {
		return "", err
	}

	for {
		candidate := filepath.Join(absDir, config.ToolVersionFile)
		if _, err := os.Stat(candidate); err == nil {
			return candidate, nil
		}

		parent := filepath.Dir(absDir)
		if parent == absDir {
			break // ルートに到達
		}
		absDir = parent
	}

	return "", fmt.Errorf("%s が見つかりません (%s から / まで検索)", config.ToolVersionFile, dir)
}

// ファイルフォーマットを読み込む:
//
//	node 20.10.0
//	go 1.22.0
//	python 3.12.0
func parseToolVersionsFile(path string) (*ToolVersions, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	tv := &ToolVersions{
		Tools: make(map[string]string),
	}

	scanner := bufio.NewScanner(f)
	lineNum := 0
	for scanner.Scan() {
		lineNum++
		line := strings.TrimSpace(scanner.Text())

		// 空行とコメントをスキップ
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}

		parts := strings.Fields(line)
		if len(parts) != 2 {
			return nil, fmt.Errorf("%s:%d: '<ツール> <バージョン>' を期待、'%s' を取得", path, lineNum, line)
		}

		tv.Tools[parts[0]] = parts[1]
	}

	return tv, scanner.Err()
}
