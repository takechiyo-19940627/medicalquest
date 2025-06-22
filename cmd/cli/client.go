package main

import (
    "encoding/json"
    "fmt"
    "io"
    "net/http"
    "time"
)

// APIクライアント
type APIClient struct {
    baseURL    string
    httpClient *http.Client
}

// 新しいAPIクライアントを作成
func NewAPIClient(baseURL string) *APIClient {
    return &APIClient{
        baseURL: baseURL,
        httpClient: &http.Client{
            Timeout: 10 * time.Second,
        },
    }
}

// HTTPリクエストを実行
func (c *APIClient) doRequest(method, path string, body io.Reader) (*http.Response, error) {
    url := c.baseURL + path
    req, err := http.NewRequest(method, url, body)
    if err != nil {
        return nil, fmt.Errorf("リクエストの作成に失敗しました: %w", err)
    }

    req.Header.Set("Content-Type", "application/json")
    req.Header.Set("Accept", "application/json")

    resp, err := c.httpClient.Do(req)
    if err != nil {
        return nil, fmt.Errorf("リクエストの送信に失敗しました: %w", err)
    }

    return resp, nil
}

// レスポンス型
type Question struct {
    UID           string   `json:"uid"`
    ReferenceCode string   `json:"reference_code"`
    Title         string   `json:"title"`
    Content       string   `json:"content"`
    Choices       []Choice `json:"choices,omitempty"`
}

type Choice struct {
    UID       string `json:"uid"`
    Content   string `json:"content"`
    IsCorrect bool   `json:"is_correct"`
}

type ListResponse struct {
    Data []Question `json:"data"`
}

type QuestionResponse struct {
    Data Question `json:"data"`
}

// 質問一覧を取得
func (c *APIClient) GetQuestions() ([]Question, error) {
    resp, err := c.doRequest("GET", "/api/questions", nil)
    if err != nil {
        return nil, err
    }
    defer resp.Body.Close()

    if resp.StatusCode != http.StatusOK {
        body, _ := io.ReadAll(resp.Body)
        return nil, fmt.Errorf("APIエラー: ステータスコード %d, レスポンス: %s", resp.StatusCode, body)
    }

    var result ListResponse
    if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
        return nil, fmt.Errorf("レスポンスのデコードに失敗しました: %w", err)
    }

    return result.Data, nil
}

// 特定の質問を取得
func (c *APIClient) GetQuestion(id string) (*Question, error) {
    resp, err := c.doRequest("GET", fmt.Sprintf("/api/questions/%s", id), nil)
    if err != nil {
        return nil, err
    }
    defer resp.Body.Close()

    if resp.StatusCode != http.StatusOK {
        body, _ := io.ReadAll(resp.Body)
        return nil, fmt.Errorf("APIエラー: ステータスコード %d, レスポンス: %s", resp.StatusCode, body)
    }

    var result QuestionResponse
    if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
        return nil, fmt.Errorf("レスポンスのデコードに失敗しました: %w", err)
    }

    return &result.Data, nil
}