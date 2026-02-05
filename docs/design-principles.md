# 設計方針・制約

## 軽量化

- 外部依存は最小限（cobra + toml のみ）
- シングルバイナリ
- shims を使わない（symlink 方式で高速化）

### なぜ shims を使わないか

- shims 方式: 毎回プロセス起動が必要で遅い
- symlink 方式: 直接バイナリを呼び出すため高速
- PATH に `~/.arsenal/current/*/bin` を追加するだけ

## 拡張性

- プラグインは TOML で宣言的に定義
- ユーザーが `~/.arsenal/plugins/` に TOML を置けば独自ツールを追加可能
- 組み込みプラグインは `go:embed` で同梱

## Bastion 連携インターフェース

Bastion から Arsenal を呼び出す際のインターフェース（将来実装）:

- `arsenal sync` - .toolversions or bastion.yaml のランタイム定義から一括セットアップ
- `arsenal doctor` - 環境チェック結果を返す
- `--output=json` フラグ（未実装）- Bastion がパースしやすい出力

## エラーハンドリング

- エラーは呼び出し元に返す
- コンテキスト情報を付加してラップ
- ユーザー向けエラーメッセージは日本語
