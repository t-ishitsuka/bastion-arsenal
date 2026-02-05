# プラグインシステム

## 概要

各ツールは TOML ファイルで定義。`go:embed` で組み込み、ユーザー定義で上書き可能。

## プラグイン定義フォーマット

```toml
# 例: node.toml
name = "node"
display_name = "Node.js"
description = "JavaScript runtime"
list_url = "https://nodejs.org/dist/index.json"
list_format = "json"
download_url = "https://nodejs.org/dist/v{{version}}/node-v{{version}}-{{os}}-{{arch}}.tar.gz"
bin_path = "bin"
archive_type = "tar.gz"
version_prefix = "v"

[os_map]
darwin = "darwin"
linux = "linux"
windows = "win"

[arch_map]
amd64 = "x64"
arm64 = "arm64"
```

## テンプレート変数

- `{{version}}` - バージョン番号
- `{{os}}` - OS 名（マッピング後）
- `{{arch}}` - アーキテクチャ（マッピング後）

OS/Arch は `runtime.GOOS` / `runtime.GOARCH` から取得し、マッピングで変換。

## フィールド説明

### 基本情報

- `name`: ツール名（コマンド引数で使用）
- `display_name`: 表示名
- `description`: 説明

### ダウンロード

- `list_url`: バージョン一覧取得 URL（ls-remote 用）
- `list_format`: 一覧のフォーマット（"json", "html", "github"）
- `download_url`: ダウンロード URL テンプレート

### インストール

- `bin_path`: アーカイブ内のバイナリパス
- `archive_type`: アーカイブ形式（"tar.gz", "tar.xz", "zip"）
- `version_prefix`: バージョン番号のプレフィックス（削除用）
- `version_regex`: バージョン抽出用正規表現

### マッピング

- `os_map`: OS 名のマッピング
- `arch_map`: アーキテクチャのマッピング

### 実行

- `post_install`: インストール後に実行するコマンド
- `env_vars`: 設定する環境変数

## プラグインの読み込み順序

1. 組み込みプラグイン（`internal/plugin/builtin/*.toml`）を `go:embed` で読み込み
2. ユーザープラグイン（`~/.arsenal/plugins/*.toml`）を読み込み（上書き）

## カスタムプラグインの追加

ユーザーは `~/.arsenal/plugins/` に TOML ファイルを配置することで、
独自ツールを追加または既存ツールの定義を上書きできる。
