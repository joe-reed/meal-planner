package shops_test

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/joe-reed/meal-planner/apps/api/shops"
	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
)

func TestStartingShop(t *testing.T) {
	shop1 := shops.Shop{Id: 1}

	r := shops.NewFakeShopRepository()
	r.Add(&shop1)

	e := echo.New()
	req := httptest.NewRequest("POST", "/shops", nil)
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	h := &shops.Handler{ShopRepository: r}

	if assert.NoError(t, h.StartShop(c)) {
		assert.Equal(t, http.StatusOK, rec.Code)
		c, _ := r.Current()
		assert.Equal(t, c.Id, 2)
	}
}
