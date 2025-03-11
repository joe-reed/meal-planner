package meals_test

import (
	"bytes"
	"github.com/joe-reed/meal-planner/apps/api/ingredients"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/joe-reed/meal-planner/apps/api/meals"
	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/require"
)

func TestUploadingMeals(t *testing.T) {
	repo := meals.NewFakeMealRepository()
	ingredientRepo := ingredients.NewFakeIngredientRepository()

	err := ingredientRepo.Add(ingredients.NewIngredientBuilder().WithName("Abc Name").WithId("abc").Build())
	require.NoError(t, err)

	err = ingredientRepo.Add(ingredients.NewIngredientBuilder().WithName("Def Name").WithId("def").Build())
	require.NoError(t, err)

	err = ingredientRepo.Add(ingredients.NewIngredientBuilder().WithName("Ghi Name").WithId("ghi").Build())
	require.NoError(t, err)

	e := echo.New()
	body := &bytes.Buffer{}
	w := multipart.NewWriter(body)

	part, err := w.CreateFormFile("meals", "meals.csv")
	require.NoError(t, err)

	_, err = part.Write([]byte("name,ingredient,amount,unit\nfoo,Abc Name,300,Gram\nfoo,Def Name,5,Tbsp\nbar,Def Name,400,Gram\nbar,Ghi Name,6,Tbsp"))
	require.NoError(t, err)

	err = w.Close()
	require.NoError(t, err)

	req := httptest.NewRequest("POST", "/meals/upload", body)

	req.Header.Set(echo.HeaderContentType, w.FormDataContentType())

	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	h := &meals.Handler{MealRepository: repo, IngredientRepository: ingredientRepo}

	err = h.UploadMeals(c)

	require.NoError(t, err)

	require.Equal(t, http.StatusCreated, rec.Code)

	m, err := repo.Get()

	require.NoError(t, err)
	require.Len(t, m, 2)

	require.Equal(t, m[0].Name, "bar")
	require.Len(t, m[0].MealIngredients, 2)
	require.Equal(t, []meals.MealIngredient{
		*meals.NewMealIngredient("def").WithQuantity(400, meals.Gram),
		*meals.NewMealIngredient("ghi").WithQuantity(6, meals.Tbsp),
	}, m[0].MealIngredients)

	require.Equal(t, m[1].Name, "foo")
	require.Len(t, m[1].MealIngredients, 2)
	require.Equal(t, []meals.MealIngredient{
		*meals.NewMealIngredient("abc").WithQuantity(300, meals.Gram),
		*meals.NewMealIngredient("def").WithQuantity(5, meals.Tbsp),
	}, m[1].MealIngredients)
}

func TestIngredientsNotExisting(t *testing.T) {
	repo := meals.NewFakeMealRepository()

	e := echo.New()
	body := &bytes.Buffer{}
	w := multipart.NewWriter(body)

	part, err := w.CreateFormFile("meals", "meals.csv")
	require.NoError(t, err)

	_, err = part.Write([]byte("name,ingredient,amount,unit\nfoo,Abc Name,300,Gram\nfoo,Def Name,5,Tbsp\nbar,Def Name,400,Gram\nbar,Ghi Name,6,Tbsp"))
	require.NoError(t, err)

	err = w.Close()
	require.NoError(t, err)

	req := httptest.NewRequest("POST", "/meals/upload", body)

	req.Header.Set(echo.HeaderContentType, w.FormDataContentType())

	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	h := &meals.Handler{MealRepository: repo, IngredientRepository: ingredients.NewFakeIngredientRepository()}

	err = h.UploadMeals(c)

	require.NoError(t, err)

	require.Equal(t, http.StatusBadRequest, rec.Code)

	m, err := repo.Get()

	require.Len(t, m, 0)

	require.Equal(t, "{\"notFoundIngredients\":[\"Abc Name\",\"Def Name\",\"Def Name\",\"Ghi Name\"]}\n", rec.Body.String())
}
