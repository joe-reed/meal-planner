package shops

import (
	"database/sql"
	"fmt"

	_ "github.com/mattn/go-sqlite3"
)

type ShopRepository interface {
	Current() (*Shop, error)
	Add(*Shop) error
}

type SqliteShopRepository struct {
	db *sql.DB
}

func NewSqliteShopRepository(dbFile string) (*SqliteShopRepository, error) {
	db, err := sql.Open("sqlite3", dbFile)
	if err != nil {
		return nil, err
	}
	if _, err := db.Exec("CREATE TABLE IF NOT EXISTS shops (id BIGINT NOT NULL PRIMARY KEY)"); err != nil {
		return nil, err
	}

	return &SqliteShopRepository{db: db}, nil
}

func (r SqliteShopRepository) Current() (*Shop, error) {
	row := r.db.QueryRow("SELECT * FROM shops ORDER BY id DESC LIMIT 1")

	var s Shop
	err := row.Scan(&s.Id)
	if err != nil {
		return nil, err
	}

	return &s, nil
}

func (r SqliteShopRepository) Add(s *Shop) error {
	_, err := r.db.Exec("INSERT INTO shops (id) VALUES (?)", s.Id)

	return err
}

type FakeShopRepository struct {
	shops []*Shop
}

func NewFakeShopRepository() *FakeShopRepository {
	return &FakeShopRepository{shops: []*Shop{}}
}

func (r *FakeShopRepository) Current() (*Shop, error) {
	if len(r.shops) > 0 {
		return r.shops[0], nil
	}

	return nil, fmt.Errorf("Shop not found")
}

func (r *FakeShopRepository) Add(s *Shop) error {
	r.shops = append([]*Shop{s}, r.shops...)

	return nil
}
