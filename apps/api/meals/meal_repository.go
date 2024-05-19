package meals

import (
	"database/sql"
	"github.com/hallgren/eventsourcing"
	"github.com/hallgren/eventsourcing/core"
	"github.com/hallgren/eventsourcing/eventstore/memory"
	sqlStore "github.com/hallgren/eventsourcing/eventstore/sql"
	_ "github.com/mattn/go-sqlite3"
	"sort"
)

type MealRepository struct {
	er  *eventsourcing.EventRepository
	all func() (core.Iterator, error)
}

func NewMealRepository(es core.EventStore, all func() (core.Iterator, error)) *MealRepository {
	er := eventsourcing.NewEventRepository(es)
	er.Register(&Meal{})
	r := &MealRepository{er, all}
	return r
}

func NewSqliteMealRepository(dbFile string) (*MealRepository, error) {
	db, err := sql.Open("sqlite3", dbFile)
	if err != nil {
		return nil, err
	}

	es := sqlStore.Open(db)

	rows, err := db.Query("SELECT name FROM sqlite_master WHERE type='table' AND name='events'")
	if err != nil {
		return nil, err
	}

	if !rows.Next() {
		err = es.Migrate()
		if err != nil {
			return nil, err
		}
	}

	return NewMealRepository(es, func() (core.Iterator, error) {
		return es.All(0, 100)
	}), nil
}

func NewFakeMealRepository() *MealRepository {
	es := memory.Create()

	return NewMealRepository(es, es.All(0, 100))
}

func (r MealRepository) Get() ([]*Meal, error) {
	meals := map[string]*Meal{}

	p := r.er.Projections.Projection(
		r.all,
		func(e eventsourcing.Event) error {
			meal, ok := meals[e.AggregateID()]
			if !ok {
				meal = &Meal{}
				meal.MealIngredients = []MealIngredient{}
				meals[e.AggregateID()] = meal
			}
			meal.Transition(e)

			return nil
		})

	p.RunOnce()

	result := make([]*Meal, 0, len(meals))
	for _, m := range meals {
		result = append(result, m)
	}

	sort.Slice(result, func(i, j int) bool {
		return result[i].Name < result[j].Name
	})

	return result, nil
}

func (r MealRepository) Find(id string) (*Meal, error) {
	m := &Meal{}
	err := r.er.Get(id, m)
	if err != nil {
		return nil, err
	}

	return m, nil
}

func (r MealRepository) Save(m *Meal) error {
	return r.er.Save(m)
}
