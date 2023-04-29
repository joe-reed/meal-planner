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

func (h *Handler) AddMeal(c echo.Context) error {
	m := new(Meal)
	if err := c.Bind(m); err != nil {
		return err
	}
	h.MealRepository.Add(m)
	return c.JSON(http.StatusAccepted, "Meal added")
}
