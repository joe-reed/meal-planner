package ingredients

import (
	"net/http"

	"github.com/labstack/echo/v4"
)

type Handler struct {
	IngredientRepository IngredientRepository
}

func (h *Handler) GetIngredients(c echo.Context) error {
	ingredients, err := h.IngredientRepository.Get()

	if err != nil {
		return err
	}

	return c.JSON(http.StatusOK, ingredients)
}
