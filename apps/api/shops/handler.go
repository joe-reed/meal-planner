package shops

import (
	"net/http"
	"strconv"

	"github.com/labstack/echo/v4"
)

type Handler struct {
	ShopRepository ShopRepository
}

func (h *Handler) CurrentShop(c echo.Context) error {
	shop, err := h.ShopRepository.Current()

	if err != nil {
		return err
	}

	return c.JSON(http.StatusOK, shop)
}

func (h *Handler) StartShop(c echo.Context) error {
	shop, err := h.ShopRepository.Current()
	if err != nil {
		return err
	}

	if shop == nil {
		err = h.ShopRepository.Add(&Shop{Id: 1})
	} else {
		err = h.ShopRepository.Add(&Shop{Id: shop.Id + 1})
	}

	if err != nil {
		return err
	}

	return c.JSON(http.StatusOK, shop)
}

func (h *Handler) AddMealToShop(c echo.Context) error {
	i, err := strconv.Atoi(c.Param("shopId"))
	if err != nil {
		return err
	}

	shop, err := h.ShopRepository.Find(i)

	if err != nil {
		return err
	}

	shopMeal := new(ShopMeal)
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
