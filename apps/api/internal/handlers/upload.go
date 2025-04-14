package handlers

import (
	"errors"
	"github.com/joe-reed/meal-planner/apps/api/internal/application"
	"github.com/joe-reed/meal-planner/apps/api/internal/domain/ingredient"
	"github.com/labstack/echo/v4"
	"net/http"
)

type UploadHandler struct {
	Application *application.UploadMealsApplication
}

func (h *UploadHandler) UploadMeals(c echo.Context) error {
	file, err := c.FormFile("meals")

	if err != nil {
		return err
	}

	src, err := file.Open()

	defer src.Close()

	err = h.Application.UploadMeals(src)

	if err != nil {
		var mealAlreadyExists *application.MealAlreadyExists
		if errors.As(err, &mealAlreadyExists) {
			return c.JSON(http.StatusBadRequest, struct {
				Error    string `json:"error"`
				MealName string `json:"mealName"`
			}{mealAlreadyExists.Error(), mealAlreadyExists.MealName})
		}

		var ingredientsNotFound *application.IngredientsNotFound
		if errors.As(err, &ingredientsNotFound) {
			return c.JSON(http.StatusBadRequest, struct {
				NotFoundIngredients []ingredient.IngredientName `json:"notFoundIngredients"`
			}{ingredientsNotFound.NotFoundIngredients})
		}
	}

	return c.NoContent(http.StatusCreated)
}
