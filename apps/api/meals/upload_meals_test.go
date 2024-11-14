package meals_test

import (
	"bytes"
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

	e := echo.New()
	body := &bytes.Buffer{}
	w := multipart.NewWriter(body)

	part, err := w.CreateFormFile("meals", "meals.csv")
	require.NoError(t, err)

	_, err = part.Write([]byte("name,ingredient,amount,unit\nfoo,abc,300,Gram\nfoo,def,5,Tbsp\nbar,def,400,Gram\nbar,ghi,6,Tbsp"))
	require.NoError(t, err)

	err = w.Close()
	require.NoError(t, err)

	req := httptest.NewRequest("POST", "/meals/upload", body)

	req.Header.Set(echo.HeaderContentType, w.FormDataContentType())

	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	h := &meals.Handler{MealRepository: repo}

	err = h.UploadMeals(c)

	require.NoError(t, err)

	require.Equal(t, http.StatusCreated, rec.Code)

	m, err := repo.Get()

	require.NoError(t, err)
	require.Len(t, m, 2)

	require.Equal(t, m[0].Name, "bar")
	require.Len(t, m[0].MealIngredients, 2)
	require.Equal(t, m[0].MealIngredients, []meals.MealIngredient{
		*meals.NewMealIngredient("def").WithQuantity(400, meals.Gram),
		*meals.NewMealIngredient("ghi").WithQuantity(6, meals.Tbsp),
	})

	require.Equal(t, m[1].Name, "foo")
	require.Len(t, m[1].MealIngredients, 2)
	require.Equal(t, m[1].MealIngredients, []meals.MealIngredient{
		*meals.NewMealIngredient("abc").WithQuantity(300, meals.Gram),
		*meals.NewMealIngredient("def").WithQuantity(5, meals.Tbsp),
	})
}
