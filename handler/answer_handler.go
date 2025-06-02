package handlers

import (
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/takechiyo-19940627/medicalquest/handler/request"
	"github.com/takechiyo-19940627/medicalquest/handler/response"
	"github.com/takechiyo-19940627/medicalquest/service"
	"github.com/takechiyo-19940627/medicalquest/service/dto"
)

type AnswerHandler struct {
	service *service.AnswerService
}

func NewAnswerHandler(service *service.AnswerService) *AnswerHandler {
	return &AnswerHandler{
		service,
	}
}

func (h *AnswerHandler) Submit(c echo.Context) error {
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
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	response := response.NewAnswerResponse(result)
	return c.JSON(http.StatusOK, response)
}
