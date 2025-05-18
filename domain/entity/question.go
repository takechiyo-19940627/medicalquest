package entity

type Question struct {
	UID           UID
	ReferenceCode string
	Title         string
	Content       string
}

func NewQuestion(referenceCode, title, content string) Question {
	return Question{
		UID:           GenerateUID(),
		ReferenceCode: referenceCode,
		Title:         title,
		Content:       content,
	}
}

func NewQuestionFromPersistence(uid, referenceCode, title, content string) Question {
	return Question{
		UID:           ToUID(uid),
		ReferenceCode: referenceCode,
		Title:         title,
		Content:       content,
	}
}
