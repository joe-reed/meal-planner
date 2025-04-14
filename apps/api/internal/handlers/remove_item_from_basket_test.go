package handlers_test

import (
	"github.com/joe-reed/meal-planner/apps/api/internal/application"
	"github.com/joe-reed/meal-planner/apps/api/internal/domain/basket"
	"github.com/joe-reed/meal-planner/apps/api/internal/handlers"
	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestRemovingItemFromBasket(t *testing.T) {
	b, err := basket.NewBasket(1)
	assert.NoError(t, err)

	b.AddItem(&basket.BasketItem{IngredientId: "ing-1"})
	b.AddItem(&basket.BasketItem{IngredientId: "ing-2"})

	br := basket.NewFakeBasketRepository()

	err = br.Save(b)
	assert.NoError(t, err)

	e := echo.New()
	req := httptest.NewRequest("DELETE", "/baskets/1/items/ing-1", nil)

	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.SetParamNames("shopId", "ingredientId")
	c.SetParamValues("1", "ing-1")
	h := &handlers.BasketHandler{Application: application.NewBasketApplication(br)}

	if assert.NoError(t, h.RemoveItemFromBasket(c)) {
		assert.Equal(t, http.StatusOK, rec.Code)
		b, _ := br.FindByShopId(1)
		assert.Equal(t, []*basket.BasketItem{{IngredientId: "ing-2"}}, b.Items)
	}
}
