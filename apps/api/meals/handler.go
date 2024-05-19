package meals

import (
  "net/http"

  "github.com/labstack/echo/v4"
)

type Handler struct {
	MealRepository *MealRepository
}

func (h *Handler) GetMeals(c echo.Context) error {
	meals, err := h.MealRepository.Get()

	if err != nil {
		return err
	}

	return c.JSON(http.StatusOK, meals)
}

func (h *Handler) GetMeal(c echo.Context) error {
	meal, err := h.MealRepository.Find(c.Param("id"))

	if err != nil {
		return err
	}

	return c.JSON(http.StatusOK, meal)
}

func (h *Handler) AddMeal(c echo.Context) error {
	body := new(Meal)
	if err := c.Bind(body); err != nil {
		return err
	}

	m, err := NewMeal(body.Id, body.Name, body.MealIngredients)
	if err != nil {
		return err
	}

	c.Logger().Debugf("Adding meal: %v", m)

	if err := h.MealRepository.Save(m); err != nil {
		return err
	}

	return c.JSON(http.StatusAccepted, m)
}

func (h *Handler) AddIngredientToMeal(c echo.Context) error {
	mealId := c.Param("mealId")

	ingredient := new(MealIngredient)
	if err := c.Bind(ingredient); err != nil {
		return err
	}

	meal, err := h.MealRepository.Find(mealId)
	if err != nil {
		return err
	}

	meal.AddIngredient(*ingredient)

	if err := h.MealRepository.Save(meal); err != nil {
		return err
	}

	return c.JSON(http.StatusOK, meal)
}

func (h *Handler) RemoveIngredientFromMeal(c echo.Context) error {
	mealId := c.Param("mealId")
	ingredientId := c.Param("ingredientId")

	meal, err := h.MealRepository.Find(mealId)
	if err != nil {
		return err
	}

	meal.RemoveIngredient(ingredientId)

	if err := h.MealRepository.Save(meal); err != nil {
		return err
	}

	return c.JSON(http.StatusOK, meal)
}
