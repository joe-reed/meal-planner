package handlers

import (
	"errors"
	"github.com/joe-reed/meal-planner/apps/api/internal/application"
	"github.com/joe-reed/meal-planner/apps/api/internal/domain/product"
	"github.com/labstack/echo/v4"
	"net/http"
)

type ProductHandler struct {
	Application *application.ProductApplication
}

func (h *ProductHandler) GetProducts(c echo.Context) error {
	products, err := h.Application.GetProducts()

	if err != nil {
		return err
	}

	return c.JSON(http.StatusOK, products)
}

func (h *ProductHandler) AddProduct(c echo.Context) error {
	body := new(product.Product)

	if err := c.Bind(body); err != nil {
		return c.JSON(
			http.StatusBadRequest,
			struct {
				Error string `json:"error"`
			}{
				Error: "Invalid request body",
			})
	}

	p, err := h.Application.AddProduct(body.Id, body.Name, body.Category)

	if err != nil {
		var validationError *application.ValidationError
		if errors.As(err, &validationError) {
			return c.JSON(http.StatusBadRequest, struct {
				Error string `json:"error"`
			}{
				Error: validationError.Error(),
			})
		}
		return err
	}

	return c.JSON(http.StatusAccepted, p)
}
