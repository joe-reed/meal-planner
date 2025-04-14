package handlers_test

import (
	"github.com/joe-reed/meal-planner/apps/api/internal/application"
	"github.com/joe-reed/meal-planner/apps/api/internal/domain/shop"
	"github.com/joe-reed/meal-planner/apps/api/internal/handlers"
	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestRemovingMealFromCurrentShop(t *testing.T) {
	s, err := shop.NewShop(1)
	assert.NoError(t, err)

	s.AddMeal(&shop.ShopMeal{MealId: "abc"}).AddMeal(&shop.ShopMeal{MealId: "def"})

	r := shop.NewFakeShopRepository()
	err = r.Save(s)
	assert.NoError(t, err)

	e := echo.New()
	req := httptest.NewRequest("DELETE", "/shops/current/meals/abc", nil)

	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.SetParamNames("mealId")
	c.SetParamValues("abc")
	h := &handlers.ShopsHandler{Application: application.NewShopApplication(r, func(string) {})}

	if assert.NoError(t, h.RemoveMealFromCurrentShop(c)) {
		assert.Equal(t, http.StatusOK, rec.Code)
		s, _ := r.Find(1)
		assert.Equal(t, []*shop.ShopMeal{{MealId: "def"}}, s.Meals)
	}
}

func TestRemovingAllMealsFromCurrentShop(t *testing.T) {
	s, err := shop.NewShop(1)
	assert.NoError(t, err)

	s.AddMeal(&shop.ShopMeal{MealId: "abc"})

	r := shop.NewFakeShopRepository()
	err = r.Save(s)
	assert.NoError(t, err)

	e := echo.New()
	req := httptest.NewRequest("DELETE", "/shops/current/meals/abc", nil)

	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.SetParamNames("mealId")
	c.SetParamValues("abc")
	h := &handlers.ShopsHandler{Application: application.NewShopApplication(r, func(string) {})}

	if assert.NoError(t, h.RemoveMealFromCurrentShop(c)) {
		assert.Equal(t, http.StatusOK, rec.Code)
		s, _ := r.Find(1)
		assert.Equal(t, []*shop.ShopMeal{}, s.Meals)
	}
}
