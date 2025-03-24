package handlers_test

import (
	"github.com/joe-reed/meal-planner/apps/api/internal/domain/shops"
	"github.com/joe-reed/meal-planner/apps/api/internal/handlers"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
)

func TestRemovingMealFromCurrentShop(t *testing.T) {
	shop, err := shops.NewShop(1)
	assert.NoError(t, err)

	shop.AddMeal(&shops.ShopMeal{MealId: "abc"}).AddMeal(&shops.ShopMeal{MealId: "def"})

	r := shops.NewFakeShopRepository()
	err = r.Save(shop)
	assert.NoError(t, err)

	e := echo.New()
	req := httptest.NewRequest("DELETE", "/shops/current/meals/abc", nil)

	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.SetParamNames("mealId")
	c.SetParamValues("abc")
	h := &handlers.ShopsHandler{ShopRepository: r}

	if assert.NoError(t, h.RemoveMealFromCurrentShop(c)) {
		assert.Equal(t, http.StatusOK, rec.Code)
		s, _ := r.Find(1)
		assert.Equal(t, []*shops.ShopMeal{{MealId: "def"}}, s.Meals)
	}
}

func TestRemovingAllMealsFromCurrentShop(t *testing.T) {
	shop, err := shops.NewShop(1)
	assert.NoError(t, err)

	shop.AddMeal(&shops.ShopMeal{MealId: "abc"})

	r := shops.NewFakeShopRepository()
	err = r.Save(shop)
	assert.NoError(t, err)

	e := echo.New()
	req := httptest.NewRequest("DELETE", "/shops/current/meals/abc", nil)

	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.SetParamNames("mealId")
	c.SetParamValues("abc")
	h := &handlers.ShopsHandler{ShopRepository: r}

	if assert.NoError(t, h.RemoveMealFromCurrentShop(c)) {
		assert.Equal(t, http.StatusOK, rec.Code)
		s, _ := r.Find(1)
		assert.Equal(t, []*shops.ShopMeal{}, s.Meals)
	}
}
