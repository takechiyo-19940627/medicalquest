package main

import (
    "fmt"
    "os"
    "text/tabwriter"

    "github.com/spf13/cobra"
)

// listコマンド - 質問一覧を取得して表示
func listCmd() *cobra.Command {
    cmd := &cobra.Command{
        Use:   "list",
        Short: "質問一覧を表示",
        Long:  `APIサーバーから全ての質問を取得してテーブル形式で表示します。`,
        RunE:  runList,
    }

    return cmd
}

func runList(cmd *cobra.Command, args []string) error {
    client := NewAPIClient(apiURL)

    // 質問一覧を取得
    questions, err := client.GetQuestions()
    if err != nil {
        return fmt.Errorf("質問一覧の取得に失敗しました: %w", err)
    }

    if len(questions) == 0 {
        fmt.Println("質問が見つかりませんでした。")
        return nil
    }

    // テーブル形式で表示
    w := tabwriter.NewWriter(os.Stdout, 0, 0, 2, ' ', 0)
    fmt.Fprintln(w, "UID\t参照コード\tタイトル\t内容")
    fmt.Fprintln(w, "---\t--------\t------\t----")

    for _, q := range questions {
        // 長い文字列は切り詰める
        content := q.Content
        if len(content) > 50 {
            content = content[:47] + "..."
        }
        
        fmt.Fprintf(w, "%s\t%s\t%s\t%s\n", 
            q.UID, 
            q.ReferenceCode, 
            q.Title, 
            content)
    }

    w.Flush()
    fmt.Printf("\n合計: %d 件の質問\n", len(questions))

    return nil
}