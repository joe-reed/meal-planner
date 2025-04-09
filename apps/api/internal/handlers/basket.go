package handlers

import (
	"github.com/joe-reed/meal-planner/apps/api/internal/application"
	"github.com/joe-reed/meal-planner/apps/api/internal/domain/basket"
	"github.com/labstack/echo/v4"
	"net/http"
	"strconv"
)

type BasketHandler struct {
	Application *application.BasketApplication
}

func (h *BasketHandler) AddItemToBasket(c echo.Context) error {
	shopId, err := getShopIdFromContext(c)

	if err != nil {
		return err
	}

	i := new(basket.BasketItem)
	if err := c.Bind(i); err != nil {
		return err
	}

	b, err := h.Application.AddItemToBasket(shopId, i)

	if err != nil {
		return err
	}

	return c.JSON(http.StatusOK, b)
}

func (h *BasketHandler) RemoveItemFromBasket(c echo.Context) error {
	shopId, err := getShopIdFromContext(c)

	if err != nil {
		return err
	}

	ingredientId := c.Param("ingredientId")

	b, err := h.Application.RemoveItemFromBasket(shopId, ingredientId)

	if err != nil {
		return err
	}

	return c.JSON(http.StatusOK, b)
}

func (h *BasketHandler) GetBasket(c echo.Context) error {
	shopId, err := getShopIdFromContext(c)

	if err != nil {
		return err
	}

	b, err := h.Application.GetBasket(shopId)

	if err != nil {
		return err
	}

	return c.JSON(http.StatusOK, b)
}

func getShopIdFromContext(c echo.Context) (int, error) {
	return strconv.Atoi(c.Param("shopId"))
}
