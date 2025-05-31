package entity

type Question struct {
	UID           UID
	ReferenceCode string
	Title         string
	Content       string
	Choices       []Choice
}

func NewQuestion(referenceCode, title, content string) Question {
	return Question{
		UID:           GenerateUID(),
		ReferenceCode: referenceCode,
		Title:         title,
		Content:       content,
	}
}
