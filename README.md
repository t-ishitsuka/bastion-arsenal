# Arsenal

[![Test](https://github.com/t-ishitsuka/bastion-arsenal/workflows/Test/badge.svg)](https://github.com/t-ishitsuka/bastion-arsenal/actions)

軽量マルチランタイムバージョンマネージャー。asdf/mise に依存せず、シンプルで軽量な独自実装を目指す。Bastion エコシステムの一部として、開発環境の統一管理を実現。

## インストール

### プリビルドバイナリ（推奨）

Linux / macOS:
```bash
curl -fsSL https://raw.githubusercontent.com/t-ishitsuka/bastion-arsenal/main/scripts/install.sh | bash
```

Windows (PowerShell):
```powershell
iwr -useb https://raw.githubusercontent.com/t-ishitsuka/bastion-arsenal/main/scripts/install.ps1 | iex
```

### GitHub Releases から手動ダウンロード

[Releases ページ](https://github.com/t-ishitsuka/bastion-arsenal/releases)から最新版をダウンロードして展開し、バイナリを PATH が通ったディレクトリに配置してください。

### Go からビルド（開発者向け）

```bash
go install github.com/arsenal/cmd/arsenal@latest
```

## クイックスタート

```bash
# シェル設定
eval "$(bastion-arsenal init-shell bash)"  # ~/.bashrc に追加推奨

# Node.js をインストール
bastion-arsenal install node 20.10.0
bastion-arsenal use node 20.10.0

# 確認
node --version
```

## 主要機能

- **バージョン管理**: install/use/uninstall/ls コマンドで簡単管理
- **プロジェクト同期**: .toolversions から一括セットアップ（`bastion-arsenal sync`）
- **自動更新**: GitHub Releases から最新版に自動更新（`bastion-arsenal self update`）
- **シェル統合**: bash/zsh/fish 対応
- **リッチUI**: カラー出力、プログレスバー、LTSフィルタリング
- **プラグインシステム**: TOML で簡単にツールを追加可能

## 基本コマンド

| コマンド                                   | 説明                      |
| ------------------------------------------ | ------------------------- |
| `bastion-arsenal install <tool> <version>` | バージョンをインストール  |
| `bastion-arsenal use <tool> <version>`     | バージョン切り替え        |
| `bastion-arsenal ls-remote <tool>`         | リモートのバージョン一覧  |
| `bastion-arsenal sync`                     | .toolversions から同期    |
| `bastion-arsenal self update`              | Arsenal を最新版に更新    |
| `bastion-arsenal doctor`                   | 環境チェック              |
| `bastion-arsenal version`                  | バージョン情報を表示      |

詳細は `bastion-arsenal --help` を参照。

## アーキテクチャ

symlink 方式で高速にバージョンを切り替え（shims 不使用）。

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

PATH に `~/.arsenal/current/*/bin` を追加するだけで動作。

## .toolversions

プロジェクトルートに配置して `bastion-arsenal sync` で一括セットアップ。

```
# プロジェクトのツール要件
node 20.10.0
go 1.22.0
python 3.12.0
```

## Bastion 連携

Arsenal は Bastion 初期化時に自動実行され、開発環境を整備。

```
bastion init
  └─→ bastion-arsenal sync  # .toolversions から自動セットアップ
```

## 対応ツール

| ツール  | 状態     |
| ------- | -------- |
| Node.js | 対応済み |
| Go      | 準備中   |
| Python  | 準備中   |
| Rust    | 準備中   |

## 開発

```bash
# ビルド（バイナリ名: bastion-arsenal）
make build

# テスト
make test

# クロスコンパイル（全プラットフォーム）
make build-all

# リリースアーカイブ作成
make release

# カバレッジ: 73%+ (目標50%達成)
```

## ライセンス

MIT
