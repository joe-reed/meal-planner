package handlers_test

import (
	"errors"
	"github.com/joe-reed/meal-planner/apps/api/internal/application"
	"github.com/joe-reed/meal-planner/apps/api/internal/domain/category"
	"github.com/joe-reed/meal-planner/apps/api/internal/domain/ingredient"
	"github.com/joe-reed/meal-planner/apps/api/internal/handlers"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
)

func TestAddingIngredient(t *testing.T) {
	repo := ingredient.NewFakeIngredientRepository()

	e := echo.New()
	req := httptest.NewRequest("POST", "/ingredients", strings.NewReader(`{"id": "123","name":"foo","category":"Fruit"}`))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	h := &handlers.IngredientsHandler{Application: application.NewIngredientApplication(repo)}

	if assert.NoError(t, h.AddIngredient(c)) {
		m, err := repo.Get()
		assert.NoError(t, err)
		assert.Len(t, m, 1)
		assert.Equal(t, http.StatusAccepted, rec.Code)
		assert.Equal(t, "{\"id\":\"123\",\"name\":\"foo\",\"category\":\"Fruit\"}\n", rec.Body.String())
		assert.EqualExportedValues(t, &ingredient.Ingredient{Id: "123", Name: "foo", Category: category.Fruit}, m[0])
	}
}

func TestAddingIngredientWithEmptyCategory(t *testing.T) {
	repo := ingredient.NewFakeIngredientRepository()

	e := echo.New()
	req := httptest.NewRequest("POST", "/ingredients", strings.NewReader(`{"id": "123","name":"foo","category":""}`))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	h := &handlers.IngredientsHandler{Application: application.NewIngredientApplication(repo)}

	if assert.NoError(t, h.AddIngredient(c)) {
		m, err := repo.Get()
		assert.NoError(t, err)
		assert.Len(t, m, 0)
		assert.Equal(t, http.StatusBadRequest, rec.Code)
	}
}

func TestAddingIngredientWithEmptyName(t *testing.T) {
	repo := ingredient.NewFakeIngredientRepository()

	e := echo.New()
	req := httptest.NewRequest("POST", "/ingredients", strings.NewReader(`{"id": "123","name":"","category":"Fruit"}`))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	h := &handlers.IngredientsHandler{Application: application.NewIngredientApplication(repo)}

	if assert.NoError(t, h.AddIngredient(c)) {
		m, err := repo.Get()
		assert.NoError(t, err)
		assert.Len(t, m, 0)
		assert.Equal(t, http.StatusBadRequest, rec.Code)
	}
}

func TestAddingIngredientWithEmptyId(t *testing.T) {
	repo := ingredient.NewFakeIngredientRepository()

	e := echo.New()
	req := httptest.NewRequest("POST", "/ingredients", strings.NewReader(`{"id": "","name":"foo","category":"Fruit"}`))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	h := &handlers.IngredientsHandler{Application: application.NewIngredientApplication(repo)}

	if assert.NoError(t, h.AddIngredient(c)) {
		m, err := repo.Get()
		assert.NoError(t, err)
		assert.Len(t, m, 0)
		assert.Equal(t, http.StatusBadRequest, rec.Code)
	}
}

type RepoWithError struct{}

func (r RepoWithError) Add(i *ingredient.Ingredient) error {
	return errors.New("error")
}
func (r RepoWithError) Get() ([]*ingredient.Ingredient, error) {
	return nil, errors.New("error")
}
func (r RepoWithError) GetByName(name ingredient.IngredientName) (*ingredient.Ingredient, error) {
	return nil, errors.New("error")
}

func TestHandlingUnknownError(t *testing.T) {
	repo := RepoWithError{}

	e := echo.New()
	req := httptest.NewRequest("POST", "/ingredients", strings.NewReader(`{"id": "123","name":"foo","category":"Fruit"}`))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	h := &handlers.IngredientsHandler{Application: application.NewIngredientApplication(repo)}

	assert.Error(t, h.AddIngredient(c), "error")
}
