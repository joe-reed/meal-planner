package handlers_test

import (
	"github.com/joe-reed/meal-planner/apps/api/internal/application"
	"github.com/joe-reed/meal-planner/apps/api/internal/domain/quantity"
	"github.com/joe-reed/meal-planner/apps/api/internal/domain/shop"
	"github.com/joe-reed/meal-planner/apps/api/internal/handlers"
	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestRemovingItemFromCurrentShop(t *testing.T) {
	s, err := shop.NewShop(1)
	assert.NoError(t, err)

	s.AddItem(
		&shop.Item{
			ProductId: "abc",
			Quantity:  quantity.Quantity{Amount: 3, Unit: quantity.Cup},
		},
	)
	s.AddItem(
		&shop.Item{
			ProductId: "xyz",
			Quantity:  quantity.Quantity{Amount: 2, Unit: quantity.Gram},
		},
	)

	r := shop.NewFakeShopRepository()
	err = r.Save(s)
	assert.NoError(t, err)

	e := echo.New()
	req := httptest.NewRequest("DELETE", "/shops/current/items/abc", nil)

	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.SetParamNames("productId")
	c.SetParamValues("abc")
	h := &handlers.ShopsHandler{Application: application.NewShopApplication(r, func(string) {})}

	if assert.NoError(t, h.RemoveItemFromCurrentShop(c)) {
		assert.Equal(t, http.StatusOK, rec.Code)
		s, _ := r.Find(1)
		assert.Equal(t, []*shop.Item{
			{
				ProductId: "xyz",
				Quantity:  quantity.Quantity{Amount: 2, Unit: quantity.Gram},
			},
		}, s.Items)
	}
}
