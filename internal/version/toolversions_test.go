package version

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/arsenal/internal/config"
)

// .toolversions ファイルのパースをテストする
func TestParseToolVersionsFile(t *testing.T) {
	// 一時ファイルを作成
	tmpDir := t.TempDir()
	tvFile := filepath.Join(tmpDir, ".toolversions")

	content := `# コメント
node 20.10.0
go 1.22.0

# 別のコメント
python 3.12.0
`
	if err := os.WriteFile(tvFile, []byte(content), 0644); err != nil {
		t.Fatalf("テストファイル作成エラー: %v", err)
	}

	tv, err := parseToolVersionsFile(tvFile)
	if err != nil {
		t.Fatalf("parseToolVersionsFile() エラー: %v", err)
	}

	// 期待されるツールとバージョン
	expected := map[string]string{
		"node":   "20.10.0",
		"go":     "1.22.0",
		"python": "3.12.0",
	}

	if len(tv.Tools) != len(expected) {
		t.Errorf("ツール数 = %d, want %d", len(tv.Tools), len(expected))
	}

	for tool, version := range expected {
		if tv.Tools[tool] != version {
			t.Errorf("Tools[%q] = %q, want %q", tool, tv.Tools[tool], version)
		}
	}
}

// 不正なフォーマットでエラーが返されるかテストする
func TestParseToolVersionsFileInvalid(t *testing.T) {
	tests := []struct {
		name    string
		content string
	}{
		{
			name:    "バージョンなし",
			content: "node\n",
		},
		{
			name:    "余分なフィールド",
			content: "node 20.10.0 extra\n",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tmpDir := t.TempDir()
			tvFile := filepath.Join(tmpDir, ".toolversions")

			if err := os.WriteFile(tvFile, []byte(tt.content), 0644); err != nil {
				t.Fatalf("テストファイル作成エラー: %v", err)
			}

			_, err := parseToolVersionsFile(tvFile)
			if err == nil {
				t.Error("不正なフォーマットでエラーが返されませんでした")
			}
		})
	}
}

// .toolversions ファイルの書き込みをテストする
func TestWriteToolVersions(t *testing.T) {
	tmpDir := t.TempDir()

	tv := &ToolVersions{
		Tools: map[string]string{
			"node":   "20.10.0",
			"go":     "1.22.0",
			"python": "3.12.0",
		},
	}

	if err := WriteToolVersions(tmpDir, tv); err != nil {
		t.Fatalf("WriteToolVersions() エラー: %v", err)
	}

	// ファイルが作成されているか確認
	tvFile := filepath.Join(tmpDir, config.ToolVersionFile)
	if _, err := os.Stat(tvFile); os.IsNotExist(err) {
		t.Fatalf(".toolversions ファイルが作成されていません")
	}

	// ファイルを読み込んで検証
	parsed, err := parseToolVersionsFile(tvFile)
	if err != nil {
		t.Fatalf("書き込んだファイルの読み込みエラー: %v", err)
	}

	for tool, version := range tv.Tools {
		if parsed.Tools[tool] != version {
			t.Errorf("Tools[%q] = %q, want %q", tool, parsed.Tools[tool], version)
		}
	}
}

// .toolversions ファイルの検索をテストする
func TestFindToolVersionsFile(t *testing.T) {
	// ディレクトリ構造を作成
	tmpDir := t.TempDir()
	subDir := filepath.Join(tmpDir, "project", "src")
	if err := os.MkdirAll(subDir, 0755); err != nil {
		t.Fatalf("ディレクトリ作成エラー: %v", err)
	}

	// ルートに .toolversions を作成
	tvFile := filepath.Join(tmpDir, "project", ".toolversions")
	if err := os.WriteFile(tvFile, []byte("node 20.10.0\n"), 0644); err != nil {
		t.Fatalf("テストファイル作成エラー: %v", err)
	}

	// サブディレクトリから検索
	found, err := findToolVersionsFile(subDir)
	if err != nil {
		t.Fatalf("findToolVersionsFile() エラー: %v", err)
	}

	if found != tvFile {
		t.Errorf("findToolVersionsFile() = %q, want %q", found, tvFile)
	}
}

// .toolversions が見つからない場合をテストする
func TestFindToolVersionsFileNotFound(t *testing.T) {
	tmpDir := t.TempDir()

	_, err := findToolVersionsFile(tmpDir)
	if err == nil {
		t.Error(".toolversions が存在しないのにエラーが返されませんでした")
	}
}

// ReadToolVersions 関数をテストする
func TestReadToolVersions(t *testing.T) {
	tmpDir := t.TempDir()
	tvFile := filepath.Join(tmpDir, ".toolversions")

	content := "node 20.10.0\ngo 1.22.0\n"
	if err := os.WriteFile(tvFile, []byte(content), 0644); err != nil {
		t.Fatalf("テストファイル作成エラー: %v", err)
	}

	tv, path, err := ReadToolVersions(tmpDir)
	if err != nil {
		t.Fatalf("ReadToolVersions() エラー: %v", err)
	}

	if path != tvFile {
		t.Errorf("path = %q, want %q", path, tvFile)
	}

	if tv.Tools["node"] != "20.10.0" {
		t.Errorf("node バージョンが正しくありません")
	}

	if tv.Tools["go"] != "1.22.0" {
		t.Errorf("go バージョンが正しくありません")
	}
}
