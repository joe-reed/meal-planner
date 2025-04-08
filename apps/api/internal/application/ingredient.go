package application

import (
	"github.com/joe-reed/meal-planner/apps/api/internal/domain/categories"
	"github.com/joe-reed/meal-planner/apps/api/internal/domain/ingredients"
	"log/slog"
)

type IngredientApplication struct {
	r *ingredients.IngredientRepository
}

func NewIngredientApplication(r *ingredients.IngredientRepository) *IngredientApplication {
	return &IngredientApplication{r: r}
}

func (a *IngredientApplication) AddIngredient(id string, name ingredients.IngredientName, category categories.Category) (*ingredients.Ingredient, error) {
	i, err := ingredients.NewIngredient(id, name, category)
	if err != nil {
		return nil, err
	}

	slog.Debug("Adding ingredient", "ingredient", i)

	if err := a.r.Add(i); err != nil {
		return i, err
	}

	return i, nil
}

func (a *IngredientApplication) GetIngredients() ([]*ingredients.Ingredient, error) {
	ings, err := a.r.Get()

	if err != nil {
		return nil, err
	}

	return ings, nil
}
