package handlers_test

import (
	"github.com/joe-reed/meal-planner/apps/api/internal/application"
	"github.com/joe-reed/meal-planner/apps/api/internal/domain/meal"
	"github.com/joe-reed/meal-planner/apps/api/internal/domain/quantity"
	"github.com/joe-reed/meal-planner/apps/api/internal/handlers"
	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestAddingIngredientToMeal(t *testing.T) {
	repo := meal.NewFakeMealRepository()

	err := repo.Save(meal.NewMealBuilder().WithId("123").Build())
	assert.NoError(t, err)

	e := echo.New()
	req := httptest.NewRequest("POST", "/meals/123/ingredients", strings.NewReader(`{"id": "ing-1", "quantity": {"amount": 3, "unit": "Cup"}}`))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.SetParamNames("mealId")
	c.SetParamValues("123")
	h := &handlers.MealsHandler{Application: application.NewMealApplication(repo)}

	if assert.NoError(t, h.AddIngredientToMeal(c)) {
		m, err := repo.Find("123")
		assert.NoError(t, err)
		assert.Equal(t, http.StatusOK, rec.Code)
		assert.Equal(t, []meal.Ingredient{{ProductId: "ing-1", Quantity: quantity.Quantity{Amount: 3, Unit: quantity.Cup}}}, m.Ingredients)
	}
}

func TestAddingIngredientToMealWithoutQuantity(t *testing.T) {
	repo := meal.NewFakeMealRepository()

	err := repo.Save(meal.NewMealBuilder().WithId("123").Build())
	assert.NoError(t, err)

	e := echo.New()
	req := httptest.NewRequest("POST", "/meals/123/ingredients", strings.NewReader(`{"id": "ing-1" }`))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.SetParamNames("mealId")
	c.SetParamValues("123")
	h := &handlers.MealsHandler{Application: application.NewMealApplication(repo)}

	if assert.NoError(t, h.AddIngredientToMeal(c)) {
		m, err := repo.Find("123")
		assert.NoError(t, err)
		assert.Equal(t, http.StatusOK, rec.Code)
		assert.Equal(t, []meal.Ingredient{{ProductId: "ing-1", Quantity: quantity.Quantity{Amount: 1, Unit: quantity.Number}}}, m.Ingredients)
	}
}
