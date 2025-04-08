package handlers_test

import (
	"fmt"
	"github.com/joe-reed/meal-planner/apps/api/internal/application"
	"github.com/joe-reed/meal-planner/apps/api/internal/handlers"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/joe-reed/meal-planner/apps/api/internal/domain/meals"
	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
)

func TestViewingMeal(t *testing.T) {
	meal := meals.NewMealBuilder().WithName("Burritos").AddIngredient(*meals.NewMealIngredient("ing-123")).Build()

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
	h := &handlers.MealsHandler{Application: application.NewMealApplication(repo)}

	if assert.NoError(t, h.FindMeal(c)) {
		assert.Equal(t, http.StatusOK, rec.Code)
		assert.Equal(t, fmt.Sprintf(`{"id":"%s","name":"Burritos","ingredients":[{"id":"ing-123","quantity":{"amount":1,"unit":"Number"}}]}`+"\n", meal.Id), rec.Body.String())
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
	h := &handlers.MealsHandler{Application: application.NewMealApplication(repo)}

	if assert.NoError(t, h.FindMeal(c)) {
		assert.Equal(t, http.StatusOK, rec.Code)
		assert.Equal(t, fmt.Sprintf(`{"id":"%s","name":"Burritos","ingredients":[]}`+"\n", meal.Id), rec.Body.String())
	}
}
