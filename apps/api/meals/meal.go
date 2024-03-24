package meals

import (
	"github.com/google/uuid"
)

type Meal struct {
	Id              string           `json:"id"`
	Name            string           `json:"name"`
	MealIngredients []MealIngredient `json:"ingredients"`
}

type MealIngredient struct {
	IngredientId string `json:"id"`
}

type MealBuilder struct {
	id              string
	name            string
	mealIngredients []MealIngredient
}

func (b *MealBuilder) WithName(name string) *MealBuilder {
	b.name = name
	return b
}

func (b *MealBuilder) AddIngredient(i MealIngredient) *MealBuilder {
	b.mealIngredients = append(b.mealIngredients, i)
	return b
}

func (b *MealBuilder) Build() *Meal {
	id := uuid.New().String()

	if b.id != "" {
		id = b.id
	}

	return &Meal{
		Id:              id,
		Name:            b.name,
		MealIngredients: b.mealIngredients,
	}
}

func NewMealBuilder() *MealBuilder {
	return &MealBuilder{"", "", []MealIngredient{}}
}
