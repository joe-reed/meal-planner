package shops_test

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/joe-reed/meal-planner/apps/api/shops"
	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
)

func TestRemovingMealFromCurrentShop(t *testing.T) {
	shop := shops.NewShop(1).AddMeal(&shops.ShopMeal{MealId: "abc"}).AddMeal(&shops.ShopMeal{MealId: "def"})

	r := shops.NewFakeShopRepository()
	err := r.Add(shop)
	assert.NoError(t, err)

	e := echo.New()
	req := httptest.NewRequest("DELETE", "/shops/current/meals/abc", nil)

	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.SetParamNames("mealId")
	c.SetParamValues("abc")
	h := &shops.Handler{ShopRepository: r}

	if assert.NoError(t, h.RemoveMealFromCurrentShop(c)) {
		assert.Equal(t, http.StatusOK, rec.Code)
		s, _ := r.Find(1)
		assert.Equal(t, []*shops.ShopMeal{{MealId: "def"}}, s.Meals)
	}
}
