package meals_test

import (
	"os"
	"testing"

	"github.com/joe-reed/meal-planner/meals"
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

		if err != nil {
			t.Fatalf("Error: %v", err)
		}

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
	assert.Contains(t, r.Get(), m)
}

func testGettingMeals(t *testing.T, r meals.MealRepository) {
	m1 := meals.NewMealBuilder().WithName("c").Build()
	m2 := meals.NewMealBuilder().WithName("b").Build()
	m3 := meals.NewMealBuilder().WithName("a").Build()

	r.Add(m1)
	r.Add(m2)
	r.Add(m3)

	assert.Equal(t, []*meals.Meal{m3, m2, m1}, r.Get())
}
