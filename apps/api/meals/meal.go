package meals

import (
	"github.com/google/uuid"
)

type Meal struct {
	Id   string `json:"id"`
	Name string `json:"name"`
}

type MealBuilder struct {
	id   string
	name string
}

func (b *MealBuilder) WithName(name string) *MealBuilder {
	b.name = name
	return b
}

func (b *MealBuilder) Build() *Meal {
	id := uuid.New().String()

	if b.id != "" {
		id = b.id
	}

	return &Meal{
		Id:   id,
		Name: b.name,
	}
}

func NewMealBuilder() *MealBuilder {
	return &MealBuilder{}
}
