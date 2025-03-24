package handlers

import (
	"github.com/joe-reed/meal-planner/apps/api/internal/domain/basket"
	"github.com/labstack/echo/v4"
	"net/http"
	"strconv"
)

type BasketHandler struct {
	BasketRepository *basket.BasketRepository
}

func (h *BasketHandler) AddItemToBasket(c echo.Context) error {
	b, err := getBasketFromContext(c, h)
	if err != nil {
		return err
	}

	i := new(basket.BasketItem)
	if err := c.Bind(i); err != nil {
		return err
	}

	c.Logger().Debugf("Adding item to basket: shopId: %d ingredientId: %s", b.ShopId, i.IngredientId)

	b.AddItem(i)

	if err := h.BasketRepository.Save(b); err != nil {
		return err
	}

	return c.JSON(http.StatusOK, b)
}

func (h *BasketHandler) RemoveItemFromBasket(c echo.Context) error {
	b, err := getBasketFromContext(c, h)
	if err != nil {
		return err
	}

	ingredientId := c.Param("ingredientId")

	c.Logger().Debugf("Removing item from basket: shopId: %d ingredientId: %s", b.ShopId, ingredientId)

	b.RemoveItem(ingredientId)

	if err := h.BasketRepository.Save(b); err != nil {
		return err
	}

	return c.JSON(http.StatusOK, b)
}

func (h *BasketHandler) GetBasketItems(c echo.Context) error {
	b, err := getBasketFromContext(c, h)
	if err != nil {
		return err
	}

	return c.JSON(http.StatusOK, b)
}

func getBasketFromContext(c echo.Context, h *BasketHandler) (*basket.Basket, error) {
	shopId, err := strconv.Atoi(c.Param("shopId"))

	if err != nil {
		return nil, err
	}

	b, err := h.BasketRepository.FindByShopId(shopId)

	if err != nil {
		return nil, err
	}

	return b, nil
}
