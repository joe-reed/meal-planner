package projections

import (
	"github.com/hallgren/eventsourcing"
	"github.com/hallgren/eventsourcing/core"
	sqlStore "github.com/hallgren/eventsourcing/eventstore/sql"
	"github.com/joe-reed/meal-planner/apps/api/internal/domain/category"
	"github.com/joe-reed/meal-planner/apps/api/internal/domain/product"
)

type ProductProjectionOutput map[category.CategoryName][]product.Product

func CreateProductProjection(es *sqlStore.SQL) (*eventsourcing.Projection, ProductProjectionOutput) {
	prods := ProductProjectionOutput{}

	start := core.Version(0)

	p := eventsourcing.NewProjection(func() (core.Iterator, error) {
		return es.All(start, 10)
	}, func(ev eventsourcing.Event) error {
		switch event := ev.Data().(type) {
		case *product.Created:
			prod := product.Product{}
			prod.Transition(ev)
			prods[event.Category] = append(prods[event.Category], prod)
		}

		start = core.Version(ev.GlobalVersion() + 1)

		return nil
	})

	return p, prods
}
