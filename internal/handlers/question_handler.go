package handlers

import (
    "net/http"
    "strconv"
    
    "github.com/labstack/echo/v4"
    
    "github.com/takechiyo-19940627/medicalquest/internal/ent"
    "github.com/takechiyo-19940627/medicalquest/internal/ent/question"
)

// QuestionHandler handles HTTP requests related to questions
type QuestionHandler struct {
    client *ent.Client
}

// NewQuestionHandler creates a new QuestionHandler
func NewQuestionHandler(client *ent.Client) *QuestionHandler {
    return &QuestionHandler{
        client: client,
    }
}

// GetAll returns all questions
func (h *QuestionHandler) GetAll(c echo.Context) error {
    questions, err := h.client.Question.Query().All(c.Request().Context())
    if err != nil {
        return echo.NewHTTPError(http.StatusInternalServerError, "Failed to fetch questions")
    }
    
    return c.JSON(http.StatusOK, questions)
}

// GetByID returns a question by ID
func (h *QuestionHandler) GetByID(c echo.Context) error {
    id, err := strconv.Atoi(c.Param("id"))
    if err != nil {
        return echo.NewHTTPError(http.StatusBadRequest, "Invalid question ID")
    }
    
    q, err := h.client.Question.Query().Where(question.ID(id)).Only(c.Request().Context())
    if err != nil {
        return echo.NewHTTPError(http.StatusNotFound, "Question not found")
    }
    
    return c.JSON(http.StatusOK, q)
}

// Create creates a new question
func (h *QuestionHandler) Create(c echo.Context) error {
    type request struct {
        ReferenceCode string `json:"reference_code"`
        Title         string `json:"title"`
        Content       string `json:"content"`
    }
    
    var req request
    if err := c.Bind(&req); err != nil {
        return echo.NewHTTPError(http.StatusBadRequest, "Invalid request body")
    }
    
    q, err := h.client.Question.Create().
        SetReferenceCode(req.ReferenceCode).
        SetTitle(req.Title).
        SetContent(req.Content).
        Save(c.Request().Context())
    
    if err != nil {
        return echo.NewHTTPError(http.StatusInternalServerError, "Failed to create question")
    }
    
    return c.JSON(http.StatusCreated, q)
}

// Update updates an existing question
func (h *QuestionHandler) Update(c echo.Context) error {
    id, err := strconv.Atoi(c.Param("id"))
    if err != nil {
        return echo.NewHTTPError(http.StatusBadRequest, "Invalid question ID")
    }
    
    type request struct {
        ReferenceCode string `json:"reference_code"`
        Title         string `json:"title"`
        Content       string `json:"content"`
    }
    
    var req request
    if err := c.Bind(&req); err != nil {
        return echo.NewHTTPError(http.StatusBadRequest, "Invalid request body")
    }
    
    q, err := h.client.Question.UpdateOneID(id).
        SetReferenceCode(req.ReferenceCode).
        SetTitle(req.Title).
        SetContent(req.Content).
        Save(c.Request().Context())
    
    if err != nil {
        return echo.NewHTTPError(http.StatusInternalServerError, "Failed to update question")
    }
    
    return c.JSON(http.StatusOK, q)
}

// Delete deletes a question
func (h *QuestionHandler) Delete(c echo.Context) error {
    id, err := strconv.Atoi(c.Param("id"))
    if err != nil {
        return echo.NewHTTPError(http.StatusBadRequest, "Invalid question ID")
    }
    
    err = h.client.Question.DeleteOneID(id).Exec(c.Request().Context())
    if err != nil {
        return echo.NewHTTPError(http.StatusInternalServerError, "Failed to delete question")
    }
    
    return c.NoContent(http.StatusNoContent)
}