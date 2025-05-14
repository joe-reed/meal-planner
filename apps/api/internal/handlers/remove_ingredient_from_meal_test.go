package handlers_test

import (
	"github.com/joe-reed/meal-planner/apps/api/internal/application"
	"github.com/joe-reed/meal-planner/apps/api/internal/domain/meal"
	"github.com/joe-reed/meal-planner/apps/api/internal/handlers"
	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestRemovingIngredientFromMeal(t *testing.T) {
	repo := meal.NewFakeMealRepository()

	err := repo.Save(meal.NewMealBuilder().
		WithId("123").
		AddIngredient(meal.Ingredient{ProductId: "ing-1"}).
		AddIngredient(meal.Ingredient{ProductId: "ing-2"}).
		Build())

	assert.NoError(t, err)

	e := echo.New()
	req := httptest.NewRequest("DELETE", "/meals/123/ingredients/ing-1", nil)
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.SetParamNames("mealId", "ingredientId")
	c.SetParamValues("123", "ing-1")
	h := &handlers.MealsHandler{Application: application.NewMealApplication(repo)}

	if assert.NoError(t, h.RemoveIngredientFromMeal(c)) {
		m, err := repo.Find("123")
		assert.NoError(t, err)
		assert.Equal(t, http.StatusOK, rec.Code)
		assert.Len(t, m.Ingredients, 1)
		assert.Equal(t, []meal.Ingredient{{ProductId: "ing-2"}}, m.Ingredients)
	}
}

func TestRemovingAllIngredientsFromMeal(t *testing.T) {
	repo := meal.NewFakeMealRepository()

	err := repo.Save(meal.NewMealBuilder().
		WithId("123").
		AddIngredient(meal.Ingredient{ProductId: "ing-1"}).
		Build())

	assert.NoError(t, err)

	e := echo.New()
	req := httptest.NewRequest("DELETE", "/meals/123/ingredients/ing-1", nil)
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.SetParamNames("mealId", "ingredientId")
	c.SetParamValues("123", "ing-1")
	h := &handlers.MealsHandler{Application: application.NewMealApplication(repo)}

	if assert.NoError(t, h.RemoveIngredientFromMeal(c)) {
		m, err := repo.Find("123")
		assert.NoError(t, err)
		assert.Equal(t, http.StatusOK, rec.Code)
		assert.Len(t, m.Ingredients, 0)
		assert.Equal(t, []meal.Ingredient{}, m.Ingredients)
	}
}
