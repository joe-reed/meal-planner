package shops_test

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/joe-reed/meal-planner/apps/api/shops"
	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
)

func TestAddingMealToShop(t *testing.T) {
	shop := shops.NewShop(1)

	r := shops.NewFakeShopRepository()
	r.Add(shop)

	e := echo.New()
	req := httptest.NewRequest("POST", "/shops/1/meals", strings.NewReader(`{"id":"abc"}`))

	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.SetParamNames("shopId")
	c.SetParamValues("1")
	h := &shops.Handler{ShopRepository: r}

	if assert.NoError(t, h.AddMealToShop(c)) {
		assert.Equal(t, http.StatusOK, rec.Code)
		s, _ := r.Find(1)
		assert.Equal(t, []*shops.ShopMeal{{MealId: "abc"}}, s.Meals)
	}
}
