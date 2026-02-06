# CLAUDE.md - Arsenal Project

## 言語設定
**日本語で応答してください。**

---

## 重要なルール

### 1. ドキュメントファイルでの絵文字禁止

**ドキュメントファイル（.md）には絵文字を一切使用しない**
- CLAUDE.md, README.md, docs/ 配下の全ファイルで絵文字禁止
- ターミナル出力（CLI実行時の表示）のみ絵文字使用可
- 理由: ドキュメントでの絵文字は好まれない
- 違反例: ✅, ❌, 🔧, 📦 など全ての絵文字
- 正しい表記: [実装済み], (未実装), ※注意 など

### 2. 実装とドキュメントの同期

**コード変更時は必ず関連ドキュメントを同時更新する**
- 新機能実装 → README.md, CLAUDE.md の「実装済み」「CLI コマンド一覧」を更新
- コマンド追加 → 使用例、コマンド一覧表を更新
- 仕様変更 → 該当する docs/ 配下のファイルを更新
- 実装状況が変わったら即座にドキュメント反映（後回しにしない）

対象ドキュメント: CLAUDE.md, README.md, docs/ 配下の全ファイル

---

## プロジェクト概要

**Arsenal（アーセナル）** は軽量マルチランタイムバージョンマネージャー。
asdf/mise 等の既存ツールに依存せず、自前で管理する学習目的＋軽量化がモチベーション。

### 上位プロジェクト：Bastion エコシステム

Arsenal は「要塞シリーズ」の一部。全体構成：

```
┌─────────────────────────────────────────────────────────┐
│                      BASTION（司令塔）                    │
│  Claude Code マルチエージェント管理システム                │
│  Envoy（指揮官）→ Marshall（監督）→ Specialist（実行者）  │
│  tmux + git worktree で並列管理                          │
│  実装言語: Go                                            │
├─────────────────────────────────────────────────────────┤
│         │                           │                    │
│         ▼                           ▼                    │
│  ┌─────────────┐             ┌─────────────┐            │
│  │   ARSENAL   │             │   CITADEL   │            │
│  │  （武器庫）  │             │  （城塞）   │            │
│  │ ランタイム   │             │  Docker     │            │
│  │ バージョン管理│             │  環境管理   │            │
│  └─────────────┘             └─────────────┘            │
└─────────────────────────────────────────────────────────┘
```

**連携フロー:**
```
bastion init
  ├─→ arsenal sync    (.toolversions からランタイム整備)
  ├─→ citadel up      (Docker サービス起動)
  └─→ worktree + tmux (エージェント並列起動)
```

**将来追加予定のツール:**
- Vault（金庫）: シークレット管理
- Forge（鍛冶場）: タスクランナー
- Sentinel（歩哨）: ヘルスチェック/待機
- Watchtower（監視塔）: ログ集約
- Rampart（城壁）: セキュリティ/証明書
- Courier（伝令）: 通知/外部連携
- Drawbridge（跳ね橋）: ネットワーク/トンネル

---

## 詳細ドキュメント

詳細な設計・実装情報は `docs/` ディレクトリを参照：

- **[architecture.md](docs/architecture.md)** - ディレクトリ構成、データ構造、パッケージ依存関係
- **[design-principles.md](docs/design-principles.md)** - 設計方針、制約、軽量化・拡張性の方針
- **[coding-standards.md](docs/coding-standards.md)** - コーディング規約、命名規則
- **[plugin-system.md](docs/plugin-system.md)** - プラグインシステムの仕様、TOML 定義
- **[toolversions.md](docs/toolversions.md)** - .toolversions ファイルフォーマット
- **[development.md](docs/development.md)** - ビルド、テスト、開発手順

---

## 現在の実装状況

### 実装済み

#### コアロジック (internal/version/)
- `manager.go` - Install, Use, Uninstall, List, Current, CurrentAll, Doctor メソッド実装済み
- `toolversions.go` - .toolversions パーサーと Sync 機能実装済み

#### プラグインシステム (internal/plugin/)
- `plugin.go` - プラグインレジストリ、go:embed によるビルトインプラグイン読み込み実装済み
- `builtin/node.toml` - Node.js プラグイン定義のみ実装済み

#### 設定管理 (internal/config/)
- `config.go` - パス管理、ディレクトリ構造実装済み

#### CLI フレームワーク (internal/cli/)
- `root.go` - ルートコマンドと初期化ロジック実装済み

#### エントリポイント
- `cmd/arsenal/main.go` - メイン関数実装済み

#### CLI コマンド (internal/cli/)
実装済み:
- `install.go` - arsenal install コマンド [実装済み]
- `use.go` - arsenal use コマンド [実装済み]
- `uninstall.go` - arsenal uninstall コマンド [実装済み]
- `lsremote.go` - arsenal ls-remote コマンド [実装済み]
- `plugin.go` - arsenal plugin list コマンド [実装済み]
- `current.go` - arsenal current コマンド [実装済み]
- `list.go` - arsenal ls コマンド [実装済み]
- `sync.go` - arsenal sync コマンド [実装済み]
- `doctor.go` - arsenal doctor コマンド [実装済み]
- `initshell.go` - arsenal init-shell コマンド [実装済み]

### 未実装

#### CLI コマンド
全ての基本コマンドが実装済み

#### プラグイン定義 (internal/plugin/builtin/)
- `go.toml` - Go プラグイン定義
- `python.toml` - Python プラグイン定義
- `rust.toml` - Rust プラグイン定義
- `php.toml` - PHP プラグイン定義

#### その他機能
- post_install 実行機能 - Python/Rust/PHP のビルド処理
- プログレスバー付きダウンロード

---

## 技術仕様概要

### 実装言語・依存

- **Go 1.22+**
- `github.com/spf13/cobra` - CLI フレームワーク
- `github.com/BurntSushi/toml` - プラグイン定義パーサー
- それ以外の外部依存は極力避ける（軽量化方針）

詳細は [architecture.md](docs/architecture.md) と [design-principles.md](docs/design-principles.md) を参照。

### テスト・CI/CD

- **テストカバレッジ**: 全体 41%+ (CLI: 73.3%, config: 84.6%, plugin: 66.1%)
- **GitHub Actions**: PR/push 時に自動テスト・lint・ビルド実行
- **golangci-lint**: errcheck, staticcheck, unused など標準リンター有効化
- **カバレッジ目標**: 最低 25%、目標 50% [達成]

### CLI コマンド一覧

| コマンド | 説明 | 状態 |
|---------|------|------|
| `arsenal install <tool> <version>` | バージョンをインストール | [実装済み] |
| `arsenal use <tool> <version>` | バージョン切り替え (symlink) | [実装済み] |
| `arsenal use <tool> <version> --local` | 切り替え + .toolversions に書き込み | [実装済み] |
| `arsenal uninstall <tool> <version>` | バージョン削除 | [実装済み] |
| `arsenal ls <tool>` | インストール済みバージョン一覧 | [実装済み] |
| `arsenal ls-remote <tool>` | リモートの利用可能バージョン取得 | [実装済み] |
| `arsenal current` | 全ツールのアクティブバージョン表示 | [実装済み] |
| `arsenal sync` | .toolversions から一括セットアップ | [実装済み] |
| `arsenal doctor` | 環境ヘルスチェック | [実装済み] |
| `arsenal plugin list` | 対応ツール一覧 | [実装済み] |
| `arsenal init-shell [bash\|zsh\|fish]` | シェル設定スクリプト出力 | [実装済み] |

### 対応ツールと状態

| ツール | インストール方式 | 状態 |
|--------|----------------|------|
| Node.js | プリビルドバイナリ | [プラグイン定義のみ実装] |
| Go | プリビルドバイナリ | プラグイン定義未実装 |
| Python | ソースからビルド | プラグイン定義未実装 |
| Rust | スタンドアロンインストーラ | プラグイン定義未実装 |
| PHP | ソースからビルド | プラグイン定義未実装 |

---

## 未実装・TODO

### 優先度高

1. **`post_install` コマンド実行** - Python/Rust/PHP のビルド
   - `os/exec` でシェルコマンド実行
   - `{{install_dir}}` テンプレート変数の置換
   - 作業ディレクトリを展開先に設定

4. **プログレスバー付きダウンロード**
   - `Content-Length` から総サイズ取得
   - `io.TeeReader` でプログレス表示

### 優先度中

5. **追加プラグイン定義** - go.toml, python.toml, rust.toml, php.toml
6. **`--output=json` フラグ** - Bastion 連携用
7. **tar.xz 展開サポート** - Python ソース配布用
8. **エラーハンドリング強化** - ネットワークエラーのリトライ等
9. **バージョンのエイリアス** - `arsenal use node lts` 等

### 優先度低

10. **自動バージョン切り替え** - `cd` 時に .toolversions を検知して自動 sync
11. **アップデートチェック** - 新しいバージョンの通知
12. **キャッシュ** - ダウンロード済みアーカイブの再利用
