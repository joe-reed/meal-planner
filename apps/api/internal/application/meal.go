package application

import (
	"github.com/joe-reed/meal-planner/apps/api/internal/domain/meal"
	"log/slog"
)

type MealApplication struct {
	r *meal.MealRepository
}

func NewMealApplication(r *meal.MealRepository) *MealApplication {
	return &MealApplication{r: r}
}

type MealAlreadyExists struct {
	MealName string
}

func (*MealAlreadyExists) Error() string {
	return "meal already exists"
}

func (a *MealApplication) AddMeal(id string, name string, mealIngredients []meal.MealIngredient) (*meal.Meal, error) {
	existingMeal, err := a.r.FindByName(name)

	if err != nil {
		return nil, err
	}

	if existingMeal != nil {
		return nil, &MealAlreadyExists{
			MealName: name,
		}
	}

	m, err := meal.NewMeal(id, name, mealIngredients)
	if err != nil {
		return nil, err
	}

	slog.Debug("Adding meal", "meal", m)

	if err := a.r.Save(m); err != nil {
		return nil, err
	}

	return m, nil
}

func (a *MealApplication) BulkAddMeals(meals []*meal.Meal) ([]*meal.Meal, error) {
	for _, m := range meals {
		if err := a.r.Save(m); err != nil {
			return nil, err
		}
	}

	return meals, nil
}

func (a *MealApplication) GetMeals() ([]*meal.Meal, error) {
	m, err := a.r.Get()

	if err != nil {
		return nil, err
	}

	return m, nil
}

func (a *MealApplication) FindMeal(id string) (*meal.Meal, error) {
	m, err := a.r.Find(id)

	if err != nil {
		return nil, err
	}

	return m, nil
}

func (a *MealApplication) AddIngredientToMeal(mealId string, mealIngredient meal.MealIngredient) (*meal.Meal, error) {
	meal, err := a.r.Find(mealId)
	if err != nil {
		return nil, err
	}

	meal.AddIngredient(mealIngredient)

	if err := a.r.Save(meal); err != nil {
		return nil, err
	}

	return meal, nil
}

func (a *MealApplication) RemoveIngredientFromMeal(mealId string, ingredientId string) (*meal.Meal, error) {
	meal, err := a.r.Find(mealId)
	if err != nil {
		return nil, err
	}

	meal.RemoveIngredient(ingredientId)

	if err := a.r.Save(meal); err != nil {
		return nil, err
	}

	return meal, nil
}
