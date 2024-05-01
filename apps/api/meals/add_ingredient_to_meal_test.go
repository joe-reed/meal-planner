package meals_test

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/joe-reed/meal-planner/apps/api/meals"
	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
)

func TestAddingIngredientToMeal(t *testing.T) {
	repo := meals.NewFakeMealRepository()

	err := repo.Add(meals.NewMealBuilder().WithId("123").Build())
	assert.NoError(t, err)

	e := echo.New()
	req := httptest.NewRequest("POST", "/meals/123/ingredients", strings.NewReader(`{"id": "ing-1"}`))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.SetParamNames("mealId")
	c.SetParamValues("123")
	h := &meals.Handler{MealRepository: repo}

	if assert.NoError(t, h.AddIngredientToMeal(c)) {
		m, err := repo.Find("123")
		assert.NoError(t, err)
		assert.Equal(t, http.StatusOK, rec.Code)
		assert.Equal(t, []meals.MealIngredient{{IngredientId: "ing-1"}}, m.MealIngredients)
	}
}
