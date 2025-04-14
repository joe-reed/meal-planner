package handlers

import (
	"github.com/joe-reed/meal-planner/apps/api/internal/application"
	"net/http"

	"github.com/labstack/echo/v4"
)

type CategoriesHandler struct {
	Application *application.CategoryApplication
}

func (h *CategoriesHandler) GetCategories(c echo.Context) error {
	return c.JSON(http.StatusOK, h.Application.GetCategories())
}
