package cli

import (
	"fmt"
	"os"
)

// ANSI エスケープコード
const (
	colorReset  = "\x1b[0m"
	colorRed    = "\x1b[31m"
	colorGreen  = "\x1b[32m"
	colorYellow = "\x1b[33m"
	colorBlue   = "\x1b[34m"
	colorCyan   = "\x1b[36m"
)

// 色付けが有効かどうかを判定する
func isColorEnabled() bool {
	// NO_COLOR 環境変数が設定されていたら無効
	if os.Getenv("NO_COLOR") != "" {
		return false
	}
	// ターミナルかどうかをチェック
	if fileInfo, _ := os.Stdout.Stat(); (fileInfo.Mode() & os.ModeCharDevice) == 0 {
		return false
	}
	return true
}

// 緑色のテキスト（成功メッセージ用）
func green(text string) string {
	if !isColorEnabled() {
		return text
	}
	return colorGreen + text + colorReset
}

// 赤色のテキスト（エラーメッセージ用）
func red(text string) string {
	if !isColorEnabled() {
		return text
	}
	return colorRed + text + colorReset
}

// 黄色のテキスト（警告メッセージ用）
func yellow(text string) string {
	if !isColorEnabled() {
		return text
	}
	return colorYellow + text + colorReset
}

// 青色のテキスト（情報メッセージ用）
func blue(text string) string {
	if !isColorEnabled() {
		return text
	}
	return colorBlue + text + colorReset
}

// シアン色のテキスト（進捗表示用）
func cyan(text string) string {
	if !isColorEnabled() {
		return text
	}
	return colorCyan + text + colorReset
}

// 成功メッセージを表示
func printSuccess(format string, args ...interface{}) {
	fmt.Printf(green("✅ "+format)+"\n", args...)
}

// エラーメッセージを表示
func printError(format string, args ...interface{}) {
	fmt.Printf(red("✗ "+format)+"\n", args...)
}

// 警告メッセージを表示
func printWarning(format string, args ...interface{}) {
	fmt.Printf(yellow("⚠️  "+format)+"\n", args...)
}

// 情報メッセージを表示
func printInfo(format string, args ...interface{}) {
	fmt.Printf(blue("📦 "+format)+"\n", args...)
}
