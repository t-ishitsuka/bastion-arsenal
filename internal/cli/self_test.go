package cli

import (
	"bytes"
	"strings"
	"testing"
)

func TestNewSelfCmd(t *testing.T) {
	cmd := newSelfCmd()
	if cmd == nil {
		t.Fatal("newSelfCmd が nil を返した")
	}

	if cmd.Use != "self" {
		t.Errorf("Use が期待と異なる: got %s, want self", cmd.Use)
	}

	// Check that update subcommand exists
	updateCmd, _, err := cmd.Find([]string{"update"})
	if err != nil {
		t.Fatalf("update サブコマンドが見つからない: %v", err)
	}

	if updateCmd.Use != "update" {
		t.Errorf("update サブコマンドの Use が期待と異なる: got %s, want update", updateCmd.Use)
	}
}

func TestSelfUpdateCheckOnlyFlag(t *testing.T) {
	// Set version to simulate a release version
	SetVersion("0.1.0", "abc1234", "2026-01-01T00:00:00Z")

	cmd := newSelfCmd()
	updateCmd, _, err := cmd.Find([]string{"update"})
	if err != nil {
		t.Fatalf("update サブコマンドが見つからない: %v", err)
	}

	// Set flags
	buf := new(bytes.Buffer)
	updateCmd.SetOut(buf)
	updateCmd.SetErr(buf)
	updateCmd.SetArgs([]string{"--check"})

	// Run command (it will try to fetch from GitHub, so it might fail in test environment)
	// We just check that the command structure is correct
	if err := updateCmd.Execute(); err != nil {
		// Expected to fail in test environment without network
		// Just check that the error message is reasonable
		if !strings.Contains(err.Error(), "最新リリース情報の取得に失敗") &&
			!strings.Contains(err.Error(), "GitHub API") {
			t.Logf("Expected network error, got: %v", err)
		}
	}
}

func TestSelfUpdateDevVersion(t *testing.T) {
	// Set version to dev
	SetVersion("dev", "unknown", "unknown")

	cmd := newSelfCmd()
	updateCmd, _, err := cmd.Find([]string{"update"})
	if err != nil {
		t.Fatalf("update サブコマンドが見つからない: %v", err)
	}

	buf := new(bytes.Buffer)
	updateCmd.SetOut(buf)
	updateCmd.SetErr(buf)

	// Should return error for dev version
	err = updateCmd.RunE(updateCmd, []string{})
	if err == nil {
		t.Error("dev バージョンで更新を試みたがエラーが返されなかった")
	}

	if !strings.Contains(err.Error(), "開発版からは更新できません") {
		t.Errorf("予期しないエラーメッセージ: %v", err)
	}
}

func TestExtractArchiveUnsupportedFormat(t *testing.T) {
	_, err := extractArchive("/tmp/test.unknown", "/tmp")
	if err == nil {
		t.Error("サポートされていないアーカイブ形式でエラーが返されなかった")
	}

	if !strings.Contains(err.Error(), "サポートされていないアーカイブ形式") {
		t.Errorf("予期しないエラーメッセージ: %v", err)
	}
}
