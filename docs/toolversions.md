# .toolversions フォーマット

## 概要

`.toolversions` ファイルはプロジェクトで必要なツールのバージョンを定義する。
asdf 互換のフォーマット。

## フォーマット

```
# コメント
node 20.10.0
go 1.22.0
python 3.12.0
```

## ルール

- 1行1ツール、スペース区切り
- `#` でコメント
- 空行は無視
- ディレクトリを遡って検索（プロジェクトルートまで）

## ファイル検索

`arsenal sync` 実行時、以下の順序でファイルを検索:

1. カレントディレクトリの `.toolversions`
2. 親ディレクトリの `.toolversions`
3. ルート（`/`）まで遡る

最初に見つかったファイルを使用。

## 使用例

### プロジェクトルート

```
my-project/
├── .toolversions    # このファイルが使われる
├── src/
└── tests/
```

### サブディレクトリ

```
my-project/
├── .toolversions    # backend/ 以下でもこのファイルが使われる
├── backend/
│   └── src/
└── frontend/
    ├── .toolversions # frontend/ 以下ではこちらが優先
    └── src/
```

## arsenal sync の動作

1. `.toolversions` を検索・読み込み
2. 各ツールについて:
   - インストールされていなければインストール
   - バージョンを切り替え（symlink 更新）
3. エラーがあっても他のツールは続行
