package handlers

import (
	"net/http"

	"github.com/labstack/echo/v4"

	"github.com/takechiyo-19940627/medicalquest/handler/request"
	"github.com/takechiyo-19940627/medicalquest/handler/response"
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
	res, err := h.service.FindAll(c.Request().Context())
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}
	response := response.NewQuestionResponse(res)

	return c.JSON(http.StatusOK, response)
}

// GetByID returns a question by ID
func (h *QuestionHandler) GetByID(c echo.Context) error {
	id := c.Param("id")
	if id == "" {
		return echo.NewHTTPError(http.StatusBadRequest, "id is required")
	}

	result, err := h.service.FindByID(c.Request().Context(), id)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}
	response := response.NewQuestionWithChoicesResponse(result)
	return c.JSON(http.StatusOK, response)
}

// Create creates a new question
func (h *QuestionHandler) Create(c echo.Context) error {
	q := new(request.CreateQuestionRequest)
	if err := c.Bind(q); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	if err := c.Validate(q); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}
	return c.JSON(http.StatusCreated, "Created")
}

// Update updates an existing question
func (h *QuestionHandler) Update(c echo.Context) error {
	return c.JSON(http.StatusOK, "")
}

// Delete deletes a question
func (h *QuestionHandler) Delete(c echo.Context) error {
	return c.NoContent(http.StatusNoContent)
}
