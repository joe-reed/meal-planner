package meal

import (
	"context"
	"database/sql"
	"github.com/hallgren/eventsourcing"
	"github.com/hallgren/eventsourcing/aggregate"
	"github.com/hallgren/eventsourcing/core"
	"github.com/hallgren/eventsourcing/eventstore/memory"
	sqlStore "github.com/hallgren/eventsourcing/eventstore/sql"
	_ "github.com/mattn/go-sqlite3"
	"sort"
)

type MealRepository interface {
	Get() ([]*Meal, error)
	Find(id string) (*Meal, error)
	Save(m *Meal) error
	FindByName(name string) (*Meal, error)
}

type EventSourcedMealRepository struct {
	es  core.EventStore
	all func() (core.Iterator, error)
}

func NewMealRepository(es core.EventStore, all func() (core.Iterator, error)) *EventSourcedMealRepository {
	aggregate.Register(&Meal{})
	r := &EventSourcedMealRepository{es, all}
	return r
}

func NewSqliteMealRepository(db *sql.DB) (*EventSourcedMealRepository, error) {
	es, err := sqlStore.NewSQLiteSingelWriter(db)

	if err != nil {
		return nil, err
	}

	return NewMealRepository(es, func() (core.Iterator, error) {
		return es.All(0, 10000)
	}), nil
}

func NewFakeMealRepository() *EventSourcedMealRepository {
	es := memory.Create()

	return NewMealRepository(es, func() (core.Iterator, error) {
		return es.All(0, 10000)()
	})
}

func (r EventSourcedMealRepository) Get() ([]*Meal, error) {
	mealMap := map[string]*Meal{}

	p := eventsourcing.NewProjection(
		r.all,
		func(e eventsourcing.Event) error {
			if e.AggregateType() != "Meal" {
				return nil
			}

			meal, ok := mealMap[e.AggregateID()]
			if !ok {
				meal = &Meal{}
				meal.Ingredients = []Ingredient{}
				mealMap[e.AggregateID()] = meal
			}

			meal.Transition(e)

			return nil
		})

	(*p).Strict = false
	_, result := p.RunOnce()

	if result.Error != nil {
		return nil, result.Error
	}

	meals := make([]*Meal, 0, len(mealMap))
	for _, m := range mealMap {
		meals = append(meals, m)
	}

	sort.Slice(meals, func(i, j int) bool {
		return meals[i].Name < meals[j].Name
	})

	return meals, nil
}

func (r EventSourcedMealRepository) Find(id string) (*Meal, error) {
	m := &Meal{}
	err := aggregate.Load(context.Background(), r.es, id, m)
	if err != nil {
		return nil, err
	}

	return m, nil
}

func (r EventSourcedMealRepository) Save(m *Meal) error {
	return aggregate.Save(r.es, m)
}

func (r EventSourcedMealRepository) FindByName(name string) (*Meal, error) {
	meals, err := r.Get()

	if err != nil {
		return nil, err
	}

	for _, m := range meals {
		if m.Name == name {
			return m, nil
		}
	}

	return nil, nil
}
