package handlers_test

import (
	"github.com/joe-reed/meal-planner/apps/api/internal/handlers"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/joe-reed/meal-planner/apps/api/internal/domain/meals"
	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
)

func TestAddingMeal(t *testing.T) {
	repo := meals.NewFakeMealRepository()

	e := echo.New()
	req := httptest.NewRequest("POST", "/meals", strings.NewReader(`{"id": "123","name":"foo"}`))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	h := &handlers.MealsHandler{MealRepository: repo}

	if assert.NoError(t, h.AddMeal(c)) {
		m, err := repo.Get()
		assert.NoError(t, err)
		assert.Len(t, m, 1)
		assert.Equal(t, http.StatusCreated, rec.Code)
		assert.Equal(t, "{\"id\":\"123\",\"name\":\"foo\",\"ingredients\":[]}\n", rec.Body.String())
		assert.EqualExportedValues(t, &meals.Meal{Id: "123", Name: "foo", MealIngredients: make([]meals.MealIngredient, 0)}, m[0])
	}
}

func TestAddingDuplicateMeal(t *testing.T) {
	repo := meals.NewFakeMealRepository()
	err := repo.Save(meals.NewMealBuilder().WithName("foo").Build())

	assert.NoError(t, err)

	e := echo.New()
	req := httptest.NewRequest("POST", "/meals", strings.NewReader(`{"id": "123","name":"foo"}`))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	h := &handlers.MealsHandler{MealRepository: repo}

	if assert.NoError(t, h.AddMeal(c)) {
		m, err := repo.Get()
		assert.NoError(t, err)
		assert.Len(t, m, 1)
		assert.Equal(t, http.StatusBadRequest, rec.Code)
	}
}
