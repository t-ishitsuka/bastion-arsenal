# CLAUDE.md - Arsenal Project

## 言語設定

**日本語で応答してください。** コード内コメントは英語OK。

---

## 重要なルール（最優先）

### 1. ドキュメントファイルでの絵文字禁止

**ドキュメントファイル（.md）には絵文字を一切使用しない**

- 対象: CLAUDE.md, README.md, docs/ 配下の全ファイル
- 理由: ユーザーがドキュメントでの絵文字を嫌う
- 違反例: ✅, ❌, 🔧, 📦, ⚠️ など全ての絵文字・Unicode記号
- 正しい表記: [実装済み], (未実装), ※注意
- 例外: ターミナル出力（CLI実行時の表示）のみ絵文字使用可

このルールは他のすべての規約より優先度が高い。

### 2. 実装とドキュメントの同期

**コード変更時は必ず関連ドキュメントを同時更新する**

- 新機能実装 → README.md, CLAUDE.md を即座に更新
- コマンド追加 → コマンド一覧表、使用例を同時に更新
- 仕様変更 → docs/ 配下のファイルを同時に更新
- 実装とドキュメントの不一致は厳禁

---

## プロジェクト概要

Arsenal は軽量マルチランタイムバージョンマネージャー。asdf/mise 等の既存ツールに依存せず、自前で管理する学習目的＋軽量化がモチベーション。

Bastion エコシステムの一部として、Claude Code マルチエージェント管理システムから呼び出される。`bastion init` 時に `arsenal sync` が自動実行され、.toolversions からランタイムを整備。

`~/.arsenal/` ディレクトリで管理し、symlink 方式でバージョンを切り替え（shims 不使用）。

---

## 実装状況

### 基本機能（全て実装済み）

| カテゴリ                               | 状態       |
| -------------------------------------- | ---------- |
| CLI コマンド（10種類）                 | [実装済み] |
| バージョン管理ロジック                 | [実装済み] |
| .toolversions パーサー                 | [実装済み] |
| シェル統合（bash/zsh/fish）            | [実装済み] |
| プラグインシステム                     | [実装済み] |
| ターミナルUI（カラー・プログレスバー） | [実装済み] |
| テスト（カバレッジ73%+）               | [実装済み] |

### 対応ツール

| ツール  | プラグイン定義 | 状態       |
| ------- | -------------- | ---------- |
| Node.js | node.toml      | [実装済み] |
| Go      | go.toml        | 未実装     |
| Python  | python.toml    | 未実装     |
| Rust    | rust.toml      | 未実装     |
| PHP     | php.toml       | 未実装     |

### TODO（優先度順）

1. **追加プラグイン定義** - go.toml, python.toml, rust.toml, php.toml
2. **post_install 実行** - Python/Rust/PHP のビルド処理
3. **エラーハンドリング強化** - ネットワークエラーのリトライ等

---

## 技術スタック

### 依存関係

- Go 1.22+
- github.com/spf13/cobra（CLI）
- github.com/BurntSushi/toml（プラグイン）
- それ以外の外部依存は極力避ける（軽量化方針）

### テスト・品質

- カバレッジ: 73%+ (CLI: 73.6%, config: 84.6%, terminal: 79.4%, plugin: 66.1%)
- GitHub Actions: PR/push 時に自動テスト・lint・ビルド
- golangci-lint: errcheck, staticcheck, unused など有効化

### ディレクトリ構造

```
arsenal/
├── cmd/arsenal/          # エントリポイント
├── internal/
│   ├── cli/              # CLI コマンド定義
│   ├── config/           # パス管理
│   ├── plugin/           # プラグインシステム
│   ├── terminal/         # ターミナルUI（カラー等）
│   └── version/          # コアバージョン管理ロジック
└── go.mod
```

---

## CLI コマンド一覧

| コマンド                         | 説明                                                   | 状態       |
| -------------------------------- | ------------------------------------------------------ | ---------- |
| `install <tool> <version>`       | バージョンをインストール                               | [実装済み] |
| `use <tool> <version> [--local]` | バージョン切り替え（--local で .toolversions に記録）  | [実装済み] |
| `uninstall <tool> <version>`     | バージョン削除（現在使用中なら自動で最新版に切り替え） | [実装済み] |
| `ls <tool>`                      | インストール済みバージョン一覧                         | [実装済み] |
| `ls-remote <tool> [--lts-only]`  | リモートのバージョン一覧（LTSフィルタ可）              | [実装済み] |
| `current`                        | 全ツールのアクティブバージョン表示                     | [実装済み] |
| `sync`                           | .toolversions から一括セットアップ                     | [実装済み] |
| `doctor`                         | 環境ヘルスチェック                                     | [実装済み] |
| `plugin list`                    | 対応ツール一覧                                         | [実装済み] |
| `init-shell [bash\|zsh\|fish]`   | シェル設定スクリプト出力                               | [実装済み] |

---

## コーディングパターン

### テスト実装

- 新しいコードを書く際は必ずテストを書く
- コードを変更したら必ず linter を実行する
- テーブル駆動テストを推奨

### エラーハンドリング

- Close(), Remove() などのクリーンアップ処理: `defer func() { _ = f.Close() }()`
- nil チェック後は early return でガード（staticcheck SA5011 対策）

### ドキュメントコメント

- 関数名を繰り返さない
- 良い例: `// 正しくパスを返すかテストする`
- 悪い例: `// TestGetPaths は GetPaths 関数が正しくパスを返すかテストする`

### ターミナルUI

- カラー: ANSI エスケープコード使用（internal/terminal/color.go）
- NO_COLOR 環境変数を尊重
- プログレスバー: \r で同じ行を上書き、io.TeeReader でストリーム追跡
- ダウンロード進捗: シアン色、完了: 緑色

---

## 詳細ドキュメント

以下のファイルで詳細情報を確認してください：

- **[architecture.md](docs/architecture.md)** - ディレクトリ構成、データ構造、パッケージ依存関係
- **[design-principles.md](docs/design-principles.md)** - 設計方針、制約、軽量化・拡張性の方針
- **[coding-standards.md](docs/coding-standards.md)** - コーディング規約、命名規則
- **[plugin-system.md](docs/plugin-system.md)** - プラグインシステムの仕様、TOML 定義
- **[toolversions.md](docs/toolversions.md)** - .toolversions ファイルフォーマット
- **[development.md](docs/development.md)** - ビルド、テスト、開発手順

---

## Bastion エコシステム連携

Arsenal は「要塞シリーズ」の一部：

```
BASTION（司令塔）
  ├─→ ARSENAL（武器庫）  # ランタイムバージョン管理
  └─→ CITADEL（城塞）    # Docker環境管理
```

将来追加予定：Vault（シークレット）、Forge（タスクランナー）、Sentinel（ヘルスチェック）など

連携フロー:

```
bastion init
  └─→ arsenal sync  # .toolversions からランタイム整備
```
