package ingredients_test

import (
	"github.com/joe-reed/meal-planner/apps/api/ingredients"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestFakeIngredientRepository(t *testing.T) {
	runSuite(t, func() ingredients.IngredientRepository {
		return ingredients.NewFakeIngredientRepository()
	}, func() {})
}

func TestSqliteIngredientRepository(t *testing.T) {
	runSuite(t, func() ingredients.IngredientRepository {
		r, err := ingredients.NewSqliteIngredientRepository("test.db")
		assert.NoError(t, err)
		return r
	}, func() {
		err := os.Remove("test.db")
		assert.NoError(t, err)
	})
}

func runSuite(t *testing.T, factory func() ingredients.IngredientRepository, teardown func()) {
	tests := []struct {
		title string
		run   func(t *testing.T, r ingredients.IngredientRepository)
	}{
		{"adding an ingredient", testAddingIngredient},
		{"getting all ingredients", testGettingIngredients},
		{"getting empty list of ingredients", testGettingZeroIngredients},
	}

	for _, test := range tests {
		t.Run(test.title, func(t *testing.T) {
			test.run(t, factory())
			teardown()
		})
	}
}

func testAddingIngredient(t *testing.T, r ingredients.IngredientRepository) {
	expected := ingredients.NewIngredientBuilder().Build()
	err := r.Add(expected)
	assert.NoError(t, err)

	actual, err := r.Get()

	assert.NoError(t, err)
	assert.Contains(t, actual, expected)
}

func testGettingIngredients(t *testing.T, r ingredients.IngredientRepository) {
	i1 := ingredients.NewIngredientBuilder().WithName("c").Build()
	i2 := ingredients.NewIngredientBuilder().WithName("b").Build()
	i3 := ingredients.NewIngredientBuilder().WithName("a").Build()

	err := r.Add(i1)
	assert.NoError(t, err)
	err = r.Add(i2)
	assert.NoError(t, err)
	err = r.Add(i3)
	assert.NoError(t, err)

	i, err := r.Get()
	assert.NoError(t, err)
	assert.Equal(t, []*ingredients.Ingredient{i3, i2, i1}, i)
}

func testGettingZeroIngredients(t *testing.T, r ingredients.IngredientRepository) {
	i, err := r.Get()
	assert.NoError(t, err)
	assert.Equal(t, []*ingredients.Ingredient{}, i)
}
