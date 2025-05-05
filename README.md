# MedicalQuest

医療系の問題を出題・回答するCLIアプリケーション

## 技術スタック

- Go言語
- Echo (Webフレームワーク)
- Ent (ORM)
- PostgreSQL (データベース)
- Atlas (データベースマイグレーション)
- Docker (開発環境)

## 開発環境のセットアップ

### 前提条件

- Docker と Docker Compose がインストールされていること
- Go 1.21 以上がインストールされていること

### 環境構築手順

1. リポジトリをクローン

```bash
git clone https://github.com/yourusername/medicalquest.git
cd medicalquest
```

2. Docker コンテナを起動

```bash
make up
```

これにより、以下のコンテナが起動します：
- API サーバー: localhost:8080
- PostgreSQL: localhost:5432

3. マイグレーションを実行（初回のみ）

```bash
make migrate-dev
```

## 開発コマンド

```bash
# Docker コンテナの起動
make up

# Docker コンテナの停止
make down

# Docker イメージのビルド
make build

# テストの実行
make test

# リントチェック
make lint

# ローカルでの実行（コンテナ外）
make run

# コード生成（Ent スキーマ変更後）
make generate

# マイグレーション（開発環境）
make migrate-dev
```

## ディレクトリ構造

```
.
├── cmd/                 # アプリケーションのエントリーポイント
│   └── api/             # API サーバー
├── internal/            # 内部パッケージ
│   ├── config/          # 設定
│   ├── handlers/        # HTTPハンドラー
│   ├── models/          # データモデル
│   ├── repositories/    # データアクセス層
│   ├── services/        # ビジネスロジック
│   └── ent/             # Ent ORM 生成コード
├── pkg/                 # 公開パッケージ
│   ├── database/        # データベース接続
│   └── logger/          # ロギング
├── migrations/          # データベースマイグレーション
└── docker-compose.yml   # Docker 設定
```

## APIエンドポイント

API サーバーは以下のエンドポイントを提供します：

### 問題 (Questions)

- `GET /api/questions` - 全ての問題を取得
- `GET /api/questions/:id` - 指定したIDの問題を取得
- `POST /api/questions` - 新しい問題を作成
- `PUT /api/questions/:id` - 問題を更新
- `DELETE /api/questions/:id` - 問題を削除

### 選択肢 (Choices)

- `GET /api/questions/:questionID/choices` - 問題に紐づく選択肢を取得
- `POST /api/questions/:questionID/choices` - 選択肢を作成
- `PUT /api/choices/:id` - 選択肢を更新
- `DELETE /api/choices/:id` - 選択肢を削除