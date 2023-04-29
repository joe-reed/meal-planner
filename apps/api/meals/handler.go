package meals

import (
	"net/http"

	"github.com/labstack/echo/v4"
)

type Handler struct {
	MealRepository MealRepository
}

func (h *Handler) GetMeals(c echo.Context) error {
	meals, err := h.MealRepository.Get()

	if err != nil {
		return err
	}

	return c.JSON(http.StatusOK, meals)
}

func (h *Handler) AddMeal(c echo.Context) error {
	m := new(Meal)
	if err := c.Bind(m); err != nil {
		return err
	}
	c.Logger().Debugf("Adding meal: %v", m)

	if err := h.MealRepository.Add(m); err != nil {
		return err
	}

	return c.JSON(http.StatusAccepted, "Meal added")
}
