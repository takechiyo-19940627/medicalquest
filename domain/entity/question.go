package entity

type Question struct {
	UID           UID
	ReferenceCode string
	Title         string
	Content       string
}

func NewQuestion(ReferenceCode, title, content string) Question {
	return Question{
		UID:           NewUID(),
		ReferenceCode: ReferenceCode,
		Title:         title,
		Content:       content,
	}
}
