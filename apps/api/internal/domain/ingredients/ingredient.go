package ingredients

import (
	"github.com/google/uuid"
	"github.com/hallgren/eventsourcing"
	"github.com/hallgren/eventsourcing/aggregate"
	"github.com/joe-reed/meal-planner/apps/api/internal/domain/categories"
)

type Ingredient struct {
	aggregate.Root
	Id       string              `json:"id"`
	Name     string              `json:"name"`
	Category categories.Category `json:"category"`
}

func (m *Ingredient) Transition(event eventsourcing.Event) {
	switch e := event.Data().(type) {
	case *Created:
		m.Id = e.Id
		m.Name = e.Name
		m.Category = e.Category
	}
}

func (m *Ingredient) Register(r aggregate.RegisterFunc) {
	r(&Created{})
}

func NewIngredient(id string, name string, category categories.Category) (*Ingredient, error) {
	i := &Ingredient{}
	err := i.SetID(id)
	if err != nil {
		return nil, err
	}
	aggregate.TrackChange(i, &Created{Id: id, Name: name, Category: category})

	return i, nil
}

type IngredientBuilder struct {
	id       string
	name     string
	category categories.Category
}

func (b *IngredientBuilder) WithName(name string) *IngredientBuilder {
	b.name = name
	return b
}

func (b *IngredientBuilder) WithId(id string) *IngredientBuilder {
	b.id = id
	return b
}

func (b *IngredientBuilder) WithCategory(category categories.Category) *IngredientBuilder {
	b.category = category
	return b
}

func (b *IngredientBuilder) Build() *Ingredient {
	id := uuid.New().String()

	if b.id != "" {
		id = b.id
	}

	i, err := NewIngredient(id, b.name, b.category)

	if err != nil {
		return nil
	}

	return i
}

func NewIngredientBuilder() *IngredientBuilder {
	return &IngredientBuilder{"", "", categories.Fruit}
}
