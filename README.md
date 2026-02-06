# Arsenal

[![Test](https://github.com/YOUR_USERNAME/arsenal/workflows/Test/badge.svg)](https://github.com/YOUR_USERNAME/arsenal/actions)

軽量マルチランタイムバージョンマネージャー。Bastion エコシステムの一部。

## 現在の状態

**開発中** - コアロジック実装済み、基本的な CLI コマンド実装済み。

### 実装済み
- コアバージョン管理ロジック (install/use/uninstall/list/current/sync/doctor)
- CLI コマンド: `install`, `use`, `ls-remote`, `plugin list`, `current`, `ls`, `doctor`
- go:embed を使ったプラグインシステム
- .toolversions パーサー
- Node.js プラグイン定義 (node.toml)
- パス管理とディレクトリ構造
- symlink ベースのバージョン切り替え
- GitHub Actions による自動テスト・lint・ビルド
- テストカバレッジ: 37%+ (CLI: 67.9%)

### TODO
- CLI コマンド実装 (uninstall.go, sync.go, initshell.go)
- 追加プラグイン定義 (go.toml, python.toml, rust.toml, php.toml)
- シェル統合 (init-shell コマンド)
- post-install コマンド実行
- ls-remote --lts-only フラグ

## インストール

```bash
go install github.com/arsenal/cmd/arsenal@latest
```

## シェル設定 (未実装)

```bash
# ~/.bashrc または ~/.zshrc に追加
eval "$(arsenal init-shell zsh)"
```

## 使用方法

### 実装済みコマンド

```bash
# リモートから利用可能なバージョンを確認
arsenal ls-remote node
arsenal ls-remote node --limit 50

# ツールのバージョンをインストール
arsenal install node 20.10.0

# バージョンを切り替え
arsenal use node 20.10.0
arsenal use go 1.22.0 --local   # .toolversions に書き込み

# 環境ヘルスチェック
arsenal doctor

# 利用可能なツール一覧
arsenal plugin list

# アクティブバージョンを表示
arsenal current

# インストール済みバージョン一覧
arsenal ls node
```

### 未実装（実装予定）

```bash
# バージョンをアンインストール
arsenal uninstall node 18.0.0

# .toolversions から同期
arsenal sync
```

## .toolversions フォーマット

```
# プロジェクトのツール要件
node 20.10.0
go 1.22.0
python 3.12.0
```

## Bastion 連携

Arsenal は `bastion init` 時に呼び出される `sync` と `doctor` を提供:

```yaml
# .bastion.yaml
environment:
  runtime:
    node: "20.10.0"
    go: "1.22.0"
```

```
bastion init → arsenal sync → 全ツール準備完了
```

## アーキテクチャ

```
~/.arsenal/
├── versions/        # インストール済みバージョン
│   ├── node/20.10.0/
│   └── go/1.22.0/
├── current/         # アクティブバージョンへの symlink
│   ├── node → ../versions/node/20.10.0
│   └── go → ../versions/go/1.22.0
└── plugins/         # カスタムツール定義 (TOML)
```

## 対応予定ツール

| ツール   | 状態 |
|---------|------|
| Node.js | プラグイン定義実装済み |
| Go      | プラグイン定義 TODO |
| Python  | プラグイン定義 TODO |
| Rust    | プラグイン定義 TODO |
| PHP     | プラグイン定義 TODO |

## ライセンス

MIT
