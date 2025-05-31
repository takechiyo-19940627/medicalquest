package persistence

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/takechiyo-19940627/medicalquest/domain/entity"
	"github.com/takechiyo-19940627/medicalquest/infrastructure/ent/enttest"
	"github.com/takechiyo-19940627/medicalquest/infrastructure/ent/question"

	_ "github.com/mattn/go-sqlite3"
)

func TestQuestionRepository_FindAll(t *testing.T) {
	tests := []struct {
		name          string
		setupData     []struct {
			refCode string
			title   string
			content string
		}
		expectedCount int
		expectedQuestions []struct {
			refCode string
			title   string
			content string
		}
	}{
		{
			name: "複数の質問を正常に取得",
			setupData: []struct {
				refCode string
				title   string
				content string
			}{
				{"REF001", "質問1", "質問1の内容"},
				{"REF002", "質問2", "質問2の内容"},
			},
			expectedCount: 2,
			expectedQuestions: []struct {
				refCode string
				title   string
				content string
			}{
				{"REF001", "質問1", "質問1の内容"},
				{"REF002", "質問2", "質問2の内容"},
			},
		},
		{
			name:          "データが存在しない場合",
			setupData:     []struct {
				refCode string
				title   string
				content string
			}{},
			expectedCount: 0,
			expectedQuestions: []struct {
				refCode string
				title   string
				content string
			}{},
		},
		{
			name: "単一の質問を取得",
			setupData: []struct {
				refCode string
				title   string
				content string
			}{
				{"REF003", "単一質問", "単一質問の内容"},
			},
			expectedCount: 1,
			expectedQuestions: []struct {
				refCode string
				title   string
				content string
			}{
				{"REF003", "単一質問", "単一質問の内容"},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// enttest を使用してインメモリSQLiteのテストクライアントを作成
			client := enttest.Open(t, "sqlite3", "file:ent?mode=memory&cache=shared&_fk=1")
			defer client.Close()

			// リポジトリのインスタンスを作成
			repo := NewQuestionRepository(client)
			ctx := context.Background()

			// テストデータの準備と保存
			for _, data := range tt.setupData {
				// テスト対象のリポジトリを使わず、直接Entクライアントでテストデータを作成
				uid := entity.GenerateUID()
				refCode := &data.refCode
				if data.refCode == "" {
					refCode = nil
				}
				_, err := client.Question.
					Create().
					SetUID(uid.String()).
					SetNillableReferenceCode(refCode).
					SetTitle(data.title).
					SetContent(data.content).
					Save(ctx)
				assert.NoError(t, err)
			}

			// テスト対象のメソッドを実行
			questions, err := repo.FindAll(ctx)

			// アサーション
			assert.NoError(t, err)
			assert.Len(t, questions, tt.expectedCount)

			// 期待される質問が含まれているかを検証
			for _, expected := range tt.expectedQuestions {
				found := false
				for _, q := range questions {
					if q.ReferenceCode == expected.refCode && q.Title == expected.title && q.Content == expected.content {
						found = true
						break
					}
				}
				assert.True(t, found, "期待される質問が見つかりませんでした: %v", expected)
			}
		})
	}
}

func TestQuestionRepository_FindByID(t *testing.T) {
	// 現状、FindByIDは実装されていないため、スキップ
	t.Skip("FindByID method is not implemented yet")
}

func TestQuestionRepository_Save(t *testing.T) {
	tests := []struct {
		name        string
		refCode     string
		title       string
		content     string
		expectError bool
	}{
		{
			name:        "正常な質問の保存",
			refCode:     "REF003",
			title:       "新しい質問",
			content:     "新しい質問の内容",
			expectError: false,
		},
		{
			name:        "空のタイトルで保存（バリデーションエラー）",
			refCode:     "REF004",
			title:       "",
			content:     "内容のみ",
			expectError: true,
		},
		{
			name:        "空の内容で保存（バリデーションエラー）",
			refCode:     "REF005",
			title:       "タイトルのみ",
			content:     "",
			expectError: true,
		},
		{
			name:        "長いタイトルと内容で保存",
			refCode:     "REF006",
			title:       "これは非常に長いタイトルです。長いタイトルのテストケース",
			content:     "これは非常に長い内容です。長い内容のテストケースとして使用します。",
			expectError: false,
		},
		{
			name:        "参照コードなしで保存",
			refCode:     "",
			title:       "参照コードなし",
			content:     "参照コードがない場合のテスト",
			expectError: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// enttest を使用してインメモリSQLiteのテストクライアントを作成
			client := enttest.Open(t, "sqlite3", "file:ent?mode=memory&cache=shared&_fk=1")
			defer client.Close()

			// リポジトリのインスタンスを作成
			repo := NewQuestionRepository(client)
			ctx := context.Background()

			// テストデータの準備
			uid := entity.GenerateUID()

			// テスト対象のメソッドを実行
			err := repo.Save(ctx, uid, tt.refCode, tt.title, tt.content)

			// アサーション
			if tt.expectError {
				assert.Error(t, err)
				return
			}

			assert.NoError(t, err)

			// 保存されたデータを検証
			savedQuestion, err := client.Question.
				Query().
				Where(
					question.UID(uid.String()),
				).
				Unique(true).
				First(ctx)
			assert.NoError(t, err)
			assert.Equal(t, uid.String(), savedQuestion.UID)
			assert.Equal(t, tt.refCode, *savedQuestion.ReferenceCode)
			assert.Equal(t, tt.title, savedQuestion.Title)
			assert.Equal(t, tt.content, savedQuestion.Content)
		})
	}
}
