package handlers_test

import (
	"bytes"
	"github.com/joe-reed/meal-planner/apps/api/internal/application"
	"github.com/joe-reed/meal-planner/apps/api/internal/domain/meal"
	"github.com/joe-reed/meal-planner/apps/api/internal/domain/product"
	"github.com/joe-reed/meal-planner/apps/api/internal/handlers"
	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/require"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestUploadingMeals(t *testing.T) {
	repo := meal.NewFakeMealRepository()
	productRepo := product.NewFakeProductRepository()

	err := productRepo.Add(product.NewProductBuilder().WithName("Abc Name").WithId("abc").Build())
	require.NoError(t, err)

	err = productRepo.Add(product.NewProductBuilder().WithName("Def Name").WithId("def").Build())
	require.NoError(t, err)

	err = productRepo.Add(product.NewProductBuilder().WithName("Ghi Name").WithId("ghi").Build())
	require.NoError(t, err)

	e := echo.New()
	body := &bytes.Buffer{}
	w := multipart.NewWriter(body)

	part, err := w.CreateFormFile("meals", "meals.csv")
	require.NoError(t, err)

	_, err = part.Write([]byte("name,product,amount,unit\nfoo,Abc Name,300,Gram\nfoo,Def Name,5,Tbsp\nbar,Def Name,400,Gram\nbar,Ghi Name,6,Tbsp"))
	require.NoError(t, err)

	err = w.Close()
	require.NoError(t, err)

	req := httptest.NewRequest("POST", "/meals/upload", body)

	req.Header.Set(echo.HeaderContentType, w.FormDataContentType())

	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	h := &handlers.UploadHandler{Application: application.NewUploadMealsApplication(productRepo, repo)}

	err = h.UploadMeals(c)

	require.NoError(t, err)

	require.Equal(t, http.StatusCreated, rec.Code)

	m, err := repo.Get()

	require.NoError(t, err)
	require.Len(t, m, 2)

	require.Equal(t, m[0].Name, "bar")
	require.Len(t, m[0].Ingredients, 2)
	require.Equal(t, []meal.Ingredient{
		*meal.NewIngredient("def").WithQuantity(400, meal.Gram),
		*meal.NewIngredient("ghi").WithQuantity(6, meal.Tbsp),
	}, m[0].Ingredients)

	require.Equal(t, m[1].Name, "foo")
	require.Len(t, m[1].Ingredients, 2)
	require.Equal(t, []meal.Ingredient{
		*meal.NewIngredient("abc").WithQuantity(300, meal.Gram),
		*meal.NewIngredient("def").WithQuantity(5, meal.Tbsp),
	}, m[1].Ingredients)
}

func TestProductsNotExisting(t *testing.T) {
	repo := meal.NewFakeMealRepository()

	e := echo.New()
	body := &bytes.Buffer{}
	w := multipart.NewWriter(body)

	part, err := w.CreateFormFile("meals", "meals.csv")
	require.NoError(t, err)

	_, err = part.Write([]byte("name,product,amount,unit\nfoo,Abc Name,300,Gram\nfoo,Def Name,5,Tbsp\nbar,Def Name,400,Gram\nbar,Ghi Name,6,Tbsp"))
	require.NoError(t, err)

	err = w.Close()
	require.NoError(t, err)

	req := httptest.NewRequest("POST", "/meals/upload", body)

	req.Header.Set(echo.HeaderContentType, w.FormDataContentType())

	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	h := &handlers.UploadHandler{Application: application.NewUploadMealsApplication(product.NewFakeProductRepository(), repo)}

	err = h.UploadMeals(c)

	require.NoError(t, err)

	require.Equal(t, http.StatusBadRequest, rec.Code)

	m, err := repo.Get()

	require.Len(t, m, 0)

	require.Equal(t, "{\"notFoundIngredients\":[\"Abc Name\",\"Def Name\",\"Def Name\",\"Ghi Name\"]}\n", rec.Body.String())
}

func TestMealAlreadyExisting(t *testing.T) {
	productRepo := product.NewFakeProductRepository()

	err := productRepo.Add(product.NewProductBuilder().WithName("Abc Name").WithId("abc").Build())
	require.NoError(t, err)

	err = productRepo.Add(product.NewProductBuilder().WithName("Def Name").WithId("def").Build())
	require.NoError(t, err)

	err = productRepo.Add(product.NewProductBuilder().WithName("Ghi Name").WithId("ghi").Build())
	require.NoError(t, err)

	repo := meal.NewFakeMealRepository()

	e := echo.New()
	body := &bytes.Buffer{}
	w := multipart.NewWriter(body)

	part, err := w.CreateFormFile("meals", "meals.csv")
	require.NoError(t, err)

	err = repo.Save(
		meal.NewMealBuilder().WithName("bar").Build(),
	)
	require.NoError(t, err)

	_, err = part.Write([]byte("name,product,amount,unit\nfoo,Abc Name,300,Gram\nfoo,Def Name,5,Tbsp\nbar,Def Name,400,Gram\nbar,Ghi Name,6,Tbsp"))
	require.NoError(t, err)

	err = w.Close()
	require.NoError(t, err)

	req := httptest.NewRequest("POST", "/meals/upload", body)

	req.Header.Set(echo.HeaderContentType, w.FormDataContentType())

	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	h := &handlers.UploadHandler{Application: application.NewUploadMealsApplication(productRepo, repo)}

	err = h.UploadMeals(c)

	m, err := repo.Get()

	require.Equal(t, http.StatusBadRequest, rec.Code)

	require.Len(t, m, 1)

	require.Equal(t, "{\"error\":\"meal already exists\",\"mealName\":\"bar\"}\n", rec.Body.String())
}
