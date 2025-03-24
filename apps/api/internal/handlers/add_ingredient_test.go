package handlers_test

import (
	"github.com/joe-reed/meal-planner/apps/api/internal/domain/categories"
	"github.com/joe-reed/meal-planner/apps/api/internal/domain/ingredients"
	"github.com/joe-reed/meal-planner/apps/api/internal/handlers"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
)

func TestAddingIngredient(t *testing.T) {
	repo := ingredients.NewFakeIngredientRepository()

	e := echo.New()
	req := httptest.NewRequest("POST", "/ingredients", strings.NewReader(`{"id": "123","name":"foo","category":"Fruit"}`))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	h := &handlers.IngredientsHandler{IngredientRepository: repo}

	if assert.NoError(t, h.AddIngredient(c)) {
		m, err := repo.Get()
		assert.NoError(t, err)
		assert.Len(t, m, 1)
		assert.Equal(t, http.StatusAccepted, rec.Code)
		assert.Equal(t, "{\"id\":\"123\",\"name\":\"foo\",\"category\":\"Fruit\"}\n", rec.Body.String())
		assert.EqualExportedValues(t, &ingredients.Ingredient{Id: "123", Name: "foo", Category: categories.Fruit}, m[0])
	}
}
