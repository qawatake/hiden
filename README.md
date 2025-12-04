# hiden

ghqで管理しているリポジトリ内の個人用メモ・スクリプト置き場（hidenディレクトリ）を横断的に検索・アクセスするためのCLIツール。

## インストール

```bash
go install github.com/qawatake/hiden@latest
```

## 使い方

### ファイル検索

```bash
# インクリメンタル検索でファイルを選択
hiden ls

# エディタで開く
vim $(hiden ls)
```

### 日付ディレクトリ作成

```bash
# 現在のリポジトリに今日の日付のディレクトリを作成
hiden mkdir
# => .hiden/2025-12-04

# 作成したディレクトリに移動
cd $(hiden mkdir)
```

## 設定

設定ファイル: `~/.config/hiden/config.json`

```json
{
  "dirname": ".hiden"
}
```

| フィールド | デフォルト | 説明 |
|-----------|-----------|------|
| `dirname` | `.hiden` | hidenディレクトリの名前 |

## ディレクトリ構成例

```
~/src/github.com/
├── org1/repo1/
│   ├── .gitignore   # .hiden/ を除外
│   └── .hiden/
│       ├── memo.md
│       └── scripts/test.sh
└── org2/repo2/
    ├── .gitignore   # .hiden/ を除外
    └── .hiden/
        └── notes.txt
```

## ライセンス

MIT
