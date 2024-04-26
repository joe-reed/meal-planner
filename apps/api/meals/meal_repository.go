package meals

import (
	"database/sql"
	"fmt"
	"sort"

	_ "github.com/mattn/go-sqlite3"
)

type MealRepository interface {
	Get() ([]*Meal, error)
	Find(id string) (*Meal, error)
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
	if _, err := db.Exec("CREATE TABLE IF NOT EXISTS meals (id VARCHAR(255) NOT NULL PRIMARY KEY, name VARCHAR(255) NOT NULL)"); err != nil {
		return nil, err
	}
	if _, err := db.Exec("CREATE TABLE IF NOT EXISTS meal_ingredients (meal_id VARCHAR(255) NOT NULL, ingredient_id VARCHAR(255) NOT NULL, PRIMARY KEY (meal_id, ingredient_id))"); err != nil {
		return nil, err
	}

	return &SqliteMealRepository{db: db}, nil
}

func (r SqliteMealRepository) Get() ([]*Meal, error) {
	rows, err := r.db.Query("SELECT id, name, ingredient_id FROM meals LEFT JOIN meal_ingredients ON meals.id = meal_ingredients.meal_id ORDER BY name")

	if err != nil {
		return nil, err
	}

	meals := map[string]*Meal{}

	for rows.Next() {
		var m Meal
		var i sql.NullString

		err = rows.Scan(&m.Id, &m.Name, &i)

		if err != nil {
			return nil, err
		}

		_, exists := meals[m.Id]
		if exists {
			m = *meals[m.Id]
		} else {
			m.MealIngredients = []MealIngredient{}
		}

		if i.Valid {
			m.MealIngredients = append(m.MealIngredients, MealIngredient{IngredientId: i.String})
		}

		meals[m.Id] = &m
	}

	result := make([]*Meal, 0, len(meals))

	for _, m := range meals {
		result = append(result, m)
	}

	return result, nil
}

func (r SqliteMealRepository) Find(id string) (*Meal, error) {
	rows, err := r.db.Query("SELECT id, name, ingredient_id FROM meals LEFT JOIN meal_ingredients ON meals.id = meal_ingredients.meal_id WHERE id = ?", id)

	if err != nil {
		return nil, err
	}

	var m Meal
	m.MealIngredients = []MealIngredient{}

	for rows.Next() {
		var i sql.NullString

		err = rows.Scan(&m.Id, &m.Name, &i)

		if err != nil {
			return nil, err
		}

		if i.Valid {
			m.MealIngredients = append(m.MealIngredients, MealIngredient{IngredientId: i.String})
		}
	}

	return &m, nil
}

func (r SqliteMealRepository) Add(m *Meal) error {
	_, err := r.db.Exec("INSERT INTO meals (id, name) VALUES (?, ?)", m.Id, m.Name)

	if err != nil {
		return err
	}

	for _, i := range m.MealIngredients {
		_, err := r.db.Exec("INSERT INTO meal_ingredients (meal_id, ingredient_id) VALUES (?, ?)", m.Id, i.IngredientId)

		if err != nil {
			return err
		}
	}

	return nil
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

func (r FakeMealRepository) Find(id string) (*Meal, error) {
	for _, m := range r.meals {
		if m.Id == id {
			return m, nil
		}
	}

	return nil, fmt.Errorf("meal not found")
}

func (r FakeMealRepository) Add(m *Meal) error {
	r.meals[m.Name] = m

	return nil
}
