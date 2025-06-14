package main

import (
    "bufio"
    "encoding/csv"
    "fmt"
    "os"
    "time"
)

func main() {
    fmt.Println("MedicalQuest CLI")
    fmt.Println("================")
    
    scanner := bufio.NewScanner(os.Stdin)
    
    // Simulate a quiz session
    correctAnswers := 0
    totalQuestions := 0
    results := [][]string{
        {"問題番号", "問題", "あなたの回答", "正解"},
    }
    
    // In a real implementation, we would fetch questions from the API
    // For now, we'll use a simple hard-coded example
    questions := []struct {
        ID       int
        Title    string
        Content  string
        Choices  []string
        Correct  int
    }{
        {
            ID:      1,
            Title:   "基本的な医学知識",
            Content: "心臓の弁は全部で何枚ありますか？",
            Choices: []string{"2枚", "3枚", "4枚", "5枚"},
            Correct: 2, // 0-indexed, so this is "4枚"
        },
        {
            ID:      2,
            Title:   "解剖学",
            Content: "ヒトの脊椎は何個の骨で構成されていますか？",
            Choices: []string{"24個", "26個", "33個", "34個"},
            Correct: 2, // 0-indexed, so this is "33個"
        },
    }
    
    // Start quiz
    fmt.Println("\n問題の開始")
    
    // Loop through questions
    for i, q := range questions {
        totalQuestions++
        
        fmt.Printf("\n問題 %d: %s\n", i+1, q.Title)
        fmt.Printf("%s\n\n", q.Content)
        
        // Display choices
        for j, choice := range q.Choices {
            fmt.Printf("%d. %s\n", j+1, choice)
        }
        
        // Get user answer
        var userAnswer int
        fmt.Print("\n回答を選択してください (1-4): ")
        for {
            scanner.Scan()
            _, err := fmt.Sscanf(scanner.Text(), "%d", &userAnswer)
            if err == nil && userAnswer >= 1 && userAnswer <= len(q.Choices) {
                break
            }
            fmt.Print("無効な入力です。もう一度入力してください (1-4): ")
        }
        
        // Check if correct
        userAnswerZeroIndex := userAnswer - 1
        isCorrect := userAnswerZeroIndex == q.Correct
        
        if isCorrect {
            correctAnswers++
            fmt.Println("正解！")
        } else {
            fmt.Printf("不正解。正解は: %d. %s\n", q.Correct+1, q.Choices[q.Correct])
        }
        
        // Add to results
        results = append(results, []string{
            fmt.Sprintf("%d", q.ID),
            q.Content,
            q.Choices[userAnswerZeroIndex],
            q.Choices[q.Correct],
        })
    }
    
    // Display summary
    fmt.Println("\n結果サマリ")
    fmt.Println("===========")
    fmt.Printf("実施日時: %s\n", time.Now().Format("2006-01-02 15:04:05"))
    fmt.Printf("正答数: %d / %d\n", correctAnswers, totalQuestions)
    fmt.Printf("正答率: %.1f%%\n", float64(correctAnswers)/float64(totalQuestions)*100)
    
    // Save results to CSV
    filename := fmt.Sprintf("quiz_results_%s.csv", time.Now().Format("20060102_150405"))
    saveResultsToCSV(filename, results)
    fmt.Printf("\n結果は %s に保存されました。\n", filename)
}

func saveResultsToCSV(filename string, data [][]string) error {
    file, err := os.Create(filename)
    if err != nil {
        return err
    }
    defer file.Close()
    
    writer := csv.NewWriter(file)
    defer writer.Flush()
    
    // Add header row with summary
    writer.Write([]string{"MedicalQuest クイズ結果"})
    writer.Write([]string{time.Now().Format("2006-01-02 15:04:05")})
    writer.Write([]string{})
    
    // Write data
    for _, row := range data {
        writer.Write(row)
    }
    
    return nil
}
