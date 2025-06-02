package request

type SubmitAnswerRequest struct {
	SelectedChoiceID string `json:"selected_choice_id" validate:"required"`
}
