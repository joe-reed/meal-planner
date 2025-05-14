package handlers

import (
	"github.com/joe-reed/meal-planner/apps/api/internal/application"
	"github.com/joe-reed/meal-planner/apps/api/internal/domain/shop"
	"github.com/labstack/echo/v4"
	"net/http"
)

type ShopsHandler struct {
	Application *application.ShopApplication
}

func (h *ShopsHandler) CurrentShop(c echo.Context) error {
	s, err := h.Application.GetCurrentShop()

	if err != nil {
		return err
	}

	return c.JSON(http.StatusOK, s)
}

func (h *ShopsHandler) StartShop(c echo.Context) error {
	s, err := h.Application.StartShop()

	if err != nil {
		return err
	}

	return c.JSON(http.StatusOK, s)
}

func (h *ShopsHandler) AddMealToCurrentShop(c echo.Context) error {
	shopMeal := new(shop.ShopMeal)
	if err := c.Bind(shopMeal); err != nil {
		return err
	}

	s, err := h.Application.AddMealToCurrentShop(shopMeal)

	if err != nil {
		return err
	}

	return c.JSON(http.StatusOK, s)
}

func (h *ShopsHandler) RemoveMealFromCurrentShop(c echo.Context) error {
	s, err := h.Application.RemoveMealFromCurrentShop(c.Param("mealId"))

	if err != nil {
		return err
	}

	return c.JSON(http.StatusOK, s)
}

func (h *ShopsHandler) AddItemToCurrentShop(c echo.Context) error {
	item := new(shop.Item)
	if err := c.Bind(item); err != nil {
		return err
	}

	s, err := h.Application.AddItemToCurrentShop(item)

	if err != nil {
		return err
	}

	return c.JSON(http.StatusOK, s)
}
