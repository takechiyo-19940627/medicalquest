package result

type QuestionResult struct {
	UID           string
	ReferenceCode string
	Title         string
	Content       string
}

type QuestionWithChoicesResult struct {
	QuestionResult
	Choices []ChoiceResult
}

type ChoiceResult struct {
	UID       string
	Content   string
	IsCorrect bool
}
