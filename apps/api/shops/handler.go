package shops

import (
	"net/http"

	"github.com/labstack/echo/v4"
)

type Handler struct {
	ShopRepository ShopRepository
}

func (h *Handler) CurrentShop(c echo.Context) error {
	shop, err := h.ShopRepository.Current()

	if err != nil {
		return err
	}

	return c.JSON(http.StatusOK, shop)
}
