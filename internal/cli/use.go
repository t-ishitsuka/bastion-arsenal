package cli

import (
	"fmt"
	"os"
	"strings"

	"github.com/arsenal/internal/terminal"
	"github.com/spf13/cobra"
)

func newUseCmd() *cobra.Command {
	var local bool

	cmd := &cobra.Command{
		Use:   "use <tool> <version>",
		Short: "ツールのアクティブバージョンを切り替え",
		Long: `指定したツールのアクティブバージョンを切り替えます。

--local フラグを指定すると、現在のディレクトリの .toolversions ファイルに
バージョンを記録します。

使用例:
  arsenal use node 20.10.0
  arsenal use go 1.22.0 --local`,
		Args: cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			return runUse(args[0], args[1], local)
		},
	}

	cmd.Flags().BoolVarP(&local, "local", "l", false, ".toolversions に記録")

	return cmd
}

func runUse(toolName, version string, local bool) error {
	// プラグイン情報を取得（存在確認）
	p, err := registry.Get(toolName)
	if err != nil {
		return err
	}

	// バージョンを切り替え
	if err := manager.Use(toolName, version); err != nil {
		return err
	}

	// --local が指定された場合は .toolversions に書き込む
	if local {
		if err := updateToolVersionsFile(toolName, version); err != nil {
			return fmt.Errorf(".toolversions 更新エラー: %w", err)
		}
		terminal.PrintSuccess(".toolversions に %s %s を記録しました", p.DisplayName, version)
	}

	return nil
}

func updateToolVersionsFile(toolName, ver string) error {
	cwd, err := os.Getwd()
	if err != nil {
		return err
	}

	// 既存の .toolversions を読み込む（なければ新規作成）
	toolversionsPath := cwd + "/.toolversions"

	// ファイルが存在しない場合は新規作成
	var tools map[string]string
	if _, err := os.Stat(toolversionsPath); os.IsNotExist(err) {
		tools = make(map[string]string)
	} else {
		// 既存のファイルを読み込む
		tools, err = readToolVersionsSimple(toolversionsPath)
		if err != nil {
			return err
		}
	}

	// バージョンを更新
	tools[toolName] = ver

	// .toolversions に書き込む
	f, err := os.Create(toolversionsPath)
	if err != nil {
		return err
	}
	defer func() { _ = f.Close() }()

	for tool, version := range tools {
		if _, err := fmt.Fprintf(f, "%s %s\n", tool, version); err != nil {
			return err
		}
	}

	return nil
}

func readToolVersionsSimple(path string) (map[string]string, error) {
	content, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	tools := make(map[string]string)
	lines := string(content)
	for _, line := range strings.Split(lines, "\n") {
		line = strings.TrimSpace(line)
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}
		parts := strings.Fields(line)
		if len(parts) == 2 {
			tools[parts[0]] = parts[1]
		}
	}

	return tools, nil
}
