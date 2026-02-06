package version

import (
	"archive/tar"
	"archive/zip"
	"compress/gzip"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"sync"
	"time"

	"github.com/arsenal/internal/config"
	"github.com/arsenal/internal/plugin"
)

// バージョンのインストールと切り替えを処理する
type Manager struct {
	paths    *config.Paths
	registry *plugin.Registry
}

// 新しいバージョンマネージャーを作成する
func NewManager(paths *config.Paths, registry *plugin.Registry) *Manager {
	return &Manager{
		paths:    paths,
		registry: registry,
	}
}

// ツールの特定バージョンをダウンロードしてインストールする
func (m *Manager) Install(toolName, version string) error {
	p, err := m.registry.Get(toolName)
	if err != nil {
		return err
	}

	installDir := m.paths.ToolVersionPath(toolName, version)

	// 既にインストール済みか確認
	if _, err := os.Stat(installDir); err == nil {
		return fmt.Errorf("%s %s は既にインストール済みです", toolName, version)
	}

	// インストールディレクトリを作成
	if err := os.MkdirAll(installDir, 0755); err != nil {
		return fmt.Errorf("インストールディレクトリ作成エラー: %w", err)
	}

	// ダウンロード URL を解決
	url := p.ResolveDownloadURL(version)
	fmt.Printf("📦 %s %s をダウンロード中...\n", p.DisplayName, version)
	fmt.Printf("   %s\n", url)

	// ダウンロード
	tmpFile, err := m.download(url)
	if err != nil {
		_ = os.RemoveAll(installDir)
		return fmt.Errorf("ダウンロードエラー: %w", err)
	}
	defer func() { _ = os.Remove(tmpFile) }()

	// 展開
	fmt.Printf("📂 展開中...\n")
	archiveType := p.ResolveArchiveType()
	if err := m.extract(tmpFile, installDir, archiveType); err != nil {
		_ = os.RemoveAll(installDir)
		return fmt.Errorf("展開エラー: %w", err)
	}

	// インストール後コマンドを実行
	if len(p.PostInstall) > 0 {
		fmt.Printf("🔧 インストール後処理を実行中...\n")
		if err := m.runPostInstall(p, installDir); err != nil {
			_ = os.RemoveAll(installDir)
			return fmt.Errorf("インストール後処理エラー: %w", err)
		}
	}

	fmt.Printf("✅ %s %s のインストールが完了しました\n", p.DisplayName, version)
	return nil
}

// symlink を更新してツールのアクティブバージョンを切り替える
func (m *Manager) Use(toolName, version string) error {
	p, err := m.registry.Get(toolName)
	if err != nil {
		return err
	}

	versionDir := m.paths.ToolVersionPath(toolName, version)

	// バージョンがインストール済みか確認
	if _, err := os.Stat(versionDir); os.IsNotExist(err) {
		return fmt.Errorf("%s %s はインストールされていません ('arsenal install %s %s' を実行)",
			toolName, version, toolName, version)
	}

	// 既存の symlink を削除
	symlinkPath := m.paths.ToolCurrentPath(toolName)
	_ = os.Remove(symlinkPath)

	// 新しい symlink を作成
	if err := os.Symlink(versionDir, symlinkPath); err != nil {
		return fmt.Errorf("symlink 作成エラー: %w", err)
	}

	fmt.Printf("✅ %s %s に切り替えました\n", p.DisplayName, version)
	return nil
}

// インストール済みバージョンを削除する
func (m *Manager) Uninstall(toolName, version string) error {
	p, err := m.registry.Get(toolName)
	if err != nil {
		return err
	}

	versionDir := m.paths.ToolVersionPath(toolName, version)
	if _, err := os.Stat(versionDir); os.IsNotExist(err) {
		return fmt.Errorf("%s %s はインストールされていません", toolName, version)
	}

	// 現在のバージョンか確認
	current, _ := m.Current(toolName)
	if current == version {
		// 先に symlink を削除
		_ = os.Remove(m.paths.ToolCurrentPath(toolName))
	}

	if err := os.RemoveAll(versionDir); err != nil {
		return fmt.Errorf("削除エラー: %w", err)
	}

	fmt.Printf("🗑️  %s %s をアンインストールしました\n", p.DisplayName, version)
	return nil
}

// ツールのインストール済みバージョンを返す
func (m *Manager) List(toolName string) ([]string, error) {
	if _, err := m.registry.Get(toolName); err != nil {
		return nil, err
	}

	toolDir := filepath.Join(m.paths.Versions, toolName)
	entries, err := os.ReadDir(toolDir)
	if err != nil {
		if os.IsNotExist(err) {
			return []string{}, nil
		}
		return nil, err
	}

	versions := make([]string, 0, len(entries))
	for _, e := range entries {
		if e.IsDir() {
			versions = append(versions, e.Name())
		}
	}

	sort.Strings(versions)
	return versions, nil
}

// ツールの現在アクティブなバージョンを返す
func (m *Manager) Current(toolName string) (string, error) {
	symlinkPath := m.paths.ToolCurrentPath(toolName)
	target, err := os.Readlink(symlinkPath)
	if err != nil {
		return "", nil // バージョン未設定
	}

	return filepath.Base(target), nil
}

// 現在アクティブな全ツールバージョンを返す
func (m *Manager) CurrentAll() (map[string]string, error) {
	result := make(map[string]string)

	entries, err := os.ReadDir(m.paths.Current)
	if err != nil {
		if os.IsNotExist(err) {
			return result, nil
		}
		return nil, err
	}

	for _, e := range entries {
		target, err := os.Readlink(filepath.Join(m.paths.Current, e.Name()))
		if err != nil {
			continue
		}
		result[e.Name()] = filepath.Base(target)
	}

	return result, nil
}

// Arsenal インストールの健全性をチェックする
func (m *Manager) Doctor() []DiagResult {
	var results []DiagResult

	// ディレクトリをチェック
	results = append(results, m.checkDir("Arsenal ルート", m.paths.Root))
	results = append(results, m.checkDir("バージョンディレクトリ", m.paths.Versions))
	results = append(results, m.checkDir("カレントディレクトリ", m.paths.Current))

	// PATH をチェック
	results = append(results, m.checkPATH())

	// インストール済みツールをチェック
	currentAll, _ := m.CurrentAll()
	for tool, ver := range currentAll {
		results = append(results, DiagResult{
			Name:    fmt.Sprintf("%s バージョン", tool),
			Status:  StatusOK,
			Message: ver,
		})
	}

	return results
}

// 診断チェックのステータスを表す
type DiagStatus int

const (
	StatusOK DiagStatus = iota
	StatusWarn
	StatusError
)

// 単一の診断結果を表す
type DiagResult struct {
	Name    string
	Status  DiagStatus
	Message string
}

func (m *Manager) checkDir(name, path string) DiagResult {
	if _, err := os.Stat(path); os.IsNotExist(err) {
		return DiagResult{Name: name, Status: StatusError, Message: "見つかりません"}
	}
	return DiagResult{Name: name, Status: StatusOK, Message: path}
}

func (m *Manager) checkPATH() DiagResult {
	pathEnv := os.Getenv("PATH")
	currentDir := m.paths.Current
	if strings.Contains(pathEnv, currentDir) {
		return DiagResult{Name: "PATH", Status: StatusOK, Message: "arsenal/current が PATH に含まれています"}
	}
	return DiagResult{
		Name:    "PATH",
		Status:  StatusWarn,
		Message: fmt.Sprintf("PATH に追加: export PATH=\"%s/**/bin:$PATH\" ('arsenal init-shell' 参照)", currentDir),
	}
}

// URL から一時ファイルにダウンロードする
func (m *Manager) download(url string) (string, error) {
	resp, err := http.Get(url)
	if err != nil {
		return "", err
	}
	defer func() { _ = resp.Body.Close() }()

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("HTTP %d: %s", resp.StatusCode, resp.Status)
	}

	tmpFile, err := os.CreateTemp("", "arsenal-download-*")
	if err != nil {
		return "", err
	}
	defer func() { _ = tmpFile.Close() }()

	// Content-Length から総ファイルサイズを取得
	totalSize := resp.ContentLength

	// プログレスバー付きでダウンロード
	if totalSize > 0 {
		pw := &progressWriter{
			total:     totalSize,
			startTime: time.Now(),
		}
		reader := io.TeeReader(resp.Body, pw)

		// 進捗表示用のゴルーチン
		done := make(chan bool)
		go func() {
			ticker := time.NewTicker(100 * time.Millisecond)
			defer ticker.Stop()
			for {
				select {
				case <-done:
					return
				case <-ticker.C:
					pw.printProgress()
				}
			}
		}()

		if _, err := io.Copy(tmpFile, reader); err != nil {
			done <- true
			_ = os.Remove(tmpFile.Name())
			return "", err
		}

		done <- true
		pw.printComplete()
	} else {
		// Content-Length がない場合は通常のコピー
		if _, err := io.Copy(tmpFile, resp.Body); err != nil {
			_ = os.Remove(tmpFile.Name())
			return "", err
		}
	}

	return tmpFile.Name(), nil
}

// プログレスバー用のライター
type progressWriter struct {
	total     int64
	current   int64
	startTime time.Time
	mu        sync.Mutex
}

func (pw *progressWriter) Write(p []byte) (int, error) {
	n := len(p)
	pw.mu.Lock()
	pw.current += int64(n)
	pw.mu.Unlock()
	return n, nil
}

func (pw *progressWriter) printProgress() {
	pw.mu.Lock()
	current := pw.current
	total := pw.total
	pw.mu.Unlock()

	if total <= 0 {
		return
	}

	percent := float64(current) / float64(total) * 100
	currentMB := float64(current) / (1024 * 1024)
	totalMB := float64(total) / (1024 * 1024)

	// 同じ行を上書き
	fmt.Printf("\r   \x1b[36mダウンロード中... %.1f MB / %.1f MB (%.0f%%)\x1b[0m", currentMB, totalMB, percent)
}

func (pw *progressWriter) printComplete() {
	pw.mu.Lock()
	total := pw.total
	pw.mu.Unlock()

	totalMB := float64(total) / (1024 * 1024)
	fmt.Printf("\r   \x1b[32mダウンロード完了 (%.1f MB)\x1b[0m\n", totalMB)
}

// アーカイブを対象ディレクトリに展開する
func (m *Manager) extract(archivePath, targetDir, archiveType string) error {
	switch archiveType {
	case "tar.gz", "tgz":
		return m.extractTarGz(archivePath, targetDir)
	case "tar.xz":
		return m.extractTarGz(archivePath, targetDir) // TODO: xz サポート
	case "zip":
		return m.extractZip(archivePath, targetDir)
	default:
		return fmt.Errorf("サポートされていないアーカイブ形式: %s", archiveType)
	}
}

func (m *Manager) extractTarGz(archivePath, targetDir string) error {
	f, err := os.Open(archivePath)
	if err != nil {
		return err
	}
	defer func() { _ = f.Close() }()

	gzr, err := gzip.NewReader(f)
	if err != nil {
		return err
	}
	defer func() { _ = gzr.Close() }()

	tr := tar.NewReader(gzr)

	// トップレベルディレクトリを検出して削除
	stripPrefix := ""

	for {
		header, err := tr.Next()
		if err == io.EOF {
			break
		}
		if err != nil {
			return err
		}

		// 最初のエントリからトップレベルディレクトリを検出
		if stripPrefix == "" {
			parts := strings.SplitN(header.Name, "/", 2)
			if len(parts) > 1 {
				stripPrefix = parts[0] + "/"
			}
		}

		// トップレベルディレクトリを削除
		name := strings.TrimPrefix(header.Name, stripPrefix)
		if name == "" || name == "." {
			continue
		}

		target := filepath.Join(targetDir, name)

		switch header.Typeflag {
		case tar.TypeDir:
			if err := os.MkdirAll(target, 0755); err != nil {
				return err
			}
		case tar.TypeReg:
			if err := os.MkdirAll(filepath.Dir(target), 0755); err != nil {
				return err
			}
			outFile, err := os.OpenFile(target, os.O_CREATE|os.O_WRONLY, os.FileMode(header.Mode))
			if err != nil {
				return err
			}
			if _, err := io.Copy(outFile, tr); err != nil {
				_ = outFile.Close()
				return err
			}
			if err := outFile.Close(); err != nil {
				return err
			}
		case tar.TypeSymlink:
			_ = os.Remove(target)
			if err := os.Symlink(header.Linkname, target); err != nil {
				return err
			}
		}
	}

	return nil
}

func (m *Manager) extractZip(archivePath, targetDir string) error {
	r, err := zip.OpenReader(archivePath)
	if err != nil {
		return err
	}
	defer func() { _ = r.Close() }()

	for _, f := range r.File {
		target := filepath.Join(targetDir, f.Name)

		if f.FileInfo().IsDir() {
			_ = os.MkdirAll(target, 0755)
			continue
		}

		if err := os.MkdirAll(filepath.Dir(target), 0755); err != nil {
			return err
		}

		outFile, err := os.OpenFile(target, os.O_CREATE|os.O_WRONLY, f.Mode())
		if err != nil {
			return err
		}

		rc, err := f.Open()
		if err != nil {
			_ = outFile.Close()
			return err
		}

		_, err = io.Copy(outFile, rc)
		_ = rc.Close()
		if closeErr := outFile.Close(); closeErr != nil && err == nil {
			err = closeErr
		}
		if err != nil {
			return err
		}
	}

	return nil
}

// インストール後コマンドを実行する
func (m *Manager) runPostInstall(p *plugin.Plugin, installDir string) error {
	// TODO: インストール後コマンド実行を実装
	// os/exec を使ってインストールディレクトリでコマンドを実行
	fmt.Printf("   ⚠️  インストール後コマンドはまだ実装されていません\n")
	return nil
}

// リモートバージョン情報を表す
type RemoteVersion struct {
	Version string
	LTS     string // "" または LTS コードネーム（"Krypton" など）
}

// リモートから利用可能なバージョン一覧を取得する
func (m *Manager) ListRemote(toolName string, limit int) ([]RemoteVersion, error) {
	p, err := m.registry.Get(toolName)
	if err != nil {
		return nil, err
	}

	if p.ListURL == "" {
		return nil, fmt.Errorf("%s は ls-remote に対応していません", toolName)
	}

	// リモートから取得
	resp, err := http.Get(p.ListURL)
	if err != nil {
		return nil, fmt.Errorf("リモート取得エラー: %w", err)
	}
	defer func() { _ = resp.Body.Close() }()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("HTTP %d: %s", resp.StatusCode, resp.Status)
	}

	// JSON フォーマットをパース
	if p.ListFormat != "json" {
		return nil, fmt.Errorf("サポートされていないフォーマット: %s (現在は json のみ対応)", p.ListFormat)
	}

	var data []map[string]interface{}
	if err := json.NewDecoder(resp.Body).Decode(&data); err != nil {
		return nil, fmt.Errorf("JSON パースエラー: %w", err)
	}

	// バージョンを抽出
	versions := make([]RemoteVersion, 0, len(data))
	for _, item := range data {
		if ver, ok := item["version"].(string); ok {
			// version_prefix を削除
			if p.VersionPrefix != "" && strings.HasPrefix(ver, p.VersionPrefix) {
				ver = strings.TrimPrefix(ver, p.VersionPrefix)
			}

			// LTS 情報を取得
			lts := ""
			if ltsVal, ok := item["lts"]; ok {
				// lts は false または文字列（コードネーム）
				if ltsStr, ok := ltsVal.(string); ok {
					lts = ltsStr
				}
			}

			versions = append(versions, RemoteVersion{
				Version: ver,
				LTS:     lts,
			})
		}
	}

	// 件数制限
	if limit > 0 && len(versions) > limit {
		versions = versions[:limit]
	}

	return versions, nil
}
