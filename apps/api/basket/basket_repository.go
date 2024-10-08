package basket

import (
	"database/sql"
	"github.com/hallgren/eventsourcing"
	"github.com/hallgren/eventsourcing/core"
	"github.com/hallgren/eventsourcing/eventstore/memory"
	sqlStore "github.com/hallgren/eventsourcing/eventstore/sql"
	_ "github.com/mattn/go-sqlite3"
	"strconv"
)

type BasketRepository struct {
	er  *eventsourcing.EventRepository
	all func() (core.Iterator, error)
}

func NewBasketRepository(es core.EventStore, all func() (core.Iterator, error)) *BasketRepository {
	er := eventsourcing.NewEventRepository(es)
	er.Register(&Basket{})
	r := &BasketRepository{er, all}
	return r
}

func NewSqliteBasketRepository(db *sql.DB) (*BasketRepository, error) {
	es := sqlStore.Open(db)

	return NewBasketRepository(es, func() (core.Iterator, error) {
		return es.All(0, 100000)
	}), nil
}

func NewFakeBasketRepository() *BasketRepository {
	es := memory.Create()

	return NewBasketRepository(es, es.All(0, 100000))
}

func (r BasketRepository) FindByShopId(shopId int) (*Basket, error) {
	b := &Basket{}
	err := r.er.Get(strconv.Itoa(shopId), b)
	if err != nil {
		return nil, err
	}

	return b, nil
}

func (r BasketRepository) Save(b *Basket) error {
	return r.er.Save(b)
}
