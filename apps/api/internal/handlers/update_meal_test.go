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

func TestUpdatingMealName(t *testing.T) {
	repo := meal.NewFakeMealRepository()

	mealId := "123"

	err := repo.Save(meal.NewMealBuilder().WithId(mealId).WithName("foo").WithUrl("foo.localhost").Build())
	assert.NoError(t, err)

	e := echo.New()
	req := httptest.NewRequest("PATCH", "/meals/123", strings.NewReader(`{"name": "bar"}`))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.SetParamNames("mealId")
	c.SetParamValues("123")
	h := &handlers.MealsHandler{Application: application.NewMealApplication(repo)}

	if assert.NoError(t, h.UpdateMeal(c)) {
		m, err := repo.Get()
		assert.NoError(t, err)
		assert.Len(t, m, 1)
		assert.Equal(t, http.StatusOK, rec.Code)
		assert.Equal(t, "{\"id\":\"123\",\"name\":\"bar\",\"url\":\"foo.localhost\",\"ingredients\":[]}\n", rec.Body.String())
		assert.EqualExportedValues(t, &meal.Meal{Id: "123", Name: "bar", Url: "foo.localhost", MealIngredients: make([]meal.MealIngredient, 0)}, m[0])
	}
}

func TestUpdatingMealUrl(t *testing.T) {
	repo := meal.NewFakeMealRepository()

	mealId := "123"

	err := repo.Save(meal.NewMealBuilder().WithId(mealId).WithName("foo").Build())
	assert.NoError(t, err)

	e := echo.New()
	req := httptest.NewRequest("PATCH", "/meals/123", strings.NewReader(`{"url": "https://test.localhost"}`))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.SetParamNames("mealId")
	c.SetParamValues("123")
	h := &handlers.MealsHandler{Application: application.NewMealApplication(repo)}

	if assert.NoError(t, h.UpdateMeal(c)) {
		m, err := repo.Get()
		assert.NoError(t, err)
		assert.Len(t, m, 1)
		assert.Equal(t, http.StatusOK, rec.Code)
		assert.Equal(t, "{\"id\":\"123\",\"name\":\"foo\",\"url\":\"https://test.localhost\",\"ingredients\":[]}\n", rec.Body.String())
		assert.EqualExportedValues(t, &meal.Meal{Id: "123", Name: "foo", Url: "https://test.localhost", MealIngredients: make([]meal.MealIngredient, 0)}, m[0])
	}
}

func TestUpdatingMealUrlToBeEmpty(t *testing.T) {
	repo := meal.NewFakeMealRepository()

	mealId := "123"

	err := repo.Save(meal.NewMealBuilder().WithId(mealId).WithName("foo").WithUrl("https://test.localhost").Build())
	assert.NoError(t, err)

	e := echo.New()
	req := httptest.NewRequest("PATCH", "/meals/123", strings.NewReader(`{"url": ""}`))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.SetParamNames("mealId")
	c.SetParamValues("123")
	h := &handlers.MealsHandler{Application: application.NewMealApplication(repo)}

	if assert.NoError(t, h.UpdateMeal(c)) {
		m, err := repo.Get()
		assert.NoError(t, err)
		assert.Len(t, m, 1)
		assert.Equal(t, http.StatusOK, rec.Code)
		assert.Equal(t, "{\"id\":\"123\",\"name\":\"foo\",\"url\":\"\",\"ingredients\":[]}\n", rec.Body.String())
		assert.EqualExportedValues(t, &meal.Meal{Id: "123", Name: "foo", Url: "", MealIngredients: make([]meal.MealIngredient, 0)}, m[0])
	}
}

func TestUpdatingMealWithUrlAndName(t *testing.T) {
	repo := meal.NewFakeMealRepository()

	mealId := "123"

	err := repo.Save(meal.NewMealBuilder().WithId(mealId).WithName("foo").WithUrl("http://foo.localhost").Build())
	assert.NoError(t, err)

	e := echo.New()
	req := httptest.NewRequest("PATCH", "/meals/123", strings.NewReader(`{"url": "https://bar.localhost", "name": "bar"}`))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.SetParamNames("mealId")
	c.SetParamValues("123")
	h := &handlers.MealsHandler{Application: application.NewMealApplication(repo)}

	if assert.NoError(t, h.UpdateMeal(c)) {
		m, err := repo.Get()
		assert.NoError(t, err)
		assert.Len(t, m, 1)
		assert.Equal(t, http.StatusOK, rec.Code)
		assert.Equal(t, "{\"id\":\"123\",\"name\":\"bar\",\"url\":\"https://bar.localhost\",\"ingredients\":[]}\n", rec.Body.String())
		assert.EqualExportedValues(t, &meal.Meal{Id: "123", Name: "bar", Url: "https://bar.localhost", MealIngredients: make([]meal.MealIngredient, 0)}, m[0])
	}
}
