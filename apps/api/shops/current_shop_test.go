package shops_test

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/joe-reed/meal-planner/apps/api/shops"
	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
)

func TestGettingCurrentShop(t *testing.T) {
	shop1 := shops.Shop{Id: 1}
	shop2 := shops.Shop{Id: 2}

	r := shops.NewFakeShopRepository()
	r.Add(&shop1)
	r.Add(&shop2)

	e := echo.New()
	req := httptest.NewRequest("GET", "/shops/current", nil)
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	h := &shops.Handler{ShopRepository: r}

	if assert.NoError(t, h.CurrentShop(c)) {
		assert.Equal(t, http.StatusOK, rec.Code)
		assert.Equal(t, `{"id":2}`+"\n", rec.Body.String())
	}
}