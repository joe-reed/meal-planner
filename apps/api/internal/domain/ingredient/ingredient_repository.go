package ingredient

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

type IngredientRepository struct {
	es  core.EventStore
	all func() (core.Iterator, error)
}

func NewIngredientRepository(es core.EventStore, all func() (core.Iterator, error)) *IngredientRepository {
	aggregate.Register(&Ingredient{})
	return &IngredientRepository{es, all}
}

func NewSqliteIngredientRepository(db *sql.DB) (*IngredientRepository, error) {
	es := sqlStore.Open(db)

	return NewIngredientRepository(es, func() (core.Iterator, error) {
		return es.All(0, 10000)
	}), nil
}

func NewFakeIngredientRepository() *IngredientRepository {
	es := memory.Create()

	return NewIngredientRepository(es, func() (core.Iterator, error) {
		return es.All(0, 10000)()
	})
}

func (r IngredientRepository) Add(i *Ingredient) error {
	return aggregate.Save(r.es, i)
}

func (r IngredientRepository) Get() ([]*Ingredient, error) {
	ingredientMap := map[string]*Ingredient{}

	p := eventsourcing.NewProjection(
		r.all,
		func(e eventsourcing.Event) error {
			if e.AggregateType() != "Ingredient" {
				return nil
			}

			ingredient, ok := ingredientMap[e.AggregateID()]
			if !ok {
				ingredient = &Ingredient{}
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

	ingredients := make([]*Ingredient, 0, len(ingredientMap))
	for _, in := range ingredientMap {
		ingredients = append(ingredients, in)
	}

	sort.Slice(ingredients, func(i, j int) bool {
		return ingredients[i].Name < ingredients[j].Name
	})

	return ingredients, nil
}

func (r IngredientRepository) GetByName(name IngredientName) (*Ingredient, error) {
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
