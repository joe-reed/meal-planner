package shops

import (
	"database/sql"
	"errors"

	_ "github.com/mattn/go-sqlite3"
)

type ShopRepository interface {
	Current() (*Shop, error)
	Find(id int) (*Shop, error)
	Add(*Shop) error
	Save(*Shop) error
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
	if _, err := db.Exec("CREATE TABLE IF NOT EXISTS shop_meals (shop_id BIGINT NOT NULL, meal_id VARCHAR(255) NOT NULL, PRIMARY KEY (shop_id, meal_id))"); err != nil {
		return nil, err
	}

	return &SqliteShopRepository{db: db}, nil
}

func (r SqliteShopRepository) Current() (*Shop, error) {
	row := r.db.QueryRow("SELECT * FROM shops ORDER BY id DESC LIMIT 1")

	var s Shop
	err := row.Scan(&s.Id)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, nil
	}

	if err != nil {
		return nil, err
	}

	return &s, nil
}

func (r SqliteShopRepository) Find(id int) (*Shop, error) {
	rows, err := r.db.Query("SELECT meal_id FROM shops LEFT JOIN shop_meals ON shops.id = shop_meals.shop_id WHERE id = ?", id)

	if err != nil {
		return nil, err
	}

	shop := NewShop(id)

	for rows.Next() {
		var m sql.NullString

		err = rows.Scan(&m)

		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}

		if err != nil {
			return nil, err
		}

		if m.Valid {
			shop.Meals = append(shop.Meals, &ShopMeal{MealId: m.String})
		}
	}

	return shop, nil
}

func (r SqliteShopRepository) Add(s *Shop) error {
	_, err := r.db.Exec("INSERT INTO shops (id) VALUES (?)", s.Id)

	if err != nil {
		return err
	}

	for _, m := range s.Meals {
		_, err := r.db.Exec("INSERT INTO shop_meals (shop_id, meal_id) VALUES (?, ?)", s.Id, m.MealId)

		if err != nil {
			return err
		}
	}

	return nil
}

func (r SqliteShopRepository) Save(s *Shop) error {
	_, err := r.db.Exec("DELETE FROM shop_meals WHERE shop_id = ?", s.Id)

	if err != nil {
		return err
	}

	for _, m := range s.Meals {
		_, err := r.db.Exec("INSERT INTO shop_meals (shop_id, meal_id) VALUES (?, ?)", s.Id, m.MealId)

		if err != nil {
			return err
		}
	}

	return nil
}

type FakeShopRepository struct {
	shops []*Shop
}

func NewFakeShopRepository() *FakeShopRepository {
	return &FakeShopRepository{shops: []*Shop{}}
}

func (r *FakeShopRepository) Current() (*Shop, error) {
	if len(r.shops) > 0 {
		s := r.shops[0]

		return NewShop(s.Id).SetMeals(s.Meals), nil
	}

	return nil, nil
}

func (r *FakeShopRepository) Find(id int) (*Shop, error) {
	if len(r.shops) > 0 {
		for _, s := range r.shops {
			if s.Id == id {
				return NewShop(s.Id).SetMeals(s.Meals), nil
			}
		}
	}

	return nil, nil
}

func (r *FakeShopRepository) Add(s *Shop) error {
	r.shops = append([]*Shop{cloneShop(s)}, r.shops...)

	return nil
}

func (r *FakeShopRepository) Save(s *Shop) error {
	var shops []*Shop

	for _, t := range r.shops {
		if t.Id == s.Id {
			shops = append(shops, cloneShop(s))
		} else {
			shops = append(shops, t)
		}
	}

	r.shops = shops

	return nil
}

func cloneShop(s *Shop) *Shop {
	return &Shop{s.Id, s.Meals}
}
