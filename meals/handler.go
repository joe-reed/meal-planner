package meals

import (
	"net/http"

	"github.com/labstack/echo/v4"
)

type (
	Handler struct {
		MealRepository MealRepository
	}
)

func (h *Handler) GetMeals(c echo.Context) error {
	return c.JSON(http.StatusOK, h.MealRepository.Get())
}
