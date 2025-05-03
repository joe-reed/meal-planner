package handlers

import (
	"errors"
	"github.com/joe-reed/meal-planner/apps/api/internal/application"
	"github.com/joe-reed/meal-planner/apps/api/internal/domain/meal"
	"github.com/labstack/echo/v4"
	"net/http"
)

type MealsHandler struct {
	Application *application.MealApplication
}

func (h *MealsHandler) GetMeals(c echo.Context) error {
	m, err := h.Application.GetMeals()

	if err != nil {
		return err
	}

	return c.JSON(http.StatusOK, m)
}

func (h *MealsHandler) FindMeal(c echo.Context) error {
	m, err := h.Application.FindMeal(c.Param("id"))

	if err != nil {
		return err
	}

	return c.JSON(http.StatusOK, m)
}

func (h *MealsHandler) AddMeal(c echo.Context) error {
	body := new(meal.Meal)
	if err := c.Bind(body); err != nil {
		return err
	}

	mealIngredients := body.MealIngredients
	if mealIngredients == nil {
		mealIngredients = make([]meal.MealIngredient, 0)
	}

	m, err := h.Application.AddMeal(
		body.Id,
		body.Name,
		body.Url,
		mealIngredients,
	)

	if err != nil {
		var mealAlreadyExists *application.MealAlreadyExists
		if errors.As(err, &mealAlreadyExists) {
			return c.JSON(http.StatusBadRequest, struct {
				Error    string `json:"error"`
				MealName string `json:"mealName"`
			}{mealAlreadyExists.Error(), mealAlreadyExists.MealName})
		}

		var validationError *application.ValidationError
		if errors.As(err, &validationError) {
			return c.JSON(http.StatusBadRequest, struct {
				Error string `json:"error"`
			}{
				Error: validationError.Error(),
			})
		}

		return err
	}

	return c.JSON(http.StatusCreated, m)
}

func (h *MealsHandler) AddIngredientToMeal(c echo.Context) error {
	mealId := c.Param("mealId")

	ingredient := meal.NewMealIngredient("")
	if err := c.Bind(&ingredient); err != nil {
		return err
	}

	m, err := h.Application.AddIngredientToMeal(mealId, *ingredient)
	if err != nil {
		return err
	}

	return c.JSON(http.StatusOK, m)
}

func (h *MealsHandler) RemoveIngredientFromMeal(c echo.Context) error {
	mealId := c.Param("mealId")
	ingredientId := c.Param("ingredientId")

	m, err := h.Application.RemoveIngredientFromMeal(mealId, ingredientId)
	if err != nil {
		return err
	}

	return c.JSON(http.StatusOK, m)
}

func (h *MealsHandler) UpdateMeal(c echo.Context) error {
	mealId := c.Param("mealId")

	body := new(application.PartialMeal)
	if err := c.Bind(body); err != nil {
		return err
	}

	m, err := h.Application.UpdateMeal(mealId, *body)
	if err != nil {
		return errors.New("error updating meal: " + err.Error())
	}

	return c.JSON(http.StatusOK, m)
}
