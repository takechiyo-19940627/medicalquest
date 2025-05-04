package handlers

import (
    "github.com/labstack/echo/v4"
    
    "github.com/medicalquest/internal/ent"
)

// RegisterRoutes sets up all routes for the application
func RegisterRoutes(e *echo.Echo, client *ent.Client) {
    // Health check
    e.GET("/health", func(c echo.Context) error {
        return c.JSON(200, map[string]string{"status": "ok"})
    })
    
    // API routes group
    api := e.Group("/api")
    
    // Questions routes
    qh := NewQuestionHandler(client)
    api.GET("/questions", qh.GetAll)
    api.GET("/questions/:id", qh.GetByID)
    api.POST("/questions", qh.Create)
    api.PUT("/questions/:id", qh.Update)
    api.DELETE("/questions/:id", qh.Delete)
    
    // Choices routes
    ch := NewChoiceHandler(client)
    api.GET("/questions/:questionID/choices", ch.GetByQuestionID)
    api.POST("/questions/:questionID/choices", ch.Create)
    api.PUT("/choices/:id", ch.Update)
    api.DELETE("/choices/:id", ch.Delete)
}