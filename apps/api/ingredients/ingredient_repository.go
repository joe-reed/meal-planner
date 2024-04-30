package ingredients

import (
	"database/sql"
	"sort"

	_ "github.com/mattn/go-sqlite3"
)

type IngredientRepository interface {
	Get() ([]*Ingredient, error)
	Add(ingredient *Ingredient) error
}

type SqliteIngredientRepository struct {
	db *sql.DB
}

func NewSqliteIngredientRepository(dbFile string) (*SqliteIngredientRepository, error) {
	db, err := sql.Open("sqlite3", dbFile)
	if err != nil {
		return nil, err
	}
	if _, err := db.Exec("CREATE TABLE IF NOT EXISTS ingredients (id VARCHAR(255) NOT NULL PRIMARY KEY, name VARCHAR(255) NOT NULL)"); err != nil {
		return nil, err
	}

	return &SqliteIngredientRepository{db: db}, nil
}

func (r SqliteIngredientRepository) Get() ([]*Ingredient, error) {
	rows, err := r.db.Query("SELECT id, name FROM ingredients ORDER BY name")

	if err != nil {
		return nil, err
	}

	result := []*Ingredient{}

	for rows.Next() {
		var i Ingredient

		err = rows.Scan(&i.Id, &i.Name)

		if err != nil {
			return nil, err
		}

		result = append(result, &i)
	}

	return result, nil
}

func (r SqliteIngredientRepository) Add(i *Ingredient) error {
	_, err := r.db.Exec("INSERT INTO ingredients (id, name) VALUES (?, ?)", i.Id, i.Name)

	return err
}

type FakeIngredientRepository struct {
	ingredients map[string]*Ingredient
}

func NewFakeIngredientRepository() FakeIngredientRepository {
	return FakeIngredientRepository{ingredients: map[string]*Ingredient{}}
}

func (r FakeIngredientRepository) Get() ([]*Ingredient, error) {
	keys := make([]string, 0, len(r.ingredients))
	for k := range r.ingredients {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	v := make([]*Ingredient, 0, len(r.ingredients))
	for _, k := range keys {
		v = append(v, r.ingredients[k])
	}
	return v, nil
}

func (r FakeIngredientRepository) Add(i *Ingredient) error {
	r.ingredients[i.Name] = i

	return nil
}
