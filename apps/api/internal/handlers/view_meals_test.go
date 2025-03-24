package handlers_test

import (
	"fmt"
	"github.com/joe-reed/meal-planner/apps/api/internal/handlers"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/joe-reed/meal-planner/apps/api/internal/domain/meals"
	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
)

func TestViewingMeals(t *testing.T) {
	meal1 := meals.NewMealBuilder().WithName("Burritos").Build()
	meal2 := meals.NewMealBuilder().WithName("Shepherd's pie").Build()
	meal3 := meals.NewMealBuilder().WithName("Tacos").Build()

	repo := meals.NewFakeMealRepository()
	err := repo.Save(meal1)
	assert.NoError(t, err)
	err = repo.Save(meal2)
	assert.NoError(t, err)
	err = repo.Save(meal3)
	assert.NoError(t, err)

	e := echo.New()
	req := httptest.NewRequest("", "/meals", nil)
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	h := &handlers.MealsHandler{MealRepository: repo}

	if assert.NoError(t, h.GetMeals(c)) {
		assert.Equal(t, http.StatusOK, rec.Code)
		assert.Equal(t, fmt.Sprintf(`[{"id":"%s","name":"Burritos","ingredients":[]},{"id":"%s","name":"Shepherd's pie","ingredients":[]},{"id":"%s","name":"Tacos","ingredients":[]}]`+"\n", meal1.Id, meal2.Id, meal3.Id), rec.Body.String())
	}
}
