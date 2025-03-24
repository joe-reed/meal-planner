package handlers

import (
	"bytes"
	"encoding/csv"
	"errors"
	"fmt"
	"github.com/joe-reed/meal-planner/apps/api/internal/domain/meals"
	"github.com/labstack/echo/v4"
	"log/slog"
	"mime/multipart"
	"net/http"
	"strconv"

	"github.com/joe-reed/meal-planner/apps/api/internal/domain/ingredients"
)

type MealsHandler struct {
	MealRepository *meals.MealRepository

	// todo: refactor to reference exported application service not repository directly
	IngredientRepository *ingredients.IngredientRepository
}

func (h *MealsHandler) GetMeals(c echo.Context) error {
	meals, err := h.MealRepository.Get()

	if err != nil {
		return err
	}

	return c.JSON(http.StatusOK, meals)
}

func (h *MealsHandler) GetMeal(c echo.Context) error {
	meal, err := h.MealRepository.Find(c.Param("id"))

	if err != nil {
		return err
	}

	return c.JSON(http.StatusOK, meal)
}

func (h *MealsHandler) AddMeal(c echo.Context) error {
	body := new(meals.Meal)
	if err := c.Bind(body); err != nil {
		return err
	}

	existingMeal, err := h.MealRepository.FindByName(body.Name)

	if err != nil {
		return err
	}

	if existingMeal != nil {
		return c.String(http.StatusBadRequest, "meal already exists")
	}

	mealIngredients := body.MealIngredients
	if mealIngredients == nil {
		mealIngredients = make([]meals.MealIngredient, 0)
	}

	m, err := meals.NewMeal(body.Id, body.Name, mealIngredients)
	if err != nil {
		return err
	}

	c.Logger().Debugf("Adding meal: %v", m)

	if err := h.MealRepository.Save(m); err != nil {
		return err
	}

	return c.JSON(http.StatusCreated, m)
}

func (h *MealsHandler) AddIngredientToMeal(c echo.Context) error {
	mealId := c.Param("mealId")

	ingredient := meals.NewMealIngredient("")
	if err := c.Bind(&ingredient); err != nil {
		return err
	}

	meal, err := h.MealRepository.Find(mealId)
	if err != nil {
		return err
	}

	meal.AddIngredient(*ingredient)

	if err := h.MealRepository.Save(meal); err != nil {
		return err
	}

	return c.JSON(http.StatusOK, meal)
}

func (h *MealsHandler) RemoveIngredientFromMeal(c echo.Context) error {
	mealId := c.Param("mealId")
	ingredientId := c.Param("ingredientId")

	meal, err := h.MealRepository.Find(mealId)
	if err != nil {
		return err
	}

	meal.RemoveIngredient(ingredientId)

	if err := h.MealRepository.Save(meal); err != nil {
		return err
	}

	return c.JSON(http.StatusOK, meal)
}

func (h *MealsHandler) UploadMeals(c echo.Context) error {
	file, err := c.FormFile("meals")

	if err != nil {
		return err
	}

	src, err := file.Open()

	if err != nil {
		return err
	}
	defer src.Close()

	meals, notFoundIngredients, err := ParseMeals(src, h.IngredientRepository)

	if err != nil {
		return err
	}

	if len(notFoundIngredients) > 0 {
		return c.JSON(http.StatusBadRequest, struct {
			NotFoundIngredients []string `json:"notFoundIngredients"`
		}{notFoundIngredients})
	}

	slog.Info("uploading meals", "meals", meals)

	for _, m := range meals {
		if err := h.MealRepository.Save(m); err != nil {
			return err
		}
	}

	return c.NoContent(http.StatusCreated)
}

func ParseMeals(src multipart.File, ingredientRepository *ingredients.IngredientRepository) (m []*meals.Meal, notFoundIngredients []string, err error) {
	var buf bytes.Buffer
	_, err = buf.ReadFrom(src)

	if err != nil {
		return nil, nil, err
	}

	var meal *meals.Meal

	csvReader := csv.NewReader(&buf)
	records, err := csvReader.ReadAll()

	if err != nil {
		return nil, nil, err
	}

	for i, record := range records {
		if i == 0 {
			if record[0] != "name" || record[1] != "ingredient" || record[2] != "amount" || record[3] != "unit" {
				return nil, nil, errors.New("invalid csv header")
			}

			continue
		}

		if len(record) != 4 {
			return nil, nil, errors.New("invalid csv row")
		}

		mealName := record[0]

		if meal == nil || mealName != meal.Name {
			if meal != nil {
				m = append(m, meal)
			}

			meal = meals.NewMealBuilder().WithName(mealName).Build()
		}

		ingredientName := record[1]
		amount, err := strconv.Atoi(record[2])

		if err != nil {
			return nil, nil, err
		}

		unit, ok := meals.UnitFromString(record[3])

		if !ok {
			return nil, nil, fmt.Errorf("invalid unit: %s", record[3])
		}

		ingredient, err := ingredientRepository.GetByName(ingredientName)

		if err != nil {
			notFoundIngredients = append(notFoundIngredients, ingredientName)
			continue
		}

		meal.AddIngredient(*meals.NewMealIngredient(ingredient.Id).WithQuantity(amount, unit))
	}

	if meal != nil {
		m = append(m, meal)
	}

	return m, notFoundIngredients, nil
}
