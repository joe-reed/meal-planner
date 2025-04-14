package handlers_test

import (
	"github.com/joe-reed/meal-planner/apps/api/internal/application"
	"github.com/joe-reed/meal-planner/apps/api/internal/domain/meal"
	"github.com/joe-reed/meal-planner/apps/api/internal/handlers"
	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestAddingMeal(t *testing.T) {
	repo := meal.NewFakeMealRepository()

	e := echo.New()
	req := httptest.NewRequest("POST", "/meals", strings.NewReader(`{"id": "123","name":"foo"}`))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	h := &handlers.MealsHandler{Application: application.NewMealApplication(repo)}

	if assert.NoError(t, h.AddMeal(c)) {
		m, err := repo.Get()
		assert.NoError(t, err)
		assert.Len(t, m, 1)
		assert.Equal(t, http.StatusCreated, rec.Code)
		assert.Equal(t, "{\"id\":\"123\",\"name\":\"foo\",\"ingredients\":[]}\n", rec.Body.String())
		assert.EqualExportedValues(t, &meal.Meal{Id: "123", Name: "foo", MealIngredients: make([]meal.MealIngredient, 0)}, m[0])
	}
}

func TestAddingDuplicateMeal(t *testing.T) {
	repo := meal.NewFakeMealRepository()
	err := repo.Save(meal.NewMealBuilder().WithName("foo").Build())

	assert.NoError(t, err)

	e := echo.New()
	req := httptest.NewRequest("POST", "/meals", strings.NewReader(`{"id": "123","name":"foo"}`))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	h := &handlers.MealsHandler{Application: application.NewMealApplication(repo)}

	if assert.NoError(t, h.AddMeal(c)) {
		m, err := repo.Get()
		assert.NoError(t, err)
		assert.Len(t, m, 1)
		assert.Equal(t, http.StatusBadRequest, rec.Code)
	}
}
