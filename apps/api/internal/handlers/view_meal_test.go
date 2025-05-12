package handlers_test

import (
	"fmt"
	"github.com/joe-reed/meal-planner/apps/api/internal/application"
	"github.com/joe-reed/meal-planner/apps/api/internal/domain/meal"
	"github.com/joe-reed/meal-planner/apps/api/internal/handlers"
	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestViewingMeal(t *testing.T) {
	m := meal.NewMealBuilder().WithName("Burritos").AddIngredient(*meal.NewIngredient("ing-123")).Build()

	repo := meal.NewFakeMealRepository()
	err := repo.Save(m)
	assert.NoError(t, err)

	e := echo.New()
	req := httptest.NewRequest("GET", "/meals/"+m.Id, nil)
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.SetParamNames("id")
	c.SetParamValues(m.Id)
	h := &handlers.MealsHandler{Application: application.NewMealApplication(repo)}

	if assert.NoError(t, h.FindMeal(c)) {
		assert.Equal(t, http.StatusOK, rec.Code)
		assert.Equal(t, fmt.Sprintf(`{"id":"%s","name":"Burritos","url":"","ingredients":[{"id":"ing-123","quantity":{"amount":1,"unit":"Number"}}]}`+"\n", m.Id), rec.Body.String())
	}
}

func TestViewingMealWithNoIngredients(t *testing.T) {
	m := meal.NewMealBuilder().WithName("Burritos").Build()

	repo := meal.NewFakeMealRepository()
	err := repo.Save(m)
	assert.NoError(t, err)

	e := echo.New()
	req := httptest.NewRequest("GET", "/meals/"+m.Id, nil)
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.SetParamNames("id")
	c.SetParamValues(m.Id)
	h := &handlers.MealsHandler{Application: application.NewMealApplication(repo)}

	if assert.NoError(t, h.FindMeal(c)) {
		assert.Equal(t, http.StatusOK, rec.Code)
		assert.Equal(t, fmt.Sprintf(`{"id":"%s","name":"Burritos","url":"","ingredients":[]}`+"\n", m.Id), rec.Body.String())
	}
}
