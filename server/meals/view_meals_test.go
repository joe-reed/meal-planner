package meals_test

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/joe-reed/meal-planner/server/meals"
	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
)

func TestViewingMeals(t *testing.T) {
	meal1 := meals.NewMealBuilder().WithName("Burritos").Build()
	meal2 := meals.NewMealBuilder().WithName("Shepherd's pie").Build()
	meal3 := meals.NewMealBuilder().WithName("Tacos").Build()

	repo := meals.NewFakeMealRepository()
	repo.Add(meal1)
	repo.Add(meal2)
	repo.Add(meal3)

	e := echo.New()
	req := httptest.NewRequest("", "/", nil)
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	h := &meals.Handler{MealRepository: repo}

	if assert.NoError(t, h.GetMeals(c)) {
		assert.Equal(t, http.StatusOK, rec.Code)
		assert.Equal(t, fmt.Sprintf(`[{"id":"%s","name":"Burritos"},{"id":"%s","name":"Shepherd's pie"},{"id":"%s","name":"Tacos"}]`+"\n", meal1.Id, meal2.Id, meal3.Id), rec.Body.String())
	}
}
