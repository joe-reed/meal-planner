package application

import (
	"bytes"
	"encoding/csv"
	"errors"
	"fmt"
	"github.com/joe-reed/meal-planner/apps/api/internal/domain/ingredients"
	"github.com/joe-reed/meal-planner/apps/api/internal/domain/meals"
	"io"
	"log/slog"
	"strconv"
)

type UploadMealsApplication struct {
	IngredientRepository *ingredients.IngredientRepository
	MealRepository       *meals.MealRepository
}

func NewUploadMealsApplication(ingredientRepository *ingredients.IngredientRepository, mealRepository *meals.MealRepository) *UploadMealsApplication {
	return &UploadMealsApplication{
		IngredientRepository: ingredientRepository,
		MealRepository:       mealRepository,
	}
}

type IngredientsNotFound struct {
	NotFoundIngredients []ingredients.IngredientName
}

func (*IngredientsNotFound) Error() string {
	return "ingredients not found"
}

func (a *UploadMealsApplication) UploadMeals(src io.Reader) error {
	ms, notFoundIngredients, err := a.parseMeals(src)

	if err != nil {
		return err
	}

	if len(notFoundIngredients) > 0 {
		return &IngredientsNotFound{
			NotFoundIngredients: notFoundIngredients,
		}
	}

	for _, m := range ms {
		meal, _ := a.MealRepository.FindByName(m.Name)

		if meal != nil {
			return &MealAlreadyExists{
				MealName: m.Name,
			}
		}
	}

	slog.Info("uploading meals", "meals", ms)

	for _, m := range ms {
		if err := a.MealRepository.Save(m); err != nil {
			return err
		}
	}

	return nil
}

func (a *UploadMealsApplication) parseMeals(src io.Reader) (m []*meals.Meal, notFoundIngredients []ingredients.IngredientName, err error) {
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

		ingredientName, err := ingredients.NewIngredientName(record[1])

		if err != nil {
			return nil, nil, err
		}

		amount, err := strconv.Atoi(record[2])

		if err != nil {
			return nil, nil, err
		}

		unit, ok := meals.UnitFromString(record[3])

		if !ok {
			return nil, nil, fmt.Errorf("invalid unit: %s", record[3])
		}

		ingredient, err := a.IngredientRepository.GetByName(ingredientName)

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
