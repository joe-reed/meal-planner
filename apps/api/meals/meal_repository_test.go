package meals_test

import (
	"os"
	"testing"

	"github.com/joe-reed/meal-planner/apps/api/meals"
	"github.com/stretchr/testify/assert"
)

func TestFakeMealRepository(t *testing.T) {
	runSuite(t, func() *meals.MealRepository {
		return meals.NewFakeMealRepository()
	}, func() {})
}

func TestSqliteMealRepository(t *testing.T) {
	runSuite(t, func() *meals.MealRepository {
		r, err := meals.NewSqliteMealRepository("test.db")
		assert.NoError(t, err)
		return r
	}, func() {
		err := os.Remove("test.db")
		assert.NoError(t, err)
	})
}

func runSuite(t *testing.T, factory func() *meals.MealRepository, teardown func()) {
	tests := []struct {
		title string
		run   func(t *testing.T, r *meals.MealRepository)
	}{
		{"adding a meal", testAddingMeal},
		{"getting all meals", testGettingMeals},
		{"finding a meal", testFindingMeal},
		{"saving a meal", testSavingMeal},
	}

	for _, test := range tests {
		t.Run(test.title, func(t *testing.T) {
			test.run(t, factory())
			teardown()
		})
	}
}

func testAddingMeal(t *testing.T, r *meals.MealRepository) {
	m := meals.NewMealBuilder().Build()
	err := r.Save(m)
	assert.NoError(t, err)

	retrieved, err := r.Get()

	assert.NoError(t, err)
	assert.Len(t, retrieved, 1)
	assert.EqualExportedValues(t, m, retrieved[0])
}

func testGettingMeals(t *testing.T, r *meals.MealRepository) {
	m1 := meals.NewMealBuilder().AddIngredient(meals.MealIngredient{IngredientId: "c"}).WithName("c").Build()
	m2 := meals.NewMealBuilder().AddIngredient(meals.MealIngredient{IngredientId: "b"}).WithName("b").Build()
	m3 := meals.NewMealBuilder().AddIngredient(meals.MealIngredient{IngredientId: "a"}).WithName("a").Build()

	err := r.Save(m1)
	assert.NoError(t, err)
	err = r.Save(m2)
	assert.NoError(t, err)
	err = r.Save(m3)
	assert.NoError(t, err)

	m, err := r.Get()
	assert.NoError(t, err)
	assert.Len(t, m, 3)
	assert.EqualExportedValues(t, m3, m[0])
	assert.EqualExportedValues(t, m2, m[1])
	assert.EqualExportedValues(t, m1, m[2])
}

func testFindingMeal(t *testing.T, r *meals.MealRepository) {
	m := meals.NewMealBuilder().AddIngredient(meals.MealIngredient{IngredientId: "a"}).WithName("a").Build()
	err := r.Save(m)
	assert.NoError(t, err)

	found, err := r.Find(m.Id)

	assert.NoError(t, err)
	assert.EqualExportedValues(t, m, found)
}

func testSavingMeal(t *testing.T, r *meals.MealRepository) {
	m := meals.NewMealBuilder().AddIngredient(meals.MealIngredient{IngredientId: "a"}).WithName("a").Build()
	err := r.Save(m)
	assert.NoError(t, err)

	m.AddIngredient(meals.MealIngredient{IngredientId: "b"})
	err = r.Save(m)
	assert.NoError(t, err)

	found, err := r.Find(m.Id)
	assert.NoError(t, err)
	assert.EqualExportedValues(t, m, found)
	assert.Equal(t, []meals.MealIngredient{{"a"}, {"b"}}, found.MealIngredients)
}
