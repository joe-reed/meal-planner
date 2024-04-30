package ingredients

import (
	"github.com/google/uuid"
)

type Ingredient struct {
	Id   string `json:"id"`
	Name string `json:"name"`
}

type IngredientBuilder struct {
	id   string
	name string
}

func (b *IngredientBuilder) WithName(name string) *IngredientBuilder {
	b.name = name
	return b
}

func (b *IngredientBuilder) WithId(id string) *IngredientBuilder {
	b.id = id
	return b
}

func (b *IngredientBuilder) Build() *Ingredient {
	id := uuid.New().String()

	if b.id != "" {
		id = b.id
	}

	return &Ingredient{
		Id:   id,
		Name: b.name,
	}
}

func NewIngredientBuilder() *IngredientBuilder {
	return &IngredientBuilder{"", ""}
}
