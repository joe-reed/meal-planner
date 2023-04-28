package meals

import (
	"sort"

	"github.com/google/uuid"
)

type (
	MealRepository interface {
		Get() []*Meal
		Add(*Meal)
	}

	Meal struct {
		Id   string `json:"id"`
		Name string `json:"name"`
	}
)

type FakeMealRepository struct {
	meals map[string]*Meal
}

func NewFakeMealRepository() FakeMealRepository {
	return FakeMealRepository{meals: map[string]*Meal{}}
}

func (r FakeMealRepository) Get() []*Meal {
	keys := make([]string, 0, len(r.meals))
	for k := range r.meals {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	v := make([]*Meal, 0, len(r.meals))
	for _, k := range keys {
		v = append(v, r.meals[k])
	}
	return v
}

func (r FakeMealRepository) Add(m *Meal) {
	r.meals[m.Name] = m
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
