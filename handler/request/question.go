package request

type CreateQuestionRequest struct {
	Title         string `json:"title" validate:"required,min=1,max=200"`
	Content       string `json:"content" validate:"required,min=1,max=800"`
	ReferenceCode string `json:"reference_code" validate:"required,min=1,max=10"`
}
