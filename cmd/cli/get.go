package main

import (
    "fmt"
    "os"
    "text/tabwriter"

    "github.com/spf13/cobra"
)

// getコマンド - 特定の質問と選択肢を表示
func getCmd() *cobra.Command {
    cmd := &cobra.Command{
        Use:   "get [質問ID]",
        Short: "特定の質問とその選択肢を表示",
        Long:  `指定されたIDの質問と、その質問に関連する全ての選択肢を表示します。`,
        Args:  cobra.ExactArgs(1),
        RunE:  runGet,
    }

    return cmd
}

func runGet(cmd *cobra.Command, args []string) error {
    questionID := args[0]
    client := NewAPIClient(apiURL)

    // 質問を取得
    question, err := client.GetQuestion(questionID)
    if err != nil {
        return fmt.Errorf("質問の取得に失敗しました: %w", err)
    }

    // 質問情報を表示
    fmt.Println("==============================")
    fmt.Println("質問情報")
    fmt.Println("==============================")
    fmt.Printf("UID:         %s\n", question.UID)
    fmt.Printf("参照コード:   %s\n", question.ReferenceCode)
    fmt.Printf("タイトル:     %s\n", question.Title)
    fmt.Printf("内容:        %s\n", question.Content)

    // 選択肢を表示
    if len(question.Choices) > 0 {
        fmt.Println("\n==============================")
        fmt.Println("選択肢")
        fmt.Println("==============================")

        w := tabwriter.NewWriter(os.Stdout, 0, 0, 2, ' ', 0)
        fmt.Fprintln(w, "番号\tUID\t内容\t正解")
        fmt.Fprintln(w, "----\t---\t----\t----")

        for i, choice := range question.Choices {
            correct := "×"
            if choice.IsCorrect {
                correct = "○"
            }
            fmt.Fprintf(w, "%d\t%s\t%s\t%s\n", 
                i+1, 
                choice.UID, 
                choice.Content, 
                correct)
        }
        w.Flush()

        fmt.Printf("\n選択肢数: %d\n", len(question.Choices))
    } else {
        fmt.Println("\nこの質問には選択肢がありません。")
    }

    return nil
}