package product

import (
	"errors"
	"github.com/google/uuid"
	"github.com/hallgren/eventsourcing"
	"github.com/hallgren/eventsourcing/aggregate"
	"github.com/joe-reed/meal-planner/apps/api/internal/domain/category"
)

type ProductName string

func (i ProductName) String() string {
	return string(i)
}

func NewProductName(name string) (ProductName, error) {
	if len(name) == 0 {
		return "", errors.New("name cannot be empty")
	}
	return ProductName(name), nil
}

type Product struct {
	aggregate.Root
	Id       string                `json:"id"`
	Name     ProductName           `json:"name"`
	Category category.CategoryName `json:"category"`
}

func (m *Product) Transition(event eventsourcing.Event) {
	switch e := event.Data().(type) {
	case *Created:
		m.Id = e.Id
		m.Name = ProductName(e.Name)
		m.Category = e.Category
	}
}

func (m *Product) Register(r aggregate.RegisterFunc) {
	r(&Created{})
}

func NewProduct(id string, name ProductName, category category.CategoryName) (*Product, error) {
	i := &Product{}
	err := i.SetID(id)
	if err != nil {
		return nil, err
	}
	aggregate.TrackChange(i, &Created{Id: id, Name: name.String(), Category: category})

	return i, nil
}

type ProductBuilder struct {
	id       string
	name     ProductName
	category category.CategoryName
}

func (b *ProductBuilder) WithName(name ProductName) *ProductBuilder {
	b.name = name
	return b
}

func (b *ProductBuilder) WithId(id string) *ProductBuilder {
	b.id = id
	return b
}

func (b *ProductBuilder) WithCategory(category category.CategoryName) *ProductBuilder {
	b.category = category
	return b
}

func (b *ProductBuilder) Build() *Product {
	id := uuid.New().String()

	if b.id != "" {
		id = b.id
	}

	i, err := NewProduct(id, b.name, b.category)

	if err != nil {
		return nil
	}

	return i
}

func NewProductBuilder() *ProductBuilder {
	return &ProductBuilder{"", "", category.Fruit}
}
