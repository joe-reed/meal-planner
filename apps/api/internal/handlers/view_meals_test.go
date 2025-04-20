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

func TestViewingMeals(t *testing.T) {
	meal1 := meal.NewMealBuilder().WithName("Burritos").Build()
	meal2 := meal.NewMealBuilder().WithName("Shepherd's pie").Build()
	meal3 := meal.NewMealBuilder().WithName("Tacos").Build()

	repo := meal.NewFakeMealRepository()
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
	h := &handlers.MealsHandler{Application: application.NewMealApplication(repo)}

	if assert.NoError(t, h.GetMeals(c)) {
		assert.Equal(t, http.StatusOK, rec.Code)
		assert.Equal(t, fmt.Sprintf(`[{"id":"%s","name":"Burritos","url":"","ingredients":[]},{"id":"%s","name":"Shepherd's pie","url":"","ingredients":[]},{"id":"%s","name":"Tacos","url":"","ingredients":[]}]`+"\n", meal1.Id, meal2.Id, meal3.Id), rec.Body.String())
	}
}
