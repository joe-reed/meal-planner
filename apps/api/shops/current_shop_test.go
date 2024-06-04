package shops_test

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/joe-reed/meal-planner/apps/api/shops"
	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
)

func TestGettingCurrentShop(t *testing.T) {
	shop1, err := shops.NewShop(1)
	assert.NoError(t, err)
	shop2, err := shops.NewShop(2)
	assert.NoError(t, err)

	shop2.AddMeal(&shops.ShopMeal{MealId: "123"})
	shop2.AddMeal(&shops.ShopMeal{MealId: "456"})

	r := shops.NewFakeShopRepository()
	err = r.Save(shop1)
	assert.NoError(t, err)
	err = r.Save(shop2)
	assert.NoError(t, err)

	e := echo.New()
	req := httptest.NewRequest("GET", "/shops/current", nil)
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	h := &shops.Handler{ShopRepository: r}

	if assert.NoError(t, h.CurrentShop(c)) {
		assert.Equal(t, http.StatusOK, rec.Code)
		assert.Equal(t, `{"id":2,"meals":[{"id":"123"},{"id":"456"}]}`+"\n", rec.Body.String())
	}
}

func TestGettingCurrentShopWithNoMeals(t *testing.T) {
	shop1, err := shops.NewShop(1)
	assert.NoError(t, err)

	r := shops.NewFakeShopRepository()
	err = r.Save(shop1)
	assert.NoError(t, err)

	e := echo.New()
	req := httptest.NewRequest("GET", "/shops/current", nil)
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	h := &shops.Handler{ShopRepository: r}

	if assert.NoError(t, h.CurrentShop(c)) {
		assert.Equal(t, http.StatusOK, rec.Code)
		assert.Equal(t, `{"id":1,"meals":[]}`+"\n", rec.Body.String())
	}
}
