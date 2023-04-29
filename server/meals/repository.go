package meals

import (
	"database/sql"
	"sort"

	_ "github.com/mattn/go-sqlite3"
)

type MealRepository interface {
	Get() ([]*Meal, error)
	Add(*Meal) error
}

type SqliteMealRepository struct {
	db *sql.DB
}

func NewSqliteMealRepository(dbFile string) (*SqliteMealRepository, error) {
	db, err := sql.Open("sqlite3", dbFile)
	if err != nil {
		return nil, err
	}
	if _, err := db.Exec("CREATE TABLE IF NOT EXISTS meals (id VARCHAR(255) NOT NULL PRIMARY KEY, name VARCHAR(255) NOT NULL);"); err != nil {
		return nil, err
	}

	return &SqliteMealRepository{db: db}, nil
}

func (r SqliteMealRepository) Get() ([]*Meal, error) {
	rows, err := r.db.Query("SELECT * FROM meals ORDER BY name")

	if err != nil {
		return nil, err
	}

	meals := []*Meal{}
	for rows.Next() {
		var m Meal
		err = rows.Scan(&m.Id, &m.Name)
		if err != nil {
			return nil, err
		}
		meals = append(meals, &m)
	}

	return meals, nil
}

func (r SqliteMealRepository) Add(m *Meal) error {
	_, err := r.db.Exec("INSERT INTO meals VALUES(?,?);", m.Id, m.Name)

	return err
}

type FakeMealRepository struct {
	meals map[string]*Meal
}

func NewFakeMealRepository() FakeMealRepository {
	return FakeMealRepository{meals: map[string]*Meal{}}
}

func (r FakeMealRepository) Get() ([]*Meal, error) {
	keys := make([]string, 0, len(r.meals))
	for k := range r.meals {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	v := make([]*Meal, 0, len(r.meals))
	for _, k := range keys {
		v = append(v, r.meals[k])
	}
	return v, nil
}

func (r FakeMealRepository) Add(m *Meal) error {
	r.meals[m.Name] = m

	return nil
}
