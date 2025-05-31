package entity

type Choice struct {
	UID         UID
	QuestionUID UID
	Content     string
	IsCorrect   bool
}

func NewChoice(questionUID UID, content string, isCorrect bool) Choice {
	return Choice{
		UID:         GenerateUID(),
		QuestionUID: questionUID,
		Content:     content,
		IsCorrect:   isCorrect,
	}
}
