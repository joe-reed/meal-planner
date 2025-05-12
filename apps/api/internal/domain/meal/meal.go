package meal

import (
	"github.com/google/uuid"
	"github.com/hallgren/eventsourcing"
	"github.com/hallgren/eventsourcing/aggregate"
)

type Meal struct {
	aggregate.Root
	Id          string       `json:"id"`
	Name        string       `json:"name"`
	Url         string       `json:"url"`
	Ingredients []Ingredient `json:"ingredients"`
}

func (m *Meal) Transition(event eventsourcing.Event) {
	switch e := event.Data().(type) {
	case *Created:
		m.Id = e.Id
		m.Name = e.Name
		m.Url = e.Url
		m.Ingredients = e.Ingredients
	case *IngredientAdded:
		m.Ingredients = append(m.Ingredients, e.Ingredient)
	case *IngredientRemoved:
		ingredients := []Ingredient{}
		for _, ingredient := range m.Ingredients {
			if ingredient.IngredientId != e.Id {
				ingredients = append(ingredients, ingredient)
			}
		}
		m.Ingredients = ingredients
	case *NameUpdated:
		m.Name = e.Name
	case *UrlUpdated:
		m.Url = e.Url
	}
}

func (m *Meal) Register(r aggregate.RegisterFunc) {
	r(&Created{}, &IngredientAdded{}, &IngredientRemoved{}, &NameUpdated{}, &UrlUpdated{})
}

func NewMeal(id string, name string, url string, ingredients []Ingredient) (*Meal, error) {
	m := &Meal{}
	err := m.SetID(id)
	if err != nil {
		return nil, err
	}
	aggregate.TrackChange(m, &Created{Id: id, Name: name, Url: url, Ingredients: ingredients})

	return m, nil
}

func (m *Meal) AddIngredient(ingredient Ingredient) {
	aggregate.TrackChange(m, &IngredientAdded{Ingredient: ingredient})
}

func (m *Meal) RemoveIngredient(id string) {
	aggregate.TrackChange(m, &IngredientRemoved{Id: id})
}

func (m *Meal) UpdateName(name string) {
	aggregate.TrackChange(m, &NameUpdated{Name: name})
}

func (m *Meal) UpdateUrl(url string) {
	aggregate.TrackChange(m, &UrlUpdated{Url: url})
}

type Quantity struct {
	Amount int  `json:"amount"`
	Unit   Unit `json:"unit"`
}

type Ingredient struct {
	IngredientId string   `json:"id"`
	Quantity     Quantity `json:"quantity"`
}

func NewIngredient(id string) *Ingredient {
	return &Ingredient{IngredientId: id, Quantity: Quantity{1, Number}}
}

func (m *Ingredient) WithQuantity(amount int, unit Unit) *Ingredient {
	m.Quantity = Quantity{amount, unit}

	return m
}

type MealBuilder struct {
	id          string
	name        string
	url         string
	Ingredients []Ingredient
}

func NewMealBuilder() *MealBuilder {
	return &MealBuilder{"", "", "", []Ingredient{}}
}

func (b *MealBuilder) WithName(name string) *MealBuilder {
	b.name = name
	return b
}

func (b *MealBuilder) WithUrl(url string) *MealBuilder {
	b.url = url
	return b
}

func (b *MealBuilder) AddIngredient(i Ingredient) *MealBuilder {
	b.Ingredients = append(b.Ingredients, i)
	return b
}

func (b *MealBuilder) AddIngredients(i []Ingredient) *MealBuilder {
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

	meal, err := NewMeal(id, b.name, b.url, b.Ingredients)
	if err != nil {
		return nil
	}

	return meal
}

func (b *MealBuilder) WithId(i string) *MealBuilder {
	b.id = i

	return b
}
