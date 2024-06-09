package meals_test

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/joe-reed/meal-planner/apps/api/meals"
	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
)

func TestViewingMeal(t *testing.T) {
	meal := meals.NewMealBuilder().WithName("Burritos").AddIngredient(meals.MealIngredient{IngredientId: "ing-123"}).Build()

	repo := meals.NewFakeMealRepository()
	err := repo.Save(meal)
	assert.NoError(t, err)

	e := echo.New()
	req := httptest.NewRequest("GET", "/meals/"+meal.Id, nil)
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.SetParamNames("id")
	c.SetParamValues(meal.Id)
	h := &meals.Handler{MealRepository: repo}

	if assert.NoError(t, h.GetMeal(c)) {
		assert.Equal(t, http.StatusOK, rec.Code)
		assert.Equal(t, fmt.Sprintf(`{"id":"%s","name":"Burritos","ingredients":[{"id":"ing-123"}]}`+"\n", meal.Id), rec.Body.String())
	}
}

func TestViewingMealWithNoIngredients(t *testing.T) {
	meal := meals.NewMealBuilder().WithName("Burritos").Build()

	repo := meals.NewFakeMealRepository()
	err := repo.Save(meal)
	assert.NoError(t, err)

	e := echo.New()
	req := httptest.NewRequest("GET", "/meals/"+meal.Id, nil)
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.SetParamNames("id")
	c.SetParamValues(meal.Id)
	h := &meals.Handler{MealRepository: repo}

	if assert.NoError(t, h.GetMeal(c)) {
		assert.Equal(t, http.StatusOK, rec.Code)
		assert.Equal(t, fmt.Sprintf(`{"id":"%s","name":"Burritos","ingredients":[]}`+"\n", meal.Id), rec.Body.String())
	}
}
