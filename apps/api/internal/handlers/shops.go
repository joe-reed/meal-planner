package handlers

import (
	"fmt"
	"github.com/joe-reed/meal-planner/apps/api/internal/domain/shops"
	"github.com/labstack/echo/v4"
	"net/http"
)

type ShopsHandler struct {
	ShopRepository *shops.ShopRepository
	Publisher      func(string)
}

func (h *ShopsHandler) CurrentShop(c echo.Context) error {
	shop, err := h.ShopRepository.Current()

	if err != nil {
		return err
	}

	return c.JSON(http.StatusOK, shop)
}

func (h *ShopsHandler) StartShop(c echo.Context) error {
	shop, err := h.ShopRepository.Current()
	if err != nil {
		return err
	}

	var newShop *shops.Shop

	if shop == nil {
		newShop, err = shops.NewShop(1)
	} else {
		newShop, err = shops.NewShop(shop.Id + 1)
	}

	if err != nil {
		return err
	}

	err = h.ShopRepository.Save(newShop)

	h.Publisher(fmt.Sprintf("shopStarted:%d", newShop.Id))

	return c.JSON(http.StatusOK, shop)
}

func (h *ShopsHandler) AddMealToCurrentShop(c echo.Context) error {
	shop, err := h.ShopRepository.Current()

	if err != nil {
		return err
	}

	shopMeal := new(shops.ShopMeal)
	if err := c.Bind(shopMeal); err != nil {
		return err
	}
	c.Logger().Debugf("Adding meal to shop: shopId: %d mealId: %s", shop.Id, shopMeal.MealId)

	shop.AddMeal(shopMeal)

	if err := h.ShopRepository.Save(shop); err != nil {
		return err
	}

	return c.JSON(http.StatusOK, shop)
}

func (h *ShopsHandler) RemoveMealFromCurrentShop(c echo.Context) error {
	shop, err := h.ShopRepository.Current()

	if err != nil {
		return err
	}

	mealId := c.Param("mealId")
	c.Logger().Debugf("Removing meal from shop: shopId: %d mealId: %s", shop.Id, mealId)

	shop.RemoveMeal(mealId)

	if err := h.ShopRepository.Save(shop); err != nil {
		return err
	}

	return c.JSON(http.StatusOK, shop)
}
