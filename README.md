# MedicalQuest

医療系の問題を出題・回答するCLIアプリケーション

## 技術スタック

- Go言語 (v1.23)
- Echo (Webフレームワーク)
- Ent (ORM)
- PostgreSQL (データベース)
- Atlas (データベースマイグレーション)
- Docker (開発環境)

## プロジェクト概要

MedicalQuestは医療知識に関する問題を出題し、回答するためのアプリケーションです。
RESTful APIとCLIの両方のインターフェースを提供し、問題の出題・回答・結果の保存などの機能を備えています。

## 開発環境のセットアップ

### 前提条件

- Docker と Docker Compose がインストールされていること
- Go 1.21 以上がインストールされていること

### 環境構築手順

1. リポジトリをクローン

```bash
git clone https://github.com/takechiyo-19940627/medicalquest.git
cd medicalquest
```

2. Docker コンテナを起動

```bash
make up
```

これにより、以下のコンテナが起動します：
- API サーバー: localhost:8080
- PostgreSQL: localhost:5432
- CLI: インタラクティブモードで実行

3. マイグレーションを実行（初回のみ）

```bash
make migrate-apply
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
go test -v ./...

# リントチェック
go vet ./...
golangci-lint run ./...

# Ent コード生成
./scripts/generate-ent.sh

# マイグレーション関連
make migrate-apply  # マイグレーション適用
make migrate-status # マイグレーション状態確認
make migrate-diff MIGRATION_NAME=name  # 差分マイグレーション作成
```

## ディレクトリ構造

```
.
├── cmd/                 # アプリケーションのエントリーポイント
│   ├── api/             # API サーバー
│   └── cli/             # CLIアプリケーション
├── config/              # 設定
├── domain/              # ドメイン層
│   ├── entity/          # エンティティ
│   └── repository/      # リポジトリインターフェース
├── handler/             # HTTPハンドラー
├── infrastructure/      # インフラストラクチャ層
│   ├── database/        # データベース接続
│   ├── ent/             # Ent ORM 生成コード
│   └── persistence/     # リポジトリ実装
├── migrations/          # データベースマイグレーション
├── pkg/                 # 公開パッケージ
└── service/             # サービス層
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

## CLIアプリケーション

CLI アプリケーションは以下の機能を提供します：

- 問題の出題
- 回答の受付
- 結果の表示
- 結果のCSVファイル出力

実行方法:

```bash
# Docker内で実行
docker-compose run cli

# ローカルで実行
go run cmd/cli/main.go
```