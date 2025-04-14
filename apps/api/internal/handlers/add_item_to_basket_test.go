package handlers_test

import (
	"github.com/joe-reed/meal-planner/apps/api/internal/application"
	"github.com/joe-reed/meal-planner/apps/api/internal/domain/basket"
	"github.com/joe-reed/meal-planner/apps/api/internal/handlers"
	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestAddingItemToBasket(t *testing.T) {
	b, err := basket.NewBasket(1)
	assert.NoError(t, err)

	br := basket.NewFakeBasketRepository()

	err = br.Save(b)
	assert.NoError(t, err)

	e := echo.New()
	req := httptest.NewRequest("POST", "/baskets/1/items", strings.NewReader(`{"ingredientId":"ing-1"}`))

	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.SetParamNames("shopId")
	c.SetParamValues("1")
	h := &handlers.BasketHandler{Application: application.NewBasketApplication(br)}

	if assert.NoError(t, h.AddItemToBasket(c)) {
		assert.Equal(t, http.StatusOK, rec.Code)
		b, _ := br.FindByShopId(1)
		assert.Equal(t, []*basket.BasketItem{{IngredientId: "ing-1"}}, b.Items)
	}
}
