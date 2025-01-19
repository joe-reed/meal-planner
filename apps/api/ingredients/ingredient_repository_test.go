package ingredients_test

import (
	"github.com/joe-reed/meal-planner/apps/api/categories"
	"github.com/joe-reed/meal-planner/apps/api/database"
	"github.com/joe-reed/meal-planner/apps/api/ingredients"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestFakeIngredientRepository(t *testing.T) {
	runSuite(t, func() *ingredients.IngredientRepository {
		return ingredients.NewFakeIngredientRepository()
	}, func() {})
}

func TestSqliteIngredientRepository(t *testing.T) {
	runSuite(t, func() *ingredients.IngredientRepository {
		db, err := database.CreateDatabase("test.db")
		assert.NoError(t, err)
		r, err := ingredients.NewSqliteIngredientRepository(db)
		assert.NoError(t, err)
		return r
	}, func() {
		err := os.Remove("test.db")
		assert.NoError(t, err)
	})
}

func runSuite(t *testing.T, factory func() *ingredients.IngredientRepository, teardown func()) {
	tests := []struct {
		title string
		run   func(t *testing.T, r *ingredients.IngredientRepository)
	}{
		{"adding an ingredient", testAddingIngredient},
		{"getting all ingredients", testGettingIngredients},
		{"getting empty list of ingredients", testGettingZeroIngredients},
		{"getting ingredient by name", testGettingIngredientByName},
	}

	for _, test := range tests {
		t.Run(test.title, func(t *testing.T) {
			test.run(t, factory())
			teardown()
		})
	}
}

func testAddingIngredient(t *testing.T, r *ingredients.IngredientRepository) {
	expected := ingredients.NewIngredientBuilder().Build()
	err := r.Add(expected)
	assert.NoError(t, err)

	actual, err := r.Get()

	assert.NoError(t, err)
	assert.Len(t, actual, 1)
	assert.EqualExportedValues(t, expected, actual[0])
}

func testGettingIngredients(t *testing.T, r *ingredients.IngredientRepository) {
	i1 := ingredients.NewIngredientBuilder().WithName("c").WithCategory(categories.Frozen).Build()
	i2 := ingredients.NewIngredientBuilder().WithName("b").WithCategory(categories.Vegetables).Build()
	i3 := ingredients.NewIngredientBuilder().WithName("a").WithCategory(categories.SeedsNutsAndDriedFruits).Build()

	err := r.Add(i1)
	assert.NoError(t, err)
	err = r.Add(i2)
	assert.NoError(t, err)
	err = r.Add(i3)
	assert.NoError(t, err)

	i, err := r.Get()
	assert.NoError(t, err)
	assert.Len(t, i, 3)
	assert.EqualExportedValues(t, i3, i[0])
	assert.EqualExportedValues(t, i2, i[1])
	assert.EqualExportedValues(t, i1, i[2])
}

func testGettingZeroIngredients(t *testing.T, r *ingredients.IngredientRepository) {
	i, err := r.Get()
	assert.NoError(t, err)
	assert.Len(t, i, 0)
}

func testGettingIngredientByName(t *testing.T, r *ingredients.IngredientRepository) {
	i := ingredients.NewIngredientBuilder().WithName("test name").WithCategory(categories.Frozen).Build()

	err := r.Add(i)
	assert.NoError(t, err)

	found, err := r.GetByName("test name")
	assert.NoError(t, err)
	assert.EqualExportedValues(t, found, i)
}
