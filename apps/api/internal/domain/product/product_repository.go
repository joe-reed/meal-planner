package product

import (
	"database/sql"
	"fmt"
	"github.com/hallgren/eventsourcing"
	"github.com/hallgren/eventsourcing/aggregate"
	"github.com/hallgren/eventsourcing/core"
	"github.com/hallgren/eventsourcing/eventstore/memory"
	sqlStore "github.com/hallgren/eventsourcing/eventstore/sql"
	"sort"

	_ "github.com/mattn/go-sqlite3"
)

type ProductRepository interface {
	Add(i *Product) error
	Get() ([]*Product, error)
	GetByName(name ProductName) (*Product, error)
}

type EventSourcedProductRepository struct {
	es  core.EventStore
	all func() (core.Iterator, error)
}

func NewProductRepository(es core.EventStore, all func() (core.Iterator, error)) *EventSourcedProductRepository {
	aggregate.Register(&Product{})
	return &EventSourcedProductRepository{es, all}
}

func NewSqliteProductRepository(db *sql.DB) (*EventSourcedProductRepository, error) {
	es := sqlStore.Open(db)

	return NewProductRepository(es, func() (core.Iterator, error) {
		return es.All(0, 10000)
	}), nil
}

func NewFakeProductRepository() *EventSourcedProductRepository {
	es := memory.Create()

	return NewProductRepository(es, func() (core.Iterator, error) {
		return es.All(0, 10000)()
	})
}

func (r EventSourcedProductRepository) Add(i *Product) error {
	return aggregate.Save(r.es, i)
}

func (r EventSourcedProductRepository) Get() ([]*Product, error) {
	ingredientMap := map[string]*Product{}

	p := eventsourcing.NewProjection(
		r.all,
		func(e eventsourcing.Event) error {
			if e.AggregateType() != "Product" {
				return nil
			}

			ingredient, ok := ingredientMap[e.AggregateID()]
			if !ok {
				ingredient = &Product{}
				ingredientMap[e.AggregateID()] = ingredient
			}

			ingredient.Transition(e)

			return nil
		})

	(*p).Strict = false
	_, result := p.RunOnce()

	if result.Error != nil {
		return nil, result.Error
	}

	ingredients := make([]*Product, 0, len(ingredientMap))
	for _, in := range ingredientMap {
		ingredients = append(ingredients, in)
	}

	sort.Slice(ingredients, func(i, j int) bool {
		return ingredients[i].Name < ingredients[j].Name
	})

	return ingredients, nil
}

func (r EventSourcedProductRepository) GetByName(name ProductName) (*Product, error) {
	ingredients, err := r.Get()

	if err != nil {
		return nil, err
	}

	for _, i := range ingredients {
		if i.Name == name {
			return i, nil
		}
	}

	return nil, fmt.Errorf("ingredient %s not found", name)
}
