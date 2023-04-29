package meals_test

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/joe-reed/meal-planner/server/meals"
	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
)

func TestAddingMeal(t *testing.T) {
	repo := meals.NewFakeMealRepository()

	e := echo.New()
	req := httptest.NewRequest("POST", "/", strings.NewReader(`{"id": "123","name":"foo"}`))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	h := &meals.Handler{MealRepository: repo}

	if assert.NoError(t, h.AddMeal(c)) {
		m, err := repo.Get()
		assert.NoError(t, err)
		assert.Equal(t, http.StatusAccepted, rec.Code)
		assert.Equal(t, "\"Meal added\"\n", rec.Body.String())
		assert.Equal(t, []*meals.Meal{{Id: "123", Name: "foo"}}, m)
	}
}
