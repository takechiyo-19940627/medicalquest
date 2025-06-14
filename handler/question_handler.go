package handlers

import (
	"errors"
	"net/http"

	"github.com/labstack/echo/v4"

	"github.com/takechiyo-19940627/medicalquest/handler/request"
	"github.com/takechiyo-19940627/medicalquest/handler/response"
	"github.com/takechiyo-19940627/medicalquest/service"
	"github.com/takechiyo-19940627/medicalquest/service/dto"
	serviceErrors "github.com/takechiyo-19940627/medicalquest/service/errors"
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
		return h.handleServiceError(c, err)
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
		return h.handleServiceError(c, err)
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

func (h *QuestionHandler) Submit(c echo.Context) error {
	id := c.Param("id")
	if id == "" {
		return echo.NewHTTPError(http.StatusBadRequest, "id is required")
	}

	q := new(request.SubmitAnswerRequest)
	if err := c.Bind(q); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	if err := c.Validate(q); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	req := dto.AnswerRequest{
		QuestionID:       id,
		SelectedChoiceID: q.SelectedChoiceID,
	}
	result, err := h.service.Submit(c.Request().Context(), req)
	if err != nil {
		return h.handleServiceError(c, err)
	}

	response := response.NewAnswerResponse(result)
	return c.JSON(http.StatusOK, response)
}

// Update updates an existing question
func (h *QuestionHandler) Update(c echo.Context) error {
	return c.JSON(http.StatusOK, "")
}

// Delete deletes a question
func (h *QuestionHandler) Delete(c echo.Context) error {
	return c.NoContent(http.StatusNoContent)
}

// handleServiceError サービス層のエラーを適切なHTTPレスポンスに変換
func (h *QuestionHandler) handleServiceError(c echo.Context, err error) error {
	var svcErr *serviceErrors.ServiceError
	if !errors.As(err, &svcErr) {
		// ServiceError以外のエラーは内部エラーとして扱う
		c.Logger().Errorf("Unexpected error type: %v", err)
		return echo.NewHTTPError(http.StatusInternalServerError, "Internal server error")
	}

	// エラーの詳細をログに記録
	c.Logger().Errorf("Service error: type=%s, message=%s, cause=%v", 
		svcErr.Type, svcErr.Message, svcErr.Unwrap())

	// エラータイプに応じてHTTPステータスコードを決定
	switch svcErr.Type {
	case serviceErrors.TypeNotFound:
		return echo.NewHTTPError(http.StatusNotFound, svcErr.Message)
	
	case serviceErrors.TypeValidation:
		return c.JSON(http.StatusBadRequest, map[string]interface{}{
			"error": "validation_error",
			"field": svcErr.Field,
			"message": svcErr.Message,
		})
	
	case serviceErrors.TypeInternal:
		// 内部エラーの詳細はクライアントに露出しない
		return echo.NewHTTPError(http.StatusInternalServerError, "Internal server error")
	
	default:
		return echo.NewHTTPError(http.StatusInternalServerError, "Internal server error")
	}
}
