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

CLI アプリケーションは API サーバーと連携して以下の機能を提供します：

### 機能

- 問題一覧の表示
- 特定問題の詳細表示
- APIサーバーとの連携

### 使用方法

#### ビルド

```bash
# バイナリのビルド
go build -o medicalquest-cli cmd/cli/main.go
```

#### 基本コマンド

```bash
# ヘルプの表示
./medicalquest-cli --help

# API URLを指定（デフォルト: http://localhost:8080）
./medicalquest-cli --api-url http://localhost:8080 [command]
```

#### 利用可能なコマンド

##### 問題一覧の表示

```bash
# 全ての問題一覧を表示
./medicalquest-cli list

# API URLを指定して実行
./medicalquest-cli --api-url http://localhost:8080 list
```

実行例:
```
問題一覧:
ID: abcd-1234-efgh-5678
参照コード: REF001
タイトル: 心臓の解剖学
内容: 心臓の弁は全部で何枚ありますか？

ID: ijkl-9012-mnop-3456
参照コード: REF002
タイトル: 血液循環
内容: 大動脈から出た血液が心臓に戻るまでの経路は？
```

##### 特定問題の詳細表示

```bash
# 問題IDを指定して詳細を表示
./medicalquest-cli get <問題ID>

# 例
./medicalquest-cli get abcd-1234-efgh-5678
```

実行例:
```
問題詳細:
ID: abcd-1234-efgh-5678
参照コード: REF001
タイトル: 心臓の解剖学
内容: 心臓の弁は全部で何枚ありますか？
```

#### Docker での実行

```bash
# Docker内でビルド・実行
docker-compose run cli go build -o medicalquest-cli cmd/cli/main.go
docker-compose run cli ./medicalquest-cli list

# または直接実行
docker-compose run cli go run cmd/cli/main.go list
```

#### 注意事項

- CLIを使用する前にAPIサーバーが起動している必要があります
- デフォルトのAPI URLは `http://localhost:8080` です
- APIサーバーが異なるポートで動作している場合は `--api-url` フラグで指定してください

## CI/CD

### GitHub Actions

Pull Requestに対して以下の自動チェックが実行されます：

#### 実行されるジョブ

1. **Repository Tests**
  - `infrastructure/persistence` パッケージのテストを実行
  - カバレッジレポートの生成
  - テスト結果のアーティファクトアップロード

2. **Lint**
  - `go vet` によるコード検証
  - `golangci-lint` による静的解析

3. **Build**
  - APIとCLIのビルド確認

#### ローカルでのテスト実行

CIと同じテストをローカルで実行する場合：

```bash
# Entコードの生成
go generate ./infrastructure/ent

# Repositoryテストの実行
go test -v -race -coverprofile=coverage.out ./infrastructure/persistence/...

# カバレッジレポートの確認
go tool cover -html=coverage.out