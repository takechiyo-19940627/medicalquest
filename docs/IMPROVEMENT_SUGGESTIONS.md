# MedicalQuest 改善提案レポート

生成日: 2025-11-13

## エグゼクティブサマリー

MedicalQuestは医療系問題出題アプリケーションとして、良好なアーキテクチャ設計（オニオンアーキテクチャ）を採用していますが、以下の重要な改善領域が特定されました：

### 緊急度の高い問題
1. **Entコード生成エラー** - ビルドが失敗する可能性
2. **アーキテクチャ違反** - ChoiceHandlerがService層を経由せず、直接Entクライアントに依存
3. **未実装機能** - 複数のAPIエンドポイントがスタブのみ
4. **テストカバレッジ不足** - 全体の約25-30%（Handler層は0%）
5. **セキュリティ対策不足** - 認証・認可の欠如、CORS設定が寛容すぎる

---

## 1. ビルド・コード生成の問題

### 1.1 Entコード生成エラー

**問題**: `infrastructure/ent/hook/hook.go`でビルドエラーが発生
```
infrastructure/ent/hook/hook.go:14:44: undefined: ent.ChoiceMutation
infrastructure/ent/hook/hook.go:14:65: undefined: ent.Value
```

**影響度**: 🔴 **緊急**
ビルドが失敗し、アプリケーションが起動できない可能性があります。

**対応方法**:
```bash
# Entコードの再生成
./scripts/generate-ent.sh
# または
go generate ./infrastructure/ent
```

**推奨対応**: CI/CDパイプラインに以下を追加
```yaml
# .github/workflows/ci.yml
- name: Generate Ent Code
  run: go generate ./infrastructure/ent
- name: Verify no uncommitted changes
  run: git diff --exit-code
```

---

## 2. アーキテクチャ・設計の改善

### 2.1 アーキテクチャ違反の修正

**問題**: `handler/choice_handler.go:12-19`
```go
type ChoiceHandler struct {
    client *ent.Client  // ❌ Infrastructure層への直接依存
}
```

**影響度**: 🟡 **高**
オニオンアーキテクチャの原則に違反し、レイヤー間の疎結合が損なわれています。

**改善案**:
```go
// handler/choice_handler.go
type ChoiceHandler struct {
    service *service.ChoiceService  // ✅ Service層を経由
}

func NewChoiceHandler(service *service.ChoiceService) *ChoiceHandler {
    return &ChoiceHandler{service: service}
}
```

**実装手順**:
1. `service/choice.go`を作成してChoiceServiceを実装
2. `ChoiceHandler`の依存を`*ent.Client`から`*service.ChoiceService`に変更
3. `cmd/api/main.go`で依存注入を更新

---

### 2.2 命名規則の統一

**問題**: domain/entity/question.go:9-13
```go
const MAX_CHOICE_COUNT = 4         // ❌ UPPER_SNAKE_CASE
var ERR_MAX_CHOICE_COUNT = errors.New("...")  // ❌ varには不適切
```

**影響度**: 🟢 **低**
機能には影響しませんが、Goの慣習に反しています。

**改善案**:
```go
// 定数
const MaxChoiceCount = 4  // ✅ MixedCase

// エラー変数
var (
    ErrMaxChoiceCount            = errors.New("choice count exceeds maximum")
    ErrNoChoice                  = errors.New("no choice provided")
    ErrInvalidCorrectChoiceCount = errors.New("invalid correct choice count")
)
```

---

### 2.3 コードの重複排除

**問題**: `infrastructure/persistence/question.go`
- FindAll（行37-43）とFindByID（行75-80）でエンティティ変換ロジックが重複

**改善案**:
```go
// 共通の変換ロジックを抽出
func (q QuestionRepository) toEntity(qs *ent.Question, choices []entity.Choice) entity.Question {
    return entity.Question{
        UID:           entity.ToUID(qs.UID),
        ReferenceCode: lo.FromPtr(qs.ReferenceCode),
        Title:         qs.Title,
        Content:       qs.Content,
        Choices:       choices,
    }
}

// 使用例
func (q QuestionRepository) FindAll(ctx context.Context) ([]entity.Question, error) {
    qs, err := q.db.Question.Query().All(ctx)
    // ...
    for _, q := range qs {
        questions = append(questions, q.toEntity(q, []entity.Choice{}))
    }
    return questions, nil
}
```

---

## 3. 未実装機能の完成

### 3.1 APIエンドポイントの実装

以下のエンドポイントがスタブ実装のままです：

| エンドポイント | ファイル | 行番号 | 優先度 |
|-------------|---------|--------|--------|
| `POST /api/questions` | handler/question_handler.go | 55-65 | 🟡 高 |
| `PUT /api/questions/:id` | handler/question_handler.go | 96-98 | 🟡 高 |
| `DELETE /api/questions/:id` | handler/question_handler.go | 101-103 | 🟡 高 |
| `GET /api/questions/:id/choices` | handler/choice_handler.go | 23-25 | 🟡 高 |
| `POST /api/questions/:id/choices` | handler/choice_handler.go | 28-30 | 🟡 高 |
| `PUT /api/choices/:id` | handler/choice_handler.go | 33-35 | 🟢 中 |
| `DELETE /api/choices/:id` | handler/choice_handler.go | 38-40 | 🟢 中 |

**実装例（問題作成）**:
```go
// handler/question_handler.go
func (h *QuestionHandler) Create(c echo.Context) error {
    req := new(request.CreateQuestionRequest)
    if err := c.Bind(req); err != nil {
        return echo.NewHTTPError(http.StatusBadRequest, err.Error())
    }
    if err := c.Validate(req); err != nil {
        return echo.NewHTTPError(http.StatusBadRequest, err.Error())
    }

    // DTOに変換してサービス層へ
    questionReq := dto.CreateQuestionRequest{
        ReferenceCode: req.ReferenceCode,
        Title:         req.Title,
        Content:       req.Content,
        Choices:       req.Choices,
    }

    result, err := h.service.Create(c.Request().Context(), questionReq)
    if err != nil {
        return h.handleServiceError(c, err)
    }

    response := response.NewQuestionResponse(result)
    return c.JSON(http.StatusCreated, response)
}
```

---

## 4. テストカバレッジの改善

### 4.1 現在のカバレッジ状況

| パッケージ | カバレッジ | テスト状況 |
|-----------|----------|-----------|
| `handler/` | 0% | ❌ テストなし |
| `service/` | ~80% | 🟡 部分的 |
| `infrastructure/persistence/` | ~60% | 🟡 部分的 |
| `domain/entity/` | 0% | ❌ テストなし |
| `cmd/cli/` | 0% | ❌ テストなし |
| **全体** | **~25-30%** | 🔴 不十分 |

### 4.2 優先的に追加すべきテスト

#### 4.2.1 Handler層のテスト

**ファイル**: `handler/question_handler_test.go`（新規作成）

```go
package handlers

import (
    "net/http"
    "net/http/httptest"
    "testing"

    "github.com/labstack/echo/v4"
    "github.com/stretchr/testify/assert"
    "github.com/stretchr/testify/mock"
)

func TestQuestionHandler_GetAll(t *testing.T) {
    // Arrange
    mockService := new(MockQuestionService)
    handler := NewQuestionHandler(mockService)

    e := echo.New()
    req := httptest.NewRequest(http.MethodGet, "/api/questions", nil)
    rec := httptest.NewRecorder()
    c := e.NewContext(req, rec)

    mockService.On("FindAll", mock.Anything).
        Return([]dto.QuestionResult{}, nil)

    // Act
    err := handler.GetAll(c)

    // Assert
    assert.NoError(t, err)
    assert.Equal(t, http.StatusOK, rec.Code)
    mockService.AssertExpectations(t)
}
```

#### 4.2.2 Entity層のバリデーションテスト

**ファイル**: `domain/entity/question_test.go`（新規作成）

```go
package entity

import (
    "testing"

    "github.com/stretchr/testify/assert"
)

func TestNewQuestion_ValidInput(t *testing.T) {
    // Arrange
    choices := []Choice{
        {Content: "選択肢1", IsCorrect: true},
        {Content: "選択肢2", IsCorrect: false},
    }

    // Act
    q, err := NewQuestion("REF001", "タイトル", "内容", choices)

    // Assert
    assert.NoError(t, err)
    assert.Equal(t, "REF001", q.ReferenceCode)
    assert.Len(t, q.Choices, 2)
}

func TestNewQuestion_NoChoices(t *testing.T) {
    // Act
    _, err := NewQuestion("REF001", "タイトル", "内容", []Choice{})

    // Assert
    assert.ErrorIs(t, err, ERR_NO_CHOICE)
}

func TestNewQuestion_TooManyChoices(t *testing.T) {
    // Arrange
    choices := make([]Choice, MAX_CHOICE_COUNT+1)
    for i := range choices {
        choices[i] = Choice{Content: "選択肢", IsCorrect: i == 0}
    }

    // Act
    _, err := NewQuestion("REF001", "タイトル", "内容", choices)

    // Assert
    assert.ErrorIs(t, err, ERR_MAX_CHOICE_COUNT)
}
```

---

## 5. セキュリティの強化

### 5.1 認証・認可の実装

**現状**: 認証機構が実装されていない
**影響度**: 🔴 **緊急**

**推奨実装**:

1. JWT認証ミドルウェアの追加

```go
// middleware/auth.go (新規作成)
package middleware

import (
    "net/http"

    "github.com/golang-jwt/jwt"
    "github.com/labstack/echo/v4"
)

func JWTAuth(secret string) echo.MiddlewareFunc {
    return func(next echo.HandlerFunc) echo.HandlerFunc {
        return func(c echo.Context) error {
            token := c.Request().Header.Get("Authorization")
            if token == "" {
                return echo.NewHTTPError(http.StatusUnauthorized, "missing token")
            }

            // トークン検証ロジック
            claims := jwt.MapClaims{}
            _, err := jwt.ParseWithClaims(token, claims, func(token *jwt.Token) (interface{}, error) {
                return []byte(secret), nil
            })
            if err != nil {
                return echo.NewHTTPError(http.StatusUnauthorized, "invalid token")
            }

            c.Set("user", claims)
            return next(c)
        }
    }
}
```

2. ルート設定での適用

```go
// handler/routes.go
func RegisterRoutes(e *echo.Echo, client *ent.Client) {
    api := e.Group("/api")

    // 公開エンドポイント
    public := api.Group("")
    public.GET("/questions", questionHandler.GetAll)
    public.GET("/questions/:id", questionHandler.GetByID)

    // 認証が必要なエンドポイント
    protected := api.Group("")
    protected.Use(middleware.JWTAuth(os.Getenv("JWT_SECRET")))
    protected.POST("/questions", questionHandler.Create)
    protected.PUT("/questions/:id", questionHandler.Update)
    protected.DELETE("/questions/:id", questionHandler.Delete)
}
```

---

### 5.2 CORS設定の厳格化

**問題**: `cmd/api/main.go:41`
```go
e.Use(middleware.CORS())  // ❌ すべてのオリジンを許可
```

**影響度**: 🟡 **高**
任意のウェブサイトからAPIを呼び出せる状態です。

**改善案**:
```go
// cmd/api/main.go
e.Use(middleware.CORSWithConfig(middleware.CORSConfig{
    AllowOrigins: []string{
        "https://yourdomain.com",
        "https://app.yourdomain.com",
    },
    AllowMethods: []string{http.MethodGet, http.MethodPost, http.MethodPut, http.MethodDelete},
    AllowHeaders: []string{
        echo.HeaderOrigin,
        echo.HeaderContentType,
        echo.HeaderAccept,
        echo.HeaderAuthorization,
    },
    AllowCredentials: true,
    MaxAge:           3600,
}))
```

開発環境用の設定:
```go
// 環境変数で制御
if cfg.Environment == "development" {
    e.Use(middleware.CORS())
} else {
    e.Use(middleware.CORSWithConfig(/* 本番用設定 */))
}
```

---

### 5.3 入力値検証の強化

#### 5.3.1 UIDフォーマットの検証

**問題**: `handler/question_handler.go:41-44`
```go
id := c.Param("id")
if id == "" {
    return echo.NewHTTPError(http.StatusBadRequest, "id is required")
}
// ❌ UUID形式の検証がない
```

**改善案**:
```go
// pkg/validation/uid.go (新規作成)
package validation

import (
    "github.com/google/uuid"
)

func IsValidUID(id string) bool {
    _, err := uuid.Parse(id)
    return err == nil
}

// handler/question_handler.go
func (h *QuestionHandler) GetByID(c echo.Context) error {
    id := c.Param("id")
    if id == "" || !validation.IsValidUID(id) {
        return echo.NewHTTPError(http.StatusBadRequest, "invalid id format")
    }
    // ...
}
```

#### 5.3.2 リクエストサイズの制限

**改善案**:
```go
// cmd/api/main.go
e.Use(middleware.BodyLimit("2M"))  // 2MBに制限
```

#### 5.3.3 XSS対策

**問題**: HTMLタグのサニタイゼーションがない

**推奨**: バリデーションタグの追加
```go
// handler/request/question.go
type CreateQuestionRequest struct {
    ReferenceCode string        `json:"reference_code" validate:"omitempty,max=50,excludesall=<>"`
    Title         string        `json:"title" validate:"required,max=200,excludesall=<>"`
    Content       string        `json:"content" validate:"required,max=800,excludesall=<>"`
    Choices       []ChoiceInput `json:"choices" validate:"required,min=2,max=4,dive"`
}
```

または、専用のサニタイズライブラリの使用:
```go
import "github.com/microcosm-cc/bluemonday"

func sanitizeHTML(input string) string {
    p := bluemonday.StrictPolicy()
    return p.Sanitize(input)
}
```

---

## 6. パフォーマンスの最適化

### 6.1 N+1問題の解決

**問題**: `infrastructure/persistence/question.go:25-47`
```go
func (q QuestionRepository) FindAll(ctx context.Context) ([]entity.Question, error) {
    qs, err := q.db.Question.Query().All(ctx)  // ❌ WithChoices()なし
    // ...
    Choices: []entity.Choice{},  // 常に空配列
}
```

**影響**: クライアントが問題ごとに個別リクエストを発行する必要がある

**改善案**:

Option 1: 選択肢を含めて取得
```go
func (q QuestionRepository) FindAll(ctx context.Context) ([]entity.Question, error) {
    qs, err := q.db.Question.
        Query().
        WithChoices().  // ✅ 選択肢も一緒に取得
        All(ctx)
    // ...変換処理で choices を含める
}
```

Option 2: 専用メソッドの作成
```go
// 選択肢なしのリスト（一覧表示用）
func (q QuestionRepository) FindAllSummary(ctx context.Context) ([]entity.Question, error)

// 選択肢ありのリスト（詳細表示用）
func (q QuestionRepository) FindAllWithChoices(ctx context.Context) ([]entity.Question, error)
```

**推奨**: Option 2（用途に応じて使い分け）

---

### 6.2 データベースインデックスの確認

**確認項目**:
```sql
-- UIDカラムにインデックスがあるか確認
CREATE INDEX IF NOT EXISTS idx_question_uid ON questions(uid);
CREATE INDEX IF NOT EXISTS idx_choice_uid ON choices(uid);
CREATE INDEX IF NOT EXISTS idx_choice_question_id ON choices(question_id);
```

**Entスキーマでの定義**:
```go
// infrastructure/ent/schema/question.go
func (Question) Indexes() []ent.Index {
    return []ent.Index{
        index.Fields("uid").Unique(),
    }
}
```

---

### 6.3 HTTPキャッシュの実装

**推奨**: ETagとCache-Controlヘッダーの追加

```go
// middleware/cache.go (新規作成)
package middleware

import (
    "crypto/md5"
    "encoding/hex"
    "net/http"

    "github.com/labstack/echo/v4"
)

func CacheControl(maxAge int) echo.MiddlewareFunc {
    return func(next echo.HandlerFunc) echo.HandlerFunc {
        return func(c echo.Context) error {
            // リクエストを処理
            if err := next(c); err != nil {
                return err
            }

            // GETリクエストのみキャッシュ
            if c.Request().Method == http.MethodGet {
                c.Response().Header().Set("Cache-Control",
                    fmt.Sprintf("public, max-age=%d", maxAge))
            }

            return nil
        }
    }
}

// 使用例
api.GET("/questions", questionHandler.GetAll, CacheControl(3600))  // 1時間キャッシュ
```

---

## 7. コード品質の向上

### 7.1 リポジトリインターフェースの改善

**問題**: `domain/repository/question.go`
```go
type QuestionRepository interface {
    FindAll(ctx context.Context) ([]entity.Question, error)
    FindByID(ctx context.Context, id entity.UID) (entity.Question, error)
    Save(ctx context.Context, id entity.UID, referenceCode, title, content string) error
    // ❌ Update, Deleteメソッドがない
    // ❌ Saveのシグネチャが不自然
}
```

**改善案**:
```go
type QuestionRepository interface {
    FindAll(ctx context.Context) ([]entity.Question, error)
    FindByID(ctx context.Context, id entity.UID) (entity.Question, error)
    Create(ctx context.Context, q entity.Question) error
    Update(ctx context.Context, q entity.Question) error
    Delete(ctx context.Context, id entity.UID) error
}
```

---

### 7.2 エラーメッセージの多言語化対応

**現状**: エラーメッセージが日本語ハードコード

**推奨**: メッセージキーとi18nライブラリの使用

```go
// pkg/i18n/messages.go (新規作成)
package i18n

var Messages = map[string]map[string]string{
    "ja": {
        "error.question.not_found": "指定された問題が見つかりません",
        "error.validation.required": "%sは必須です",
    },
    "en": {
        "error.question.not_found": "Question not found",
        "error.validation.required": "%s is required",
    },
}

func GetMessage(lang, key string, args ...interface{}) string {
    msg, ok := Messages[lang][key]
    if !ok {
        return key
    }
    return fmt.Sprintf(msg, args...)
}

// 使用例
return serviceErrors.NewNotFoundError(
    i18n.GetMessage(lang, "error.question.not_found"),
    err,
)
```

---

### 7.3 ログ戦略の改善

**現状**: ロギングレベルと構造化ログが未実装

**推奨**: Zapまたはzerolog の導入

```go
// pkg/logger/logger.go (新規作成)
package logger

import (
    "go.uber.org/zap"
)

var Log *zap.Logger

func Init(env string) error {
    var err error
    if env == "production" {
        Log, err = zap.NewProduction()
    } else {
        Log, err = zap.NewDevelopment()
    }
    return err
}

// 使用例
logger.Log.Error("failed to find question",
    zap.String("question_id", id),
    zap.Error(err),
    zap.String("trace_id", traceID),
)
```

---

## 8. ドキュメンテーションの改善

### 8.1 API仕様書（OpenAPI/Swagger）

**現状**: READMEのみ
**推奨**: swaggoの導入

```bash
# インストール
go install github.com/swaggo/swag/cmd/swag@latest

# アノテーションの追加例
// @title MedicalQuest API
// @version 1.0
// @description 医療問題出題APIサーバー
// @host localhost:8080
// @BasePath /api

// @Summary 全ての問題を取得
// @Description 問題一覧を取得します
// @Tags questions
// @Accept json
// @Produce json
// @Success 200 {object} response.QuestionsResponse
// @Failure 500 {object} echo.HTTPError
// @Router /questions [get]
func (h *QuestionHandler) GetAll(c echo.Context) error {
    // ...
}

# ドキュメント生成
swag init -g cmd/api/main.go
```

---

### 8.2 ADR（Architecture Decision Records）

**推奨**: 重要な設計判断を記録

```markdown
# ADR-001: オニオンアーキテクチャの採用

## Status
採用済み

## Context
アプリケーションの保守性と拡張性を高めるため、アーキテクチャパターンを選定する必要がある。

## Decision
オニオンアーキテクチャを採用し、以下のレイヤー構成とする：
- Domain層（Entity, Repository Interface）
- Service層（Use Case実装）
- Infrastructure層（Repository実装、DB接続）
- Handler層（HTTP/CLI）

## Consequences
### ポジティブ
- ビジネスロジックがインフラストラクチャから独立
- テストしやすい構造

### ネガティブ
- レイヤー間の移動でボイラープレートが増加
```

---

## 9. CI/CDの改善

### 9.1 GitHub Actionsの拡充

**現状**: テスト、Lint、ビルドのみ

**推奨追加項目**:

```yaml
# .github/workflows/ci.yml
name: CI

on:
  pull_request:
    branches: [main]

jobs:
  # 既存ジョブ...

  security-scan:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v4
      - name: Run Gosec Security Scanner
        uses: securego/gosec@master
        with:
          args: ./...

  coverage:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v4
      - name: Run tests with coverage
        run: |
          go test -race -coverprofile=coverage.out -covermode=atomic ./...
          go tool cover -func=coverage.out
      - name: Upload coverage to Codecov
        uses: codecov/codecov-action@v3
        with:
          files: ./coverage.out

  dependency-check:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v4
      - name: Check for vulnerabilities
        run: |
          go install golang.org/x/vuln/cmd/govulncheck@latest
          govulncheck ./...
```

---

## 10. 実装優先順位とロードマップ

### Phase 1: 緊急対応（1-2週間）

1. ✅ **Entコード生成エラーの修正**
   - `./scripts/generate-ent.sh`実行
   - CI/CDへの自動生成追加

2. ✅ **ChoiceHandlerのアーキテクチャ修正**
   - ChoiceServiceの実装
   - ハンドラーの依存関係変更

3. ✅ **CORS設定の厳格化**
   - 環境別の設定追加

4. ✅ **基本的な入力検証強化**
   - UID形式のバリデーション
   - リクエストサイズ制限

---

### Phase 2: 品質向上（3-4週間）

5. ✅ **未実装APIエンドポイントの完成**
   - Create, Update, Delete の実装
   - Choice系エンドポイントの実装

6. ✅ **テストカバレッジの向上（目標: 70%以上）**
   - Handler層のテスト追加
   - Entity層のテスト追加
   - 統合テストの追加

7. ✅ **エラーハンドリングの統一**
   - エラーコードの追加
   - メッセージキーの導入

8. ✅ **ロギング改善**
   - 構造化ログの導入
   - ログレベルの適切な設定

---

### Phase 3: セキュリティ強化（2-3週間）

9. ✅ **認証・認可の実装**
   - JWT認証の導入
   - ロールベースのアクセス制御

10. ✅ **セキュリティスキャンの導入**
    - Gosecの導入
    - 依存関係の脆弱性チェック

---

### Phase 4: パフォーマンス最適化（2週間）

11. ✅ **N+1問題の解決**
    - FindAllWithChoicesメソッドの実装
    - 適切なデータ取得戦略の選択

12. ✅ **キャッシング戦略の実装**
    - HTTPキャッシュヘッダー
    - （必要に応じて）Redisキャッシュ

13. ✅ **データベースインデックスの最適化**

---

### Phase 5: ドキュメントと運用（1-2週間）

14. ✅ **API仕様書の作成**
    - Swagger/OpenAPIドキュメント

15. ✅ **運用ドキュメントの整備**
    - デプロイ手順
    - トラブルシューティングガイド

---

## 11. まとめ

MedicalQuestは堅実なアーキテクチャ設計を持つプロジェクトですが、以下の改善により、**プロダクションレディな状態**に近づけることができます：

### 必須対応（Phase 1-2）
- ✅ Entコード生成エラーの修正
- ✅ アーキテクチャ違反の修正
- ✅ 未実装機能の完成
- ✅ テストカバレッジの向上
- ✅ セキュリティの基本対策

### 推奨対応（Phase 3-5）
- 🔄 認証・認可の実装
- 🔄 パフォーマンス最適化
- 🔄 ドキュメント整備

### 長期的な改善
- 🔄 多言語化対応
- 🔄 監視・アラート体制
- 🔄 パフォーマンスモニタリング

---

## 付録: チェックリスト

### コードレビュー時のチェックポイント

#### アーキテクチャ
- [ ] レイヤー間の依存方向が正しいか
- [ ] ドメインロジックがEntity/Service層に集約されているか
- [ ] Handler層がインフラ層に直接依存していないか

#### セキュリティ
- [ ] 入力値のバリデーションが実装されているか
- [ ] エラーメッセージに機密情報が含まれていないか
- [ ] 認証が必要なエンドポイントに認証チェックがあるか
- [ ] SQLインジェクションのリスクがないか（ORMの適切な使用）

#### テスト
- [ ] ユニットテストが書かれているか
- [ ] テストカバレッジが十分か（目標: 70%以上）
- [ ] エッジケースのテストがあるか

#### パフォーマンス
- [ ] N+1問題が発生していないか
- [ ] 不要なデータベースクエリがないか
- [ ] 適切なインデックスが設定されているか

#### コード品質
- [ ] 命名規則が統一されているか
- [ ] コードの重複がないか
- [ ] エラーハンドリングが適切か
- [ ] ログが適切に出力されているか

---

**生成者**: Claude Code AI
**最終更新**: 2025-11-13
