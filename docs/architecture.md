# アーキテクチャ

## ディレクトリ構成

```
arsenal/
├── cmd/arsenal/main.go              # エントリポイント
├── internal/
│   ├── cli/                         # Cobra コマンド定義
│   │   ├── root.go                  # ルートコマンド + 初期化
│   │   ├── install.go               # arsenal install <tool> <version>
│   │   ├── use.go                   # arsenal use <tool> <version> [--local]
│   │   ├── uninstall.go             # arsenal uninstall <tool> <version>
│   │   ├── list.go                  # arsenal ls <tool>
│   │   ├── current.go               # arsenal current
│   │   ├── sync.go                  # arsenal sync (.toolversions 一括適用)
│   │   ├── doctor.go                # arsenal doctor (環境ヘルスチェック)
│   │   ├── plugin.go                # arsenal plugin list
│   │   └── initshell.go             # arsenal init-shell [bash|zsh|fish]
│   ├── config/
│   │   └── config.go                # パス管理、グローバル設定
│   ├── plugin/
│   │   ├── plugin.go                # プラグインシステム (go:embed + TOML)
│   │   └── builtin/                 # 組み込みプラグイン定義
│   │       ├── node.toml
│   │       ├── go.toml
│   │       ├── python.toml
│   │       ├── rust.toml
│   │       └── php.toml
│   └── version/
│       ├── manager.go               # コアロジック (DL/展開/symlink/doctor)
│       └── toolversions.go          # .toolversions パーサー + sync
├── docs/                            # 設計文書
├── go.mod
├── Makefile
├── README.md
├── CLAUDE.md
└── .gitignore
```

## ランタイムデータ構造

```
~/.arsenal/
├── versions/              # インストール済みバージョン
│   ├── node/
│   │   ├── 20.10.0/       # 展開されたバイナリ一式
│   │   └── 18.19.0/
│   ├── go/
│   │   └── 1.22.0/
│   └── python/
│       └── 3.12.0/
├── current/               # アクティブバージョンへの symlink
│   ├── node → ../versions/node/20.10.0
│   └── go → ../versions/go/1.22.0
├── plugins/               # ユーザー定義プラグイン（TOML）
└── config.toml            # グローバル設定
```

## バージョン切り替え方式

**symlink 方式**（shims ではない）：

- `~/.arsenal/current/<tool>` → `~/.arsenal/versions/<tool>/<version>` への symlink
- PATH に `~/.arsenal/current/*/bin` を追加
- shims 方式より高速（毎回プロセス起動しない）

## パッケージ依存関係

```
cli → version → plugin → config
     ↓
   plugin
```

- `config`: 基本的なパス管理、設定（他に依存しない）
- `plugin`: プラグイン定義の読み込みと管理（config に依存）
- `version`: バージョン管理ロジック（config, plugin に依存）
- `cli`: コマンドライン UI（全てに依存）
