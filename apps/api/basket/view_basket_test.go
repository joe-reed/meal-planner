package basket_test

import (
	"github.com/joe-reed/meal-planner/apps/api/basket"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
)

func TestViewingBasket(t *testing.T) {
	b, err := basket.NewBasket(1)
	assert.NoError(t, err)

	b.AddItem(&basket.BasketItem{IngredientId: "ing-1"})
	b.AddItem(&basket.BasketItem{IngredientId: "ing-2"})

	br := basket.NewFakeBasketRepository()

	err = br.Save(b)
	assert.NoError(t, err)

	e := echo.New()
	req := httptest.NewRequest("GET", "/baskets/1", nil)

	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.SetParamNames("shopId")
	c.SetParamValues("1")
	h := &basket.Handler{BasketRepository: br}

	if assert.NoError(t, h.GetBasketItems(c)) {
		assert.Equal(t, http.StatusOK, rec.Code)
		assert.Equal(t, `{"shopId":1,"items":[{"ingredientId":"ing-1"},{"ingredientId":"ing-2"}]}`+"\n", rec.Body.String())
	}
}

func TestViewingBasketWithNoItems(t *testing.T) {
	b, err := basket.NewBasket(1)
	assert.NoError(t, err)

	br := basket.NewFakeBasketRepository()

	err = br.Save(b)
	assert.NoError(t, err)

	e := echo.New()
	req := httptest.NewRequest("GET", "/baskets/1", nil)

	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.SetParamNames("shopId")
	c.SetParamValues("1")
	h := &basket.Handler{BasketRepository: br}

	if assert.NoError(t, h.GetBasketItems(c)) {
		assert.Equal(t, http.StatusOK, rec.Code)
		assert.Equal(t, `{"shopId":1,"items":[]}`+"\n", rec.Body.String())
	}
}
