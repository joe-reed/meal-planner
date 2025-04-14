package handlers_test

import (
	"github.com/joe-reed/meal-planner/apps/api/internal/application"
	"github.com/joe-reed/meal-planner/apps/api/internal/domain/category"
	"github.com/joe-reed/meal-planner/apps/api/internal/domain/ingredient"
	"github.com/joe-reed/meal-planner/apps/api/internal/handlers"
	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestViewingIngredients(t *testing.T) {
	i1 := ingredient.NewIngredientBuilder().WithName("Chicken").WithCategory(category.Meat).WithId("8a378ac5-e0d5-405a-8cb5-f03cc1d92d8b").Build()
	i2 := ingredient.NewIngredientBuilder().WithName("Ice Cream").WithCategory(category.Frozen).WithId("ad11289f-b2b0-4195-ba89-63d92ccc64d7").Build()
	i3 := ingredient.NewIngredientBuilder().WithName("Carrot").WithCategory(category.Vegetables).WithId("57b5d842-eda3-4368-a496-21956de5e254").Build()

	repo := ingredient.NewFakeIngredientRepository()
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
	h := &handlers.IngredientsHandler{Application: application.NewIngredientApplication(repo)}

	if assert.NoError(t, h.GetIngredients(c)) {
		assert.Equal(t, http.StatusOK, rec.Code)
		assert.Equal(t, `[{"id":"57b5d842-eda3-4368-a496-21956de5e254","name":"Carrot","category":"Vegetables"},{"id":"8a378ac5-e0d5-405a-8cb5-f03cc1d92d8b","name":"Chicken","category":"Meat"},{"id":"ad11289f-b2b0-4195-ba89-63d92ccc64d7","name":"Ice Cream","category":"Frozen"}]`+"\n", rec.Body.String())
	}
}
