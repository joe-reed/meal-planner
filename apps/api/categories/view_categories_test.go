package categories_test

import (
	"github.com/joe-reed/meal-planner/apps/api/categories"
	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestViewingCategories(t *testing.T) {
	e := echo.New()
	req := httptest.NewRequest("", "/", nil)
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	h := &categories.Handler{}

	if assert.NoError(t, h.GetCategories(c)) {
		assert.Equal(t, http.StatusOK, rec.Code)
		assert.Equal(t, `[{"name":"AlcoholicDrinks"},{"name":"Bakery"},{"name":"Chilled"},{"name":"ChocolateAndSweets"},{"name":"Dairy"},{"name":"Desserts"},{"name":"Drinks"},{"name":"Eggs"},{"name":"FishAndSeafood"},{"name":"FoodCupboard"},{"name":"Frozen"},{"name":"Fruit"},{"name":"Meat"},{"name":"PastaRiceAndNoodles"},{"name":"SaucesOilsAndDressings"},{"name":"SeedsNutsAndDriedFruits"},{"name":"TeaAndCoffee"},{"name":"TinsCansAndPackets"},{"name":"Vegetables"}]`+"\n", rec.Body.String())
	}
}
