package handlers

import (
    "net/http"
    "strconv"
    
    "github.com/labstack/echo/v4"
    
    "github.com/medicalquest/internal/ent"
    "github.com/medicalquest/internal/ent/choice"
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
    questionID, err := strconv.Atoi(c.Param("questionID"))
    if err != nil {
        return echo.NewHTTPError(http.StatusBadRequest, "Invalid question ID")
    }
    
    choices, err := h.client.Choice.Query().
        Where(choice.HasQuestionWith(choice.IDEQ(questionID))).
        All(c.Request().Context())
    
    if err != nil {
        return echo.NewHTTPError(http.StatusInternalServerError, "Failed to fetch choices")
    }
    
    return c.JSON(http.StatusOK, choices)
}

// Create creates a new choice for a question
func (h *ChoiceHandler) Create(c echo.Context) error {
    questionID, err := strconv.Atoi(c.Param("questionID"))
    if err != nil {
        return echo.NewHTTPError(http.StatusBadRequest, "Invalid question ID")
    }
    
    type request struct {
        Content   string `json:"content"`
        IsCorrect bool   `json:"is_correct"`
    }
    
    var req request
    if err := c.Bind(&req); err != nil {
        return echo.NewHTTPError(http.StatusBadRequest, "Invalid request body")
    }
    
    choice, err := h.client.Choice.Create().
        SetQuestionID(questionID).
        SetContent(req.Content).
        SetIsCorrect(req.IsCorrect).
        Save(c.Request().Context())
    
    if err != nil {
        return echo.NewHTTPError(http.StatusInternalServerError, "Failed to create choice")
    }
    
    return c.JSON(http.StatusCreated, choice)
}

// Update updates an existing choice
func (h *ChoiceHandler) Update(c echo.Context) error {
    id, err := strconv.Atoi(c.Param("id"))
    if err != nil {
        return echo.NewHTTPError(http.StatusBadRequest, "Invalid choice ID")
    }
    
    type request struct {
        Content   string `json:"content"`
        IsCorrect bool   `json:"is_correct"`
    }
    
    var req request
    if err := c.Bind(&req); err != nil {
        return echo.NewHTTPError(http.StatusBadRequest, "Invalid request body")
    }
    
    choice, err := h.client.Choice.UpdateOneID(id).
        SetContent(req.Content).
        SetIsCorrect(req.IsCorrect).
        Save(c.Request().Context())
    
    if err != nil {
        return echo.NewHTTPError(http.StatusInternalServerError, "Failed to update choice")
    }
    
    return c.JSON(http.StatusOK, choice)
}

// Delete deletes a choice
func (h *ChoiceHandler) Delete(c echo.Context) error {
    id, err := strconv.Atoi(c.Param("id"))
    if err != nil {
        return echo.NewHTTPError(http.StatusBadRequest, "Invalid choice ID")
    }
    
    err = h.client.Choice.DeleteOneID(id).Exec(c.Request().Context())
    if err != nil {
        return echo.NewHTTPError(http.StatusInternalServerError, "Failed to delete choice")
    }
    
    return c.NoContent(http.StatusNoContent)
}