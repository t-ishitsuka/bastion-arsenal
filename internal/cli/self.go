package cli

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
	"runtime"
	"strings"

	"github.com/arsenal/internal/terminal"
	"github.com/spf13/cobra"
)

const (
	githubRepo    = "t-ishitsuka/bastion-arsenal"
	githubAPIBase = "https://api.github.com"
)

var (
	checkOnly   bool
	forceUpdate bool
)

type releaseInfo struct {
	TagName string `json:"tag_name"`
	Assets  []struct {
		Name               string `json:"name"`
		BrowserDownloadURL string `json:"browser_download_url"`
	} `json:"assets"`
}

func newSelfCmd() *cobra.Command {
	selfCmd := &cobra.Command{
		Use:   "self",
		Short: "Arsenal 自体の管理",
		Long:  "Arsenal バイナリの更新など、Arsenal 自体を管理するコマンドです。",
	}

	updateCmd := &cobra.Command{
		Use:   "update",
		Short: "Arsenal を最新版に更新",
		Long: `GitHub Releases から最新版をダウンロードして、Arsenal を更新します。

使用例:
  bastion-arsenal self update         最新版に更新
  bastion-arsenal self update --check 更新をチェックのみ
  bastion-arsenal self update --force バージョンが同じでも強制更新`,
		RunE: runSelfUpdate,
	}

	updateCmd.Flags().BoolVar(&checkOnly, "check", false, "更新をチェックするのみで、実際には更新しない")
	updateCmd.Flags().BoolVar(&forceUpdate, "force", false, "バージョンが同じでも強制的に更新する")

	selfCmd.AddCommand(updateCmd)
	return selfCmd
}

func runSelfUpdate(cmd *cobra.Command, _ []string) error {
	terminal.PrintInfo("更新をチェック中...")

	// Get current version
	currentVersion := versionInfo.Version
	if currentVersion == "dev" || currentVersion == "unknown" {
		return fmt.Errorf("開発版からは更新できません。GitHub Releases からインストールしてください")
	}

	// Get latest release info
	release, err := getLatestRelease()
	if err != nil {
		return fmt.Errorf("最新リリース情報の取得に失敗: %w", err)
	}

	latestVersion := strings.TrimPrefix(release.TagName, "v")
	currentVersionClean := strings.TrimPrefix(currentVersion, "v")

	terminal.PrintfBlue("現在のバージョン: %s\n", currentVersionClean)
	terminal.PrintfBlue("最新のバージョン: %s\n", latestVersion)

	// Check if update is needed
	if currentVersionClean == latestVersion && !forceUpdate {
		terminal.PrintSuccess("既に最新版です")
		return nil
	}

	if checkOnly {
		if currentVersionClean != latestVersion {
			terminal.PrintfYellow("新しいバージョンが利用可能です: %s\n", latestVersion)
			terminal.PrintfBlue("\n更新するには次のコマンドを実行してください:\n")
			terminal.PrintfBlue("  bastion-arsenal self update\n")
		}
		return nil
	}

	// Detect platform
	platform := fmt.Sprintf("%s-%s", runtime.GOOS, runtime.GOARCH)
	terminal.PrintfBlue("プラットフォーム: %s\n", platform)

	// Find the appropriate asset
	archiveName := fmt.Sprintf("bastion-arsenal-%s-%s", release.TagName, platform)
	if runtime.GOOS == "windows" {
		archiveName += ".zip"
	} else {
		archiveName += ".tar.gz"
	}

	var downloadURL string
	for _, asset := range release.Assets {
		if asset.Name == archiveName {
			downloadURL = asset.BrowserDownloadURL
			break
		}
	}

	if downloadURL == "" {
		return fmt.Errorf("このプラットフォーム用のバイナリが見つかりません: %s", platform)
	}

	// Download and install
	terminal.PrintInfo("最新版をダウンロード中...")
	if err := downloadAndInstall(downloadURL, archiveName); err != nil {
		return fmt.Errorf("更新に失敗: %w", err)
	}

	terminal.PrintSuccess("Arsenal を %s に更新しました", latestVersion)
	terminal.PrintfBlue("\n新しいバージョンを確認:\n")
	terminal.PrintfBlue("  bastion-arsenal version\n")

	return nil
}

func getLatestRelease() (*releaseInfo, error) {
	url := fmt.Sprintf("%s/repos/%s/releases/latest", githubAPIBase, githubRepo)

	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer func() { _ = resp.Body.Close() }()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("GitHub API エラー: %s", resp.Status)
	}

	var release releaseInfo
	if err := json.NewDecoder(resp.Body).Decode(&release); err != nil {
		return nil, err
	}

	return &release, nil
}

func downloadAndInstall(downloadURL, archiveName string) error {
	// Create temp directory
	tmpDir, err := os.MkdirTemp("", "arsenal-update-*")
	if err != nil {
		return fmt.Errorf("一時ディレクトリの作成に失敗: %w", err)
	}
	defer func() { _ = os.RemoveAll(tmpDir) }()

	// Download archive
	archivePath := filepath.Join(tmpDir, archiveName)
	if err := downloadFile(downloadURL, archivePath); err != nil {
		return fmt.Errorf("ダウンロードに失敗: %w", err)
	}

	// Extract archive
	terminal.PrintInfo("展開中...")
	binaryPath, err := extractArchive(archivePath, tmpDir)
	if err != nil {
		return fmt.Errorf("展開に失敗: %w", err)
	}

	// Get current executable path
	execPath, err := os.Executable()
	if err != nil {
		return fmt.Errorf("実行ファイルパスの取得に失敗: %w", err)
	}

	// Resolve symlinks
	execPath, err = filepath.EvalSymlinks(execPath)
	if err != nil {
		return fmt.Errorf("シンボリックリンクの解決に失敗: %w", err)
	}

	// Backup current binary
	backupPath := execPath + ".backup"
	terminal.PrintInfo("現在のバージョンをバックアップ中...")
	if err := os.Rename(execPath, backupPath); err != nil {
		return fmt.Errorf("バックアップに失敗: %w", err)
	}

	// Install new binary
	terminal.PrintInfo("新しいバージョンをインストール中...")
	if err := copyFile(binaryPath, execPath); err != nil {
		// Restore backup on failure
		_ = os.Rename(backupPath, execPath)
		return fmt.Errorf("インストールに失敗: %w", err)
	}

	// Set executable permission
	if err := os.Chmod(execPath, 0755); err != nil {
		return fmt.Errorf("実行権限の設定に失敗: %w", err)
	}

	// Remove backup on success
	_ = os.Remove(backupPath)

	return nil
}

func downloadFile(url, filepath string) error {
	resp, err := http.Get(url)
	if err != nil {
		return err
	}
	defer func() { _ = resp.Body.Close() }()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("ダウンロードエラー: %s", resp.Status)
	}

	out, err := os.Create(filepath)
	if err != nil {
		return err
	}
	defer func() { _ = out.Close() }()

	_, err = io.Copy(out, resp.Body)
	return err
}

func extractArchive(archivePath, destDir string) (string, error) {
	// Extract based on archive type
	if strings.HasSuffix(archivePath, ".tar.gz") {
		return extractTarGz(archivePath, destDir)
	} else if strings.HasSuffix(archivePath, ".zip") {
		return extractZip(archivePath, destDir)
	}

	return "", fmt.Errorf("サポートされていないアーカイブ形式: %s", archivePath)
}

func extractTarGz(archivePath, destDir string) (string, error) {
	// Use tar command to extract
	// The archive contains a binary with platform-specific name like "bastion-arsenal-linux-amd64"
	tmpExtractDir := filepath.Join(destDir, "extract")
	if err := os.MkdirAll(tmpExtractDir, 0755); err != nil {
		return "", err
	}

	// Extract using Go's tar/gzip packages
	file, err := os.Open(archivePath)
	if err != nil {
		return "", err
	}
	defer func() { _ = file.Close() }()

	gzr, err := gzip.NewReader(file)
	if err != nil {
		return "", err
	}
	defer func() { _ = gzr.Close() }()

	tr := tar.NewReader(gzr)
	for {
		header, err := tr.Next()
		if err == io.EOF {
			break
		}
		if err != nil {
			return "", err
		}

		// Only extract the binary file
		if header.Typeflag == tar.TypeReg {
			targetPath := filepath.Join(tmpExtractDir, filepath.Base(header.Name))
			outFile, err := os.Create(targetPath)
			if err != nil {
				return "", err
			}

			if _, err := io.Copy(outFile, tr); err != nil {
				_ = outFile.Close()
				return "", err
			}
			_ = outFile.Close()

			// Return the first regular file found
			return targetPath, nil
		}
	}

	return "", fmt.Errorf("アーカイブ内にバイナリが見つかりません")
}

func extractZip(archivePath, destDir string) (string, error) {
	r, err := zip.OpenReader(archivePath)
	if err != nil {
		return "", err
	}
	defer func() { _ = r.Close() }()

	for _, f := range r.File {
		if f.FileInfo().IsDir() {
			continue
		}

		rc, err := f.Open()
		if err != nil {
			return "", err
		}

		targetPath := filepath.Join(destDir, filepath.Base(f.Name))
		outFile, err := os.Create(targetPath)
		if err != nil {
			_ = rc.Close()
			return "", err
		}

		_, err = io.Copy(outFile, rc)
		_ = rc.Close()
		_ = outFile.Close()

		if err != nil {
			return "", err
		}

		// Return the first file found
		return targetPath, nil
	}

	return "", fmt.Errorf("アーカイブ内にバイナリが見つかりません")
}

func copyFile(src, dst string) error {
	sourceFile, err := os.Open(src)
	if err != nil {
		return err
	}
	defer func() { _ = sourceFile.Close() }()

	destFile, err := os.Create(dst)
	if err != nil {
		return err
	}
	defer func() { _ = destFile.Close() }()

	if _, err := io.Copy(destFile, sourceFile); err != nil {
		return err
	}

	return destFile.Sync()
}
