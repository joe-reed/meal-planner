package application

import (
	"errors"
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

type PartialMeal struct {
	Name *string `json:"name"`
	Url  *string `json:"url"`
}

func (a *MealApplication) AddMeal(id string, name string, url string, mealIngredients []meal.MealIngredient) (*meal.Meal, error) {
	err := validateId(id)
	if err != nil {
		return nil, err
	}

	err = validateName(name)
	if err != nil {
		return nil, err
	}

	existingMeal, err := a.r.FindByName(name)

	if err != nil {
		return nil, err
	}

	if existingMeal != nil {
		return nil, &MealAlreadyExists{
			MealName: name,
		}
	}

	m, err := meal.NewMeal(id, name, url, mealIngredients)
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
	m, err := a.r.Find(mealId)
	if err != nil {
		return nil, err
	}

	m.AddIngredient(mealIngredient)

	if err := a.r.Save(m); err != nil {
		return nil, err
	}

	return m, nil
}

func (a *MealApplication) RemoveIngredientFromMeal(mealId string, ingredientId string) (*meal.Meal, error) {
	m, err := a.r.Find(mealId)
	if err != nil {
		return nil, err
	}

	m.RemoveIngredient(ingredientId)

	if err := a.r.Save(m); err != nil {
		return nil, err
	}

	return m, nil
}

func (a *MealApplication) UpdateMeal(mealId string, body PartialMeal) (*meal.Meal, error) {
	m, err := a.r.Find(mealId)
	if err != nil {
		return nil, errors.New("error finding meal: " + err.Error())
	}

	if body.Name != nil {
		m.UpdateName(*body.Name)
	}

	if body.Url != nil {
		m.UpdateUrl(*body.Url)
	}

	if err := a.r.Save(m); err != nil {
		return nil, err
	}

	return m, nil
}
