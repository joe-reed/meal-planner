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

type ValidationError struct {
	Field   string
	Message string
}

func (e *ValidationError) Error() string {
	return e.Message
}

func (a *IngredientApplication) AddIngredient(id string, name ingredient.IngredientName, category category.CategoryName) (*ingredient.Ingredient, error) {
	err := validateId(id)
	if err != nil {
		return nil, err
	}

	err = validateName(name)
	if err != nil {
		return nil, err
	}

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

// todo: reduce duplication, standardise validation or use library
func validateId(id string) error {
	if id == "" {
		return &ValidationError{
			Field:   "id",
			Message: "id cannot be empty",
		}
	}
	return nil
}

func validateName(name ingredient.IngredientName) error {
	if name == "" {
		return &ValidationError{
			Field:   "name",
			Message: "name cannot be empty",
		}
	}
	return nil
}

func (a *IngredientApplication) GetIngredients() ([]*ingredient.Ingredient, error) {
	ings, err := a.r.Get()

	if err != nil {
		return nil, err
	}

	return ings, nil
}
