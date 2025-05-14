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
	"strings"
	"testing"
)

func TestAddingItemToCurrentShop(t *testing.T) {
	s, err := shop.NewShop(1)
	assert.NoError(t, err)

	r := shop.NewFakeShopRepository()
	err = r.Save(s)
	assert.NoError(t, err)

	e := echo.New()
	req := httptest.NewRequest("POST", "/shops/current/items", strings.NewReader(`{"productId": "abc", "quantity": {"amount": 3, "unit": "Cup"}}`))

	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	h := &handlers.ShopsHandler{Application: application.NewShopApplication(r, func(string) {})}

	if assert.NoError(t, h.AddItemToCurrentShop(c)) {
		assert.Equal(t, http.StatusOK, rec.Code)
		s, _ := r.Find(1)
		assert.Equal(t, []*shop.Item{{ProductId: "abc", Quantity: quantity.Quantity{Amount: 3, Unit: quantity.Cup}}}, s.Items)
	}
}
