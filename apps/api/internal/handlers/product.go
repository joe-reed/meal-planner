package handlers

import (
	"context"
	"errors"
	sqlStore "github.com/hallgren/eventsourcing/eventstore/sql"
	"github.com/joe-reed/meal-planner/apps/api/internal/application"
	"github.com/joe-reed/meal-planner/apps/api/internal/domain/product"
	"github.com/joe-reed/meal-planner/apps/api/internal/projections"
	"github.com/labstack/echo/v4"
	"net/http"
)

type ProductHandler struct {
	Application *application.ProductApplication
	EventStore  *sqlStore.SQLite
}

func (h *ProductHandler) GetProducts(c echo.Context) error {
	if c.QueryParam("grouped") == "true" {
		p, output := projections.CreateProductProjection(h.EventStore)

		result := p.RunToEnd(context.TODO())

		if result.Error != nil {
			return result.Error
		}

		return c.JSON(200, output)
	}

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
				Error: "Invalid request body: " + err.Error(),
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

		var productAlreadyExists *application.ProductAlreadyExists
		if errors.As(err, &productAlreadyExists) {
			return c.JSON(http.StatusConflict, struct {
				Error       string `json:"error"`
				ProductName string `json:"productName"`
			}{
				Error:       productAlreadyExists.Error(),
				ProductName: productAlreadyExists.ProductName,
			})
		}
		return err
	}

	return c.JSON(http.StatusAccepted, p)
}
