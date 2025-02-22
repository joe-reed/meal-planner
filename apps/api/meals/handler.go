package meals

import (
	"bytes"
	"encoding/csv"
	"errors"
	"fmt"
	"github.com/labstack/echo/v4"
	"log/slog"
	"mime/multipart"
	"net/http"
	"strconv"

	"github.com/joe-reed/meal-planner/apps/api/ingredients"
)

type Handler struct {
	MealRepository *MealRepository

	// todo: refactor to reference exported application service not repository directly
	IngredientRepository *ingredients.IngredientRepository
}

func (h *Handler) GetMeals(c echo.Context) error {
	meals, err := h.MealRepository.Get()

	if err != nil {
		return err
	}

	return c.JSON(http.StatusOK, meals)
}

func (h *Handler) GetMeal(c echo.Context) error {
	meal, err := h.MealRepository.Find(c.Param("id"))

	if err != nil {
		return err
	}

	return c.JSON(http.StatusOK, meal)
}

func (h *Handler) AddMeal(c echo.Context) error {
	body := new(Meal)
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

	m, err := NewMeal(body.Id, body.Name, body.MealIngredients)
	if err != nil {
		return err
	}

	c.Logger().Debugf("Adding meal: %v", m)

	if err := h.MealRepository.Save(m); err != nil {
		return err
	}

	return c.JSON(http.StatusCreated, m)
}

func (h *Handler) AddIngredientToMeal(c echo.Context) error {
	mealId := c.Param("mealId")

	ingredient := NewMealIngredient("")
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

func (h *Handler) RemoveIngredientFromMeal(c echo.Context) error {
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

func (h *Handler) UploadMeals(c echo.Context) error {
	file, err := c.FormFile("meals")

	if err != nil {
		return err
	}

	src, err := file.Open()

	if err != nil {
		return err
	}
	defer src.Close()

	meals, err := ParseMeals(src, h.IngredientRepository)

	if err != nil {
		return err
	}

	slog.Info("uploading meals", "meals", meals)

	for _, m := range meals {
		if err := h.MealRepository.Save(m); err != nil {
			return err
		}
	}

	return c.NoContent(http.StatusCreated)
}

func ParseMeals(src multipart.File, ingredientRepository *ingredients.IngredientRepository) ([]*Meal, error) {
	var buf bytes.Buffer
	_, err := buf.ReadFrom(src)

	if err != nil {
		return nil, err
	}

	var meals []*Meal
	var meal *Meal

	csvReader := csv.NewReader(&buf)
	records, err := csvReader.ReadAll()

	if err != nil {
		return nil, err
	}

	for i, record := range records {
		if i == 0 {
			if record[0] != "name" || record[1] != "ingredient" || record[2] != "amount" || record[3] != "unit" {
				return nil, errors.New("invalid csv header")
			}

			continue
		}

		if len(record) != 4 {
			return nil, errors.New("invalid csv row")
		}

		mealName := record[0]

		if meal == nil || mealName != meal.Name {
			if meal != nil {
				meals = append(meals, meal)
			}

			meal = NewMealBuilder().WithName(mealName).Build()
		}

		ingredientName := record[1]
		amount, err := strconv.Atoi(record[2])

		if err != nil {
			return nil, err
		}

		unit, ok := UnitFromString(record[3])

		if !ok {
			return nil, fmt.Errorf("invalid unit: %s", record[3])
		}

		ingredient, err := ingredientRepository.GetByName(ingredientName)

		if err != nil {
			return nil, err
		}

		meal.AddIngredient(*NewMealIngredient(ingredient.Id).WithQuantity(amount, unit))
	}

	if meal != nil {
		meals = append(meals, meal)
	}

	return meals, nil
}
