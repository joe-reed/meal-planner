package handlers_test

import (
	"github.com/joe-reed/meal-planner/apps/api/internal/application"
	"github.com/joe-reed/meal-planner/apps/api/internal/domain/shops"
	"github.com/joe-reed/meal-planner/apps/api/internal/handlers"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
)

func TestAddingMealToCurrentShop(t *testing.T) {
	shop, err := shops.NewShop(1)
	assert.NoError(t, err)

	r := shops.NewFakeShopRepository()
	err = r.Save(shop)
	assert.NoError(t, err)

	e := echo.New()
	req := httptest.NewRequest("POST", "/shops/current/meals", strings.NewReader(`{"id":"abc"}`))

	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	h := &handlers.ShopsHandler{Application: application.NewShopApplication(r, func(string) {})}

	if assert.NoError(t, h.AddMealToCurrentShop(c)) {
		assert.Equal(t, http.StatusOK, rec.Code)
		s, _ := r.Find(1)
		assert.Equal(t, []*shops.ShopMeal{{MealId: "abc"}}, s.Meals)
	}
}
