# 開発ガイド

## ビルド・インストール

```bash
# ビルド
make build

# インストール（GOPATH/bin へ）
make install

# テスト
make test

# 開発用（~/.local/bin へコピー）
make dev

# クリーンアップ
make clean

# Lint
make lint
```

## 依存関係

```bash
# 依存関係の整理
go mod tidy

# 依存関係の更新
go get -u ./...
```

## テスト

```bash
# 全テスト実行
go test ./...

# カバレッジ付き
go test -cover ./...

# 詳細出力
go test -v ./...
```

## デバッグ

```bash
# バイナリをビルドしてデバッグ
go build -o arsenal ./cmd/arsenal
./arsenal --help

# デバッグ情報付きビルド
go build -gcflags="all=-N -l" -o arsenal ./cmd/arsenal
```

## 新しいコマンドの追加

1. `internal/cli/` に新しいファイルを作成（例: `newcommand.go`）
2. `newXxxCmd()` 関数を実装
3. `root.go` の `NewRootCmd()` で `root.AddCommand(newXxxCmd())` を追加
4. 必要に応じて `internal/version/manager.go` にロジックを追加

## 新しいプラグインの追加

1. `internal/plugin/builtin/` に TOML ファイルを作成（例: `ruby.toml`）
2. プラグイン定義を記述
3. ビルドすると自動的に `go:embed` で組み込まれる

## CI/CD

### GitHub Actions

プルリクエスト作成時、自動的に以下が実行されます：

- **Test**: 全テストを実行（カバレッジ 50% 以上必須）
- **Lint**: golangci-lint でコード品質チェック
- **Build**: ビルドが成功するか確認

### ローカルで CI と同じチェックを実行

```bash
# テスト（カバレッジ付き）
go test -v -race -coverprofile=coverage.out ./...
go tool cover -func=coverage.out

# カバレッジを HTML で表示
go tool cover -html=coverage.out -o coverage.html

# Lint
golangci-lint run ./...

# ビルド
make build
```

### カバレッジ要件

- 総カバレッジ: 最低 25%（目標: 50% 以上）
- 新規コードはテストを含めること（コーディング規約参照）
- CLI コマンド実装時はテストも同時に作成
