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
	FindByName(name ProductName) (*Product, error)
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
	productMap := map[string]*Product{}

	p := eventsourcing.NewProjection(
		r.all,
		func(e eventsourcing.Event) error {
			if e.AggregateType() != "Product" {
				return nil
			}

			p, ok := productMap[e.AggregateID()]
			if !ok {
				p = &Product{}
				productMap[e.AggregateID()] = p
			}

			p.Transition(e)

			return nil
		})

	(*p).Strict = false
	_, result := p.RunOnce()

	if result.Error != nil {
		return nil, result.Error
	}

	products := make([]*Product, 0, len(productMap))
	for _, in := range productMap {
		products = append(products, in)
	}

	sort.Slice(products, func(i, j int) bool {
		return products[i].Name < products[j].Name
	})

	return products, nil
}

func (r EventSourcedProductRepository) GetByName(name ProductName) (*Product, error) {
	p, err := r.FindByName(name)

	if err != nil {
		return nil, err
	}

	if p == nil {
		return nil, fmt.Errorf("product %s not found", name)
	}

	return p, nil
}

func (r EventSourcedProductRepository) FindByName(name ProductName) (*Product, error) {
	products, err := r.Get()

	if err != nil {
		return nil, err
	}

	for _, i := range products {
		if i.Name == name {
			return i, nil
		}
	}

	return nil, nil
}
