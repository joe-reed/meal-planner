package basket

import (
	"context"
	"database/sql"
	"github.com/hallgren/eventsourcing/aggregate"
	"github.com/hallgren/eventsourcing/core"
	"github.com/hallgren/eventsourcing/eventstore/memory"
	sqlStore "github.com/hallgren/eventsourcing/eventstore/sql"
	_ "github.com/mattn/go-sqlite3"
	"strconv"
)

type BasketRepository struct {
	es  core.EventStore
	all func() (core.Iterator, error)
}

func NewBasketRepository(es core.EventStore, all func() (core.Iterator, error)) *BasketRepository {
	aggregate.Register(&Basket{})
	r := &BasketRepository{es, all}
	return r
}

func NewSqliteBasketRepository(db *sql.DB) (*BasketRepository, error) {
	es, err := sqlStore.NewSQLiteSingelWriter(db)

	if err != nil {
		return nil, err
	}

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
	err := aggregate.Load(context.Background(), r.es, strconv.Itoa(shopId), b)
	if err != nil {
		return nil, err
	}

	return b, nil
}

func (r BasketRepository) Save(b *Basket) error {
	return aggregate.Save(r.es, b)
}
