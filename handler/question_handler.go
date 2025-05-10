package handlers

import (
	"net/http"

	"github.com/labstack/echo/v4"

	"github.com/takechiyo-19940627/medicalquest/service"
)

// QuestionHandler handles HTTP requests related to questions
type QuestionHandler struct {
	service *service.QuestionService
}

// NewQuestionHandler creates a new QuestionHandler
func NewQuestionHandler(service *service.QuestionService) *QuestionHandler {
	return &QuestionHandler{
		service,
	}
}

// GetAll returns all questions
func (h *QuestionHandler) GetAll(c echo.Context) error {
	return c.JSON(http.StatusOK, "")
}

// GetByID returns a question by ID
func (h *QuestionHandler) GetByID(c echo.Context) error {
	return c.JSON(http.StatusOK, "")
}

// Create creates a new question
func (h *QuestionHandler) Create(c echo.Context) error {
	return c.JSON(http.StatusCreated, "")
}

// Update updates an existing question
func (h *QuestionHandler) Update(c echo.Context) error {
	return c.JSON(http.StatusOK, "")
}

// Delete deletes a question
func (h *QuestionHandler) Delete(c echo.Context) error {
	return c.NoContent(http.StatusNoContent)
}
