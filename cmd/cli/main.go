package main

import (
    "fmt"
    "os"

    "github.com/spf13/cobra"
)

var (
    apiURL string
)

func main() {
    var rootCmd = &cobra.Command{
        Use:   "medicalquest-cli",
        Short: "MedicalQuest CLI - APIからクイズデータを取得して表示",
        Long:  `MedicalQuest CLIはAPIサーバーからクイズの質問と選択肢を取得して表示するツールです。`,
    }

    // グローバルフラグ
    rootCmd.PersistentFlags().StringVar(&apiURL, "api-url", "http://localhost:8080", "APIサーバーのURL")

    // コマンドを追加
    rootCmd.AddCommand(listCmd())
    rootCmd.AddCommand(getCmd())

    if err := rootCmd.Execute(); err != nil {
        fmt.Fprintln(os.Stderr, err)
        os.Exit(1)
    }
}