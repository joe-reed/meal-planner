package meals_test

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/joe-reed/meal-planner/apps/api/meals"
	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
)

func TestRemovingIngredientFromMeal(t *testing.T) {
	repo := meals.NewFakeMealRepository()

	err := repo.Add(meals.NewMealBuilder().
		WithId("123").
		AddIngredient(meals.MealIngredient{IngredientId: "ing-1"}).
		AddIngredient(meals.MealIngredient{IngredientId: "ing-2"}).
		Build())

	assert.NoError(t, err)

	e := echo.New()
	req := httptest.NewRequest("DELETE", "/meals/123/ingredients/ing-1", nil)
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.SetParamNames("mealId", "ingredientId")
	c.SetParamValues("123", "ing-1")
	h := &meals.Handler{MealRepository: repo}

	if assert.NoError(t, h.RemoveIngredientFromMeal(c)) {
		m, err := repo.Find("123")
		assert.NoError(t, err)
		assert.Equal(t, http.StatusOK, rec.Code)
		assert.Equal(t, []meals.MealIngredient{{IngredientId: "ing-2"}}, m.MealIngredients)
	}
}
