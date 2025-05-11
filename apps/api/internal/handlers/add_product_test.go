package handlers_test

import (
	"errors"
	"github.com/joe-reed/meal-planner/apps/api/internal/application"
	"github.com/joe-reed/meal-planner/apps/api/internal/domain/category"
	"github.com/joe-reed/meal-planner/apps/api/internal/domain/product"
	"github.com/joe-reed/meal-planner/apps/api/internal/handlers"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
)

func TestAddingProduct(t *testing.T) {
	repo := product.NewFakeProductRepository()

	e := echo.New()
	req := httptest.NewRequest("POST", "/products", strings.NewReader(`{"id": "123","name":"foo","category":"Fruit"}`))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	h := &handlers.ProductHandler{Application: application.NewProductApplication(repo)}

	if assert.NoError(t, h.AddProduct(c)) {
		m, err := repo.Get()
		assert.NoError(t, err)
		assert.Len(t, m, 1)
		assert.Equal(t, http.StatusAccepted, rec.Code)
		assert.Equal(t, "{\"id\":\"123\",\"name\":\"foo\",\"category\":\"Fruit\"}\n", rec.Body.String())
		assert.EqualExportedValues(t, &product.Product{Id: "123", Name: "foo", Category: category.Fruit}, m[0])
	}
}

func TestAddingProductWithEmptyCategory(t *testing.T) {
	repo := product.NewFakeProductRepository()

	e := echo.New()
	req := httptest.NewRequest("POST", "/products", strings.NewReader(`{"id": "123","name":"foo","category":""}`))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	h := &handlers.ProductHandler{Application: application.NewProductApplication(repo)}

	if assert.NoError(t, h.AddProduct(c)) {
		m, err := repo.Get()
		assert.NoError(t, err)
		assert.Len(t, m, 0)
		assert.Equal(t, http.StatusBadRequest, rec.Code)
	}
}

func TestAddingProductWithEmptyName(t *testing.T) {
	repo := product.NewFakeProductRepository()

	e := echo.New()
	req := httptest.NewRequest("POST", "/products", strings.NewReader(`{"id": "123","name":"","category":"Fruit"}`))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	h := &handlers.ProductHandler{Application: application.NewProductApplication(repo)}

	if assert.NoError(t, h.AddProduct(c)) {
		m, err := repo.Get()
		assert.NoError(t, err)
		assert.Len(t, m, 0)
		assert.Equal(t, http.StatusBadRequest, rec.Code)
	}
}

func TestAddingProductWithEmptyId(t *testing.T) {
	repo := product.NewFakeProductRepository()

	e := echo.New()
	req := httptest.NewRequest("POST", "/products", strings.NewReader(`{"id": "","name":"foo","category":"Fruit"}`))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	h := &handlers.ProductHandler{Application: application.NewProductApplication(repo)}

	if assert.NoError(t, h.AddProduct(c)) {
		m, err := repo.Get()
		assert.NoError(t, err)
		assert.Len(t, m, 0)
		assert.Equal(t, http.StatusBadRequest, rec.Code)
	}
}

type ProductRepoWithError struct{}

func (r ProductRepoWithError) Add(p *product.Product) error {
	return errors.New("error")
}
func (r ProductRepoWithError) Get() ([]*product.Product, error) {
	return nil, errors.New("error")
}
func (r ProductRepoWithError) GetByName(name product.ProductName) (*product.Product, error) {
	return nil, errors.New("error")
}

func TestAddingProductWithUnknownError(t *testing.T) {
	repo := ProductRepoWithError{}

	e := echo.New()
	req := httptest.NewRequest("POST", "/products", strings.NewReader(`{"id": "123","name":"foo","category":"Fruit"}`))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	h := &handlers.ProductHandler{Application: application.NewProductApplication(repo)}

	assert.Error(t, h.AddProduct(c), "error")
}
