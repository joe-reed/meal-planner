package ingredients

import (
  "github.com/google/uuid"
  "github.com/hallgren/eventsourcing"
)

type Ingredient struct {
  eventsourcing.AggregateRoot
  Id   string `json:"id"`
  Name string `json:"name"`
}

func (m *Ingredient) Transition(event eventsourcing.Event) {
  switch e := event.Data().(type) {
  case *Created:
    m.Id = e.Id
    m.Name = e.Name
  }
}

func (m *Ingredient) Register(r eventsourcing.RegisterFunc) {
  r(&Created{})
}

func NewIngredient(id string, name string) (*Ingredient, error) {
  i := &Ingredient{}
  err := i.SetID(id)
  if err != nil {
    return nil, err
  }
  i.TrackChange(i, &Created{Id: id, Name: name})

  return i, nil
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

  i, err := NewIngredient(id, b.name)

  if err != nil {
    return nil
  }

  return i
}

func NewIngredientBuilder() *IngredientBuilder {
  return &IngredientBuilder{"", ""}
}
