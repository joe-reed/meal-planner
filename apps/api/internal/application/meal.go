package application

import (
	"errors"
	"github.com/joe-reed/meal-planner/apps/api/internal/domain/meals"
	"log/slog"
)

type MealApplication struct {
	r *meals.MealRepository
}

func NewMealApplication(r *meals.MealRepository) *MealApplication {
	return &MealApplication{r: r}
}

var MealAlreadyExists = errors.New("meal already exists")

func (a *MealApplication) AddMeal(id string, name string, mealIngredients []meals.MealIngredient) (*meals.Meal, error) {
	existingMeal, err := a.r.FindByName(name)

	if err != nil {
		return nil, err
	}

	if existingMeal != nil {
		return nil, MealAlreadyExists
	}

	m, err := meals.NewMeal(id, name, mealIngredients)
	if err != nil {
		return nil, err
	}

	slog.Debug("Adding meal", "meal", m)

	if err := a.r.Save(m); err != nil {
		return nil, err
	}

	return m, nil
}
