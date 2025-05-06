package handlers

import (
    "net/http"
    "github.com/labstack/echo/v4"

    "github.com/takechiyo-19940627/medicalquest/internal/ent"
)

// ChoiceHandler handles HTTP requests related to choices
type ChoiceHandler struct {
    client *ent.Client
}

// NewChoiceHandler creates a new ChoiceHandler
func NewChoiceHandler(client *ent.Client) *ChoiceHandler {
    return &ChoiceHandler{
        client: client,
    }
}

// GetByQuestionID returns all choices for a specific question
func (h *ChoiceHandler) GetByQuestionID(c echo.Context) error {
    return c.JSON(http.StatusOK, "")
}

// Create creates a new choice for a question
func (h *ChoiceHandler) Create(c echo.Context) error {
    return c.JSON(http.StatusCreated, "")
}

// Update updates an existing choice
func (h *ChoiceHandler) Update(c echo.Context) error {
    return c.JSON(http.StatusOK, "")
}

// Delete deletes a choice
func (h *ChoiceHandler) Delete(c echo.Context) error {
    return c.NoContent(http.StatusNoContent)
}