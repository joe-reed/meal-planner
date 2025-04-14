package application

import (
	"github.com/joe-reed/meal-planner/apps/api/internal/domain/category"
	"github.com/joe-reed/meal-planner/apps/api/internal/domain/ingredient"
	"log/slog"
)

type IngredientApplication struct {
	r *ingredient.IngredientRepository
}

func NewIngredientApplication(r *ingredient.IngredientRepository) *IngredientApplication {
	return &IngredientApplication{r: r}
}

func (a *IngredientApplication) AddIngredient(id string, name ingredient.IngredientName, category category.CategoryName) (*ingredient.Ingredient, error) {
	i, err := ingredient.NewIngredient(id, name, category)
	if err != nil {
		return nil, err
	}

	slog.Debug("Adding ingredient", "ingredient", i)

	if err := a.r.Add(i); err != nil {
		return i, err
	}

	return i, nil
}

func (a *IngredientApplication) GetIngredients() ([]*ingredient.Ingredient, error) {
	ings, err := a.r.Get()

	if err != nil {
		return nil, err
	}

	return ings, nil
}
