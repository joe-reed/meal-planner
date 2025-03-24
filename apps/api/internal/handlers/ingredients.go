package handlers

import (
	"github.com/joe-reed/meal-planner/apps/api/internal/domain/ingredients"
	"net/http"

	"github.com/labstack/echo/v4"
)

type IngredientsHandler struct {
	IngredientRepository *ingredients.IngredientRepository
}

func (h *IngredientsHandler) GetIngredients(c echo.Context) error {
	ings, err := h.IngredientRepository.Get()

	if err != nil {
		return err
	}

	return c.JSON(http.StatusOK, ings)
}

func (h *IngredientsHandler) AddIngredient(c echo.Context) error {
	body := new(ingredients.Ingredient)
	if err := c.Bind(body); err != nil {
		return err
	}

	i, err := ingredients.NewIngredient(body.Id, body.Name, body.Category)
	if err != nil {
		return err
	}

	c.Logger().Debugf("Adding ingredient: %v", i)

	if err := h.IngredientRepository.Add(i); err != nil {
		return err
	}

	return c.JSON(http.StatusAccepted, i)
}
