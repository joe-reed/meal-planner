package meal

import (
	"github.com/google/uuid"
	"github.com/hallgren/eventsourcing"
	"github.com/hallgren/eventsourcing/aggregate"
)

type Meal struct {
	aggregate.Root
	Id              string           `json:"id"`
	Name            string           `json:"name"`
	MealIngredients []MealIngredient `json:"ingredients"`
}

func (m *Meal) Transition(event eventsourcing.Event) {
	switch e := event.Data().(type) {
	case *Created:
		m.Id = e.Id
		m.Name = e.Name
		m.MealIngredients = e.MealIngredients
	case *IngredientAdded:
		m.MealIngredients = append(m.MealIngredients, e.Ingredient)
	case *IngredientRemoved:
		ingredients := []MealIngredient{}
		for _, ingredient := range m.MealIngredients {
			if ingredient.IngredientId != e.Id {
				ingredients = append(ingredients, ingredient)
			}
		}
		m.MealIngredients = ingredients
	}
}

func (m *Meal) Register(r aggregate.RegisterFunc) {
	r(&Created{}, &IngredientAdded{}, &IngredientRemoved{})
}

func NewMeal(id string, name string, mealIngredients []MealIngredient) (*Meal, error) {
	m := &Meal{}
	err := m.SetID(id)
	if err != nil {
		return nil, err
	}
	aggregate.TrackChange(m, &Created{Id: id, Name: name, MealIngredients: mealIngredients})

	return m, nil
}

func (m *Meal) AddIngredient(ingredient MealIngredient) {
	aggregate.TrackChange(m, &IngredientAdded{Ingredient: ingredient})
}

func (m *Meal) RemoveIngredient(id string) {
	aggregate.TrackChange(m, &IngredientRemoved{Id: id})
}

type Quantity struct {
	Amount int  `json:"amount"`
	Unit   Unit `json:"unit"`
}

type MealIngredient struct {
	IngredientId string   `json:"id"`
	Quantity     Quantity `json:"quantity"`
}

func NewMealIngredient(id string) *MealIngredient {
	return &MealIngredient{IngredientId: id, Quantity: Quantity{1, Number}}
}

func (m *MealIngredient) WithQuantity(amount int, unit Unit) *MealIngredient {
	m.Quantity = Quantity{amount, unit}

	return m
}

type MealBuilder struct {
	id              string
	name            string
	mealIngredients []MealIngredient
}

func NewMealBuilder() *MealBuilder {
	return &MealBuilder{"", "", []MealIngredient{}}
}

func (b *MealBuilder) WithName(name string) *MealBuilder {
	b.name = name
	return b
}

func (b *MealBuilder) AddIngredient(i MealIngredient) *MealBuilder {
	b.mealIngredients = append(b.mealIngredients, i)
	return b
}

func (b *MealBuilder) AddIngredients(i []MealIngredient) *MealBuilder {
	for _, ingredient := range i {
		b.AddIngredient(ingredient)
	}
	return b
}

func (b *MealBuilder) Build() *Meal {
	id := uuid.New().String()

	if b.id != "" {
		id = b.id
	}

	meal, err := NewMeal(id, b.name, b.mealIngredients)
	if err != nil {
		return nil
	}

	return meal
}

func (b *MealBuilder) WithId(i string) *MealBuilder {
	b.id = i

	return b
}
