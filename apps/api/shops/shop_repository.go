package shops

import (
	"database/sql"
	"github.com/hallgren/eventsourcing"
	"github.com/hallgren/eventsourcing/core"
	"github.com/hallgren/eventsourcing/eventstore/memory"
	sqlStore "github.com/hallgren/eventsourcing/eventstore/sql"
	_ "github.com/mattn/go-sqlite3"
	"strconv"
)

type ShopRepository struct {
	er  *eventsourcing.EventRepository
	all func() (core.Iterator, error)
}

func NewShopRepository(es core.EventStore, all func() (core.Iterator, error)) *ShopRepository {
	er := eventsourcing.NewEventRepository(es)
	er.Register(&Shop{})
	r := &ShopRepository{er, all}
	return r
}

func NewSqliteShopRepository(db *sql.DB) (*ShopRepository, error) {
	es := sqlStore.Open(db)

	return NewShopRepository(es, func() (core.Iterator, error) {
		return es.All(0, 100000)
	}), nil
}

func NewFakeShopRepository() *ShopRepository {
	es := memory.Create()

	return NewShopRepository(es, es.All(0, 100000))
}

func (r ShopRepository) Current() (*Shop, error) {
	currentId := 0

	p := r.er.Projections.Projection(
		r.all,
		func(e eventsourcing.Event) error {
			aId, err := strconv.Atoi(e.AggregateID())
			if err != nil {
				return err
			}

			if aId > currentId {
				currentId = aId
			}

			return nil
		})

	(*p).Strict = false
	_, result := p.RunOnce()

	if result.Error != nil {
		return nil, result.Error
	}

	if currentId == 0 {
		return nil, nil
	}

	return r.Find(currentId)
}

func (r ShopRepository) Find(id int) (*Shop, error) {
	s := &Shop{}
	err := r.er.Get(strconv.Itoa(id), s)
	if err != nil {
		return nil, err
	}

	return s, nil
}

func (r ShopRepository) Save(s *Shop) error {
	return r.er.Save(s)
}

func cloneShop(s *Shop) (*Shop, error) {
	result, err := NewShop(s.Id)
	if err != nil {
		return nil, err
	}

	for _, m := range s.Meals {
		result.AddMeal(m)
	}

	return result, nil
}
