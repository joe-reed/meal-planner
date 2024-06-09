package ingredients

import (
	"net/http"

	"github.com/labstack/echo/v4"
)

type Handler struct {
	IngredientRepository *IngredientRepository
}

func (h *Handler) GetIngredients(c echo.Context) error {
	ingredients, err := h.IngredientRepository.Get()

	if err != nil {
		return err
	}

	return c.JSON(http.StatusOK, ingredients)
}

func (h *Handler) AddIngredient(c echo.Context) error {
	body := new(Ingredient)
	if err := c.Bind(body); err != nil {
		return err
	}

	i, err := NewIngredient(body.Id, body.Name)
	if err != nil {
		return err
	}

	c.Logger().Debugf("Adding ingredient: %v", i)

	if err := h.IngredientRepository.Add(i); err != nil {
		return err
	}

	return c.JSON(http.StatusAccepted, i)
}
