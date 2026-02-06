package cli

import (
	"bytes"
	"strings"
	"testing"
)

func TestVersionCmd(t *testing.T) {
	// Set test version info
	SetVersion("1.0.0", "abc1234", "2026-01-01T00:00:00Z")

	cmd := newVersionCmd()
	buf := new(bytes.Buffer)
	cmd.SetOut(buf)
	cmd.SetErr(buf)

	err := cmd.Execute()
	if err != nil {
		t.Fatalf("version コマンド実行失敗: %v", err)
	}

	output := buf.String()

	// Check output contains version information
	if !strings.Contains(output, "Arsenal 1.0.0") {
		t.Errorf("バージョン番号が出力されていない: %s", output)
	}
	if !strings.Contains(output, "Commit: abc1234") {
		t.Errorf("コミットハッシュが出力されていない: %s", output)
	}
	if !strings.Contains(output, "Built: 2026-01-01T00:00:00Z") {
		t.Errorf("ビルド日時が出力されていない: %s", output)
	}
}
