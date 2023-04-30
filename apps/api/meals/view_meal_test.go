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
	meal := meals.NewMealBuilder().WithName("Burritos").Build()

	repo := meals.NewFakeMealRepository()
	repo.Add(meal)

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
		assert.Equal(t, fmt.Sprintf(`{"id":"%s","name":"Burritos"}`+"\n", meal.Id), rec.Body.String())
	}
}
