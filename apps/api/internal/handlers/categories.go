package handlers

import (
	"github.com/joe-reed/meal-planner/apps/api/internal/domain/categories"
	"net/http"

	"github.com/labstack/echo/v4"
)

type CategoriesHandler struct {
}

func (h *CategoriesHandler) GetCategories(c echo.Context) error {
	cats := categories.Categories()

	categoryObjects := make([]map[string]string, 0, len(cats))
	for _, category := range cats {
		categoryObjects = append(categoryObjects, map[string]string{"name": category})
	}

	return c.JSON(http.StatusOK, categoryObjects)
}
