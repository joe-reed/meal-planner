package meals_test

import (
	"os"
	"testing"

	"github.com/joe-reed/meal-planner/apps/api/meals"
	"github.com/stretchr/testify/assert"
)

func TestFakeMealRepository(t *testing.T) {
	runSuite(t, func() meals.MealRepository {
		return meals.NewFakeMealRepository()
	}, func() {})
}

func TestSqliteMealRepository(t *testing.T) {
	runSuite(t, func() meals.MealRepository {
		r, err := meals.NewSqliteMealRepository("test.db")
		assert.NoError(t, err)
		return r
	}, func() {
		os.Remove("test.db")
	})
}

func runSuite(t *testing.T, factory func() meals.MealRepository, teardown func()) {
	tests := []struct {
		title string
		run   func(t *testing.T, r meals.MealRepository)
	}{
		{"adding a meal", testAddingMeal},
		{"getting all meals", testGettingMeals},
		{"finding a meal", testFindingMeal},
	}

	for _, test := range tests {
		t.Run(test.title, func(t *testing.T) {
			test.run(t, factory())
			teardown()
		})
	}
}

func testAddingMeal(t *testing.T, r meals.MealRepository) {
	m := meals.NewMealBuilder().Build()
	r.Add(m)

	meals, err := r.Get()

	assert.NoError(t, err)
	assert.Contains(t, meals, m)
}

func testGettingMeals(t *testing.T, r meals.MealRepository) {
	m1 := meals.NewMealBuilder().AddIngredient(meals.MealIngredient{"c"}).WithName("c").Build()
	m2 := meals.NewMealBuilder().AddIngredient(meals.MealIngredient{"b"}).WithName("b").Build()
	m3 := meals.NewMealBuilder().AddIngredient(meals.MealIngredient{"a"}).WithName("a").Build()

	err := r.Add(m1)
	assert.NoError(t, err)
	err = r.Add(m2)
	assert.NoError(t, err)
	err = r.Add(m3)
	assert.NoError(t, err)

	m, err := r.Get()
	assert.NoError(t, err)
	assert.Equal(t, []*meals.Meal{m3, m2, m1}, m)
}

func testFindingMeal(t *testing.T, r meals.MealRepository) {
	m := meals.NewMealBuilder().AddIngredient(meals.MealIngredient{"a"}).WithName("a").Build()
	r.Add(m)

	found, err := r.Find(m.Id)

	assert.NoError(t, err)
	assert.Equal(t, m, found)
}
