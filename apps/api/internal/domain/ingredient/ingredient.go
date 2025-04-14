package ingredient

import (
	"errors"
	"github.com/google/uuid"
	"github.com/hallgren/eventsourcing"
	"github.com/hallgren/eventsourcing/aggregate"
	"github.com/joe-reed/meal-planner/apps/api/internal/domain/category"
)

type IngredientName string

func (i IngredientName) String() string {
	return string(i)
}

func NewIngredientName(name string) (IngredientName, error) {
	if len(name) == 0 {
		return "", errors.New("name cannot be empty")
	}
	return IngredientName(name), nil
}

type Ingredient struct {
	aggregate.Root
	Id       string                `json:"id"`
	Name     IngredientName        `json:"name"`
	Category category.CategoryName `json:"category"`
}

func (m *Ingredient) Transition(event eventsourcing.Event) {
	switch e := event.Data().(type) {
	case *Created:
		m.Id = e.Id
		m.Name = IngredientName(e.Name)
		m.Category = e.Category
	}
}

func (m *Ingredient) Register(r aggregate.RegisterFunc) {
	r(&Created{})
}

func NewIngredient(id string, name IngredientName, category category.CategoryName) (*Ingredient, error) {
	i := &Ingredient{}
	err := i.SetID(id)
	if err != nil {
		return nil, err
	}
	aggregate.TrackChange(i, &Created{Id: id, Name: name.String(), Category: category})

	return i, nil
}

type IngredientBuilder struct {
	id       string
	name     IngredientName
	category category.CategoryName
}

func (b *IngredientBuilder) WithName(name IngredientName) *IngredientBuilder {
	b.name = name
	return b
}

func (b *IngredientBuilder) WithId(id string) *IngredientBuilder {
	b.id = id
	return b
}

func (b *IngredientBuilder) WithCategory(category category.CategoryName) *IngredientBuilder {
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
	return &IngredientBuilder{"", "", category.Fruit}
}
