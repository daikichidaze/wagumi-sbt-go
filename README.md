# Wagumi Token working group (Jul-Sep 2022)

和組 DAO による Soul bound token 発行を行うためのメタデータ JSON 生成クライアント CLI です。NotionDB よりレピュテーションデータを取得し、JSON データへ整形します。

## Development Envoirment

Go 1.19.0

## Makefile

下記コンパイルに対応。Apple silicon 端末への対応は今後検討。

### Linux amd64

```
make linux-amd64
```

### Mac Intel

```
make mac-amd64
```

### Mac Apple silicon (未検証)

```
make mac-arm64
```

### Windows 64bit

```
make win
```

## Execution directory

実行のためには下記ファイル・ディレクトリの準備が必要です。

```
.env
metadata/
```

.env: Notion API を呼び出すためのトークン及び Database ID の指定

metadata/: 生成するメタデータファイルを出力するフォルダ
