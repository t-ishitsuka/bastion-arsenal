# Arsenal

軽量マルチランタイムバージョンマネージャー。Bastion エコシステムの一部。

## 現在の状態

**開発中** - コアロジック実装済み、CLI コマンド実装中。

### 実装済み
- コアバージョン管理ロジック (install/use/uninstall/list/current/sync/doctor)
- go:embed を使ったプラグインシステム
- .toolversions パーサー
- Node.js プラグイン定義 (node.toml)
- パス管理とディレクトリ構造
- symlink ベースのバージョン切り替え

### TODO
- CLI コマンド実装 (install.go, use.go など)
- 追加プラグイン定義 (go.toml, python.toml, rust.toml, php.toml)
- シェル統合 (init-shell コマンド)
- リモートバージョン一覧 (ls-remote コマンド)
- post-install コマンド実行

## インストール

```bash
go install github.com/arsenal/cmd/arsenal@latest
```

## シェル設定 (未実装)

```bash
# ~/.bashrc または ~/.zshrc に追加
eval "$(arsenal init-shell zsh)"
```

## 使用方法（予定）

```bash
# ツールのバージョンをインストール
arsenal install node 20.10.0
arsenal install go 1.22.0

# バージョンを切り替え
arsenal use node 20.10.0
arsenal use go 1.22.0 --local   # .toolversions に書き込み

# インストール済みバージョン一覧
arsenal ls node

# アクティブバージョンを表示
arsenal current

# .toolversions から同期
arsenal sync

# 環境ヘルスチェック
arsenal doctor

# 利用可能なツール一覧
arsenal plugin list
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
