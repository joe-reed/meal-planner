package meals

import (
	"database/sql"
	"sort"

	_ "github.com/mattn/go-sqlite3"
)

type MealRepository interface {
	Get() []*Meal
	Add(*Meal)
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

func (r SqliteMealRepository) Get() []*Meal {
	rows, _ := r.db.Query("SELECT * FROM meals ORDER BY name")

	meals := []*Meal{}
	for rows.Next() {
		var m Meal
		_ = rows.Scan(&m.Id, &m.Name)
		meals = append(meals, &m)
	}

	return meals
}

func (r SqliteMealRepository) Add(m *Meal) {
	r.db.Exec("INSERT INTO meals VALUES(?,?);", m.Id, m.Name)
}

type FakeMealRepository struct {
	meals map[string]*Meal
}

func NewFakeMealRepository() FakeMealRepository {
	return FakeMealRepository{meals: map[string]*Meal{}}
}

func (r FakeMealRepository) Get() []*Meal {
	keys := make([]string, 0, len(r.meals))
	for k := range r.meals {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	v := make([]*Meal, 0, len(r.meals))
	for _, k := range keys {
		v = append(v, r.meals[k])
	}
	return v
}

func (r FakeMealRepository) Add(m *Meal) {
	r.meals[m.Name] = m
}
