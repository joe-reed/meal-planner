package handlers_test

import (
	"github.com/joe-reed/meal-planner/apps/api/internal/domain/shops"
	"github.com/joe-reed/meal-planner/apps/api/internal/handlers"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
)

func TestStartingShop(t *testing.T) {
	shop1, err := shops.NewShop(1)
	assert.NoError(t, err)

	r := shops.NewFakeShopRepository()
	err = r.Save(shop1)
	assert.NoError(t, err)

	e := echo.New()
	req := httptest.NewRequest("POST", "/shops", nil)
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	h := &handlers.ShopsHandler{ShopRepository: r, Publisher: func(string) {}}

	if assert.NoError(t, h.StartShop(c)) {
		assert.Equal(t, http.StatusOK, rec.Code)
		c, err := r.Current()
		assert.NoError(t, err)
		assert.Equal(t, c.Id, 2)
	}
}
