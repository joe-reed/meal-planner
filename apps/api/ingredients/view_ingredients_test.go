package ingredients_test

import (
	"github.com/joe-reed/meal-planner/apps/api/ingredients"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
)

func TestViewingIngredients(t *testing.T) {
	i1 := ingredients.NewIngredientBuilder().WithName("Onion").WithId("8a378ac5-e0d5-405a-8cb5-f03cc1d92d8b").Build()
	i2 := ingredients.NewIngredientBuilder().WithName("Potato").WithId("ad11289f-b2b0-4195-ba89-63d92ccc64d7").Build()
	i3 := ingredients.NewIngredientBuilder().WithName("Carrot").WithId("57b5d842-eda3-4368-a496-21956de5e254").Build()

	repo := ingredients.NewFakeIngredientRepository()
	err := repo.Add(i1)
	assert.NoError(t, err)
	err = repo.Add(i2)
	assert.NoError(t, err)
	err = repo.Add(i3)
	assert.NoError(t, err)

	e := echo.New()
	req := httptest.NewRequest("", "/", nil)
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	h := &ingredients.Handler{IngredientRepository: repo}

	if assert.NoError(t, h.GetIngredients(c)) {
		assert.Equal(t, http.StatusOK, rec.Code)
		assert.Equal(t, `[{"id":"57b5d842-eda3-4368-a496-21956de5e254","name":"Carrot"},{"id":"8a378ac5-e0d5-405a-8cb5-f03cc1d92d8b","name":"Onion"},{"id":"ad11289f-b2b0-4195-ba89-63d92ccc64d7","name":"Potato"}]`+"\n", rec.Body.String())
	}
}
