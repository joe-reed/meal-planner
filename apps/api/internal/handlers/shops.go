package handlers

import (
	"github.com/joe-reed/meal-planner/apps/api/internal/application"
	"github.com/joe-reed/meal-planner/apps/api/internal/domain/shops"
	"github.com/labstack/echo/v4"
	"net/http"
)

type ShopsHandler struct {
	Application *application.ShopApplication
}

func (h *ShopsHandler) CurrentShop(c echo.Context) error {
	shop, err := h.Application.GetCurrentShop()

	if err != nil {
		return err
	}

	return c.JSON(http.StatusOK, shop)
}

func (h *ShopsHandler) StartShop(c echo.Context) error {
	shop, err := h.Application.StartShop()

	if err != nil {
		return err
	}

	return c.JSON(http.StatusOK, shop)
}

func (h *ShopsHandler) AddMealToCurrentShop(c echo.Context) error {
	shopMeal := new(shops.ShopMeal)
	if err := c.Bind(shopMeal); err != nil {
		return err
	}

	shop, err := h.Application.AddMealToCurrentShop(shopMeal)

	if err != nil {
		return err
	}

	return c.JSON(http.StatusOK, shop)
}

func (h *ShopsHandler) RemoveMealFromCurrentShop(c echo.Context) error {
	shop, err := h.Application.RemoveMealFromCurrentShop(c.Param("mealId"))

	if err != nil {
		return err
	}

	return c.JSON(http.StatusOK, shop)
}
