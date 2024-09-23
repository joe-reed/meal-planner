package categories

import (
	"net/http"

	"github.com/labstack/echo/v4"
)

type Handler struct {
}

func (h *Handler) GetCategories(c echo.Context) error {
	categories := Categories()

	categoryObjects := make([]map[string]string, 0, len(categories))
	for _, category := range categories {
		categoryObjects = append(categoryObjects, map[string]string{"name": category})
	}

	return c.JSON(http.StatusOK, categoryObjects)
}
