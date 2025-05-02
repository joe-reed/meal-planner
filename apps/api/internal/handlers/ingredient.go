package handlers

import (
	"errors"
	"github.com/joe-reed/meal-planner/apps/api/internal/application"
	"github.com/joe-reed/meal-planner/apps/api/internal/domain/ingredient"
	"github.com/labstack/echo/v4"
	"net/http"
)

type IngredientsHandler struct {
	Application *application.IngredientApplication
}

func (h *IngredientsHandler) GetIngredients(c echo.Context) error {
	ings, err := h.Application.GetIngredients()

	if err != nil {
		return err
	}

	return c.JSON(http.StatusOK, ings)
}

func (h *IngredientsHandler) AddIngredient(c echo.Context) error {
	body := new(ingredient.Ingredient)

	if err := c.Bind(body); err != nil {
		return c.JSON(
			http.StatusBadRequest,
			struct {
				Error string `json:"error"`
			}{
				Error: "Invalid request body",
			})
	}

	i, err := h.Application.AddIngredient(body.Id, body.Name, body.Category)

	if err != nil {
		var validationError *application.ValidationError
		if errors.As(err, &validationError) {
			return c.JSON(http.StatusBadRequest, struct {
				Error string `json:"error"`
			}{
				Error: validationError.Error(),
			})
		}
		return err
	}

	return c.JSON(http.StatusAccepted, i)
}
