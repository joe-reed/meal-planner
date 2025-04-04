package handlers

import (
	"github.com/joe-reed/meal-planner/apps/api/internal/application"
	"github.com/joe-reed/meal-planner/apps/api/internal/domain/ingredients"
	"net/http"

	"github.com/labstack/echo/v4"
)

type IngredientsHandler struct {
	Application *application.IngredientApplication
}

func (h *IngredientsHandler) GetIngredients(c echo.Context) error {
	ings, err := h.Application.ListIngredients()

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

	i, err := h.Application.AddIngredient(body.Id, body.Name, body.Category)

	if err != nil {
		return err
	}

	return c.JSON(http.StatusAccepted, i)
}
