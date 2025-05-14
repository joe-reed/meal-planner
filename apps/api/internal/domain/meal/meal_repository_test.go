package meal_test

import (
	"github.com/joe-reed/meal-planner/apps/api/internal/database"
	"github.com/joe-reed/meal-planner/apps/api/internal/domain/meal"
	"github.com/stretchr/testify/assert"
	"os"
	"testing"
)

func TestFakeMealRepository(t *testing.T) {
	runSuite(t, func() *meal.EventSourcedMealRepository {
		return meal.NewFakeMealRepository()
	}, func() {})
}

func TestSqliteMealRepository(t *testing.T) {
	runSuite(t, func() *meal.EventSourcedMealRepository {
		db, err := database.CreateDatabase("test.db")
		assert.NoError(t, err)
		r, err := meal.NewSqliteMealRepository(db)
		assert.NoError(t, err)
		return r
	}, func() {
		err := os.Remove("test.db")
		assert.NoError(t, err)
	})
}

func runSuite(t *testing.T, factory func() *meal.EventSourcedMealRepository, teardown func()) {
	tests := []struct {
		title string
		run   func(t *testing.T, r *meal.EventSourcedMealRepository)
	}{
		{"adding a meal", testAddingMeal},
		{"getting all meals", testGettingMeals},
		{"getting all meals without ingredients", testGettingMealsWithoutIngredients},
		{"finding a meal", testFindingMeal},
		{"saving a meal", testSavingMeal},
		{"updating a meal's url", testUpdatingMealUrl},
	}

	for _, test := range tests {
		t.Run(test.title, func(t *testing.T) {
			test.run(t, factory())
			teardown()
		})
	}
}

func testAddingMeal(t *testing.T, r *meal.EventSourcedMealRepository) {
	m := meal.NewMealBuilder().Build()
	err := r.Save(m)
	assert.NoError(t, err)

	retrieved, err := r.Get()

	assert.NoError(t, err)
	assert.Len(t, retrieved, 1)
	assert.EqualExportedValues(t, m, retrieved[0])
}

func testGettingMeals(t *testing.T, r *meal.EventSourcedMealRepository) {
	m1 := meal.NewMealBuilder().AddIngredient(meal.Ingredient{ProductId: "c"}).WithName("c").Build()
	m2 := meal.NewMealBuilder().AddIngredient(meal.Ingredient{ProductId: "b"}).WithName("b").Build()
	m3 := meal.NewMealBuilder().AddIngredient(meal.Ingredient{ProductId: "a"}).WithName("a").Build()

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

func testGettingMealsWithoutIngredients(t *testing.T, r *meal.EventSourcedMealRepository) {
	m1 := meal.NewMealBuilder().WithName("c").Build()
	m2 := meal.NewMealBuilder().WithName("b").Build()
	m3 := meal.NewMealBuilder().WithName("a").Build()

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

func testFindingMeal(t *testing.T, r *meal.EventSourcedMealRepository) {
	m := meal.NewMealBuilder().AddIngredient(*meal.NewIngredient("a")).WithName("a").Build()
	err := r.Save(m)
	assert.NoError(t, err)

	found, err := r.Find(m.Id)

	assert.NoError(t, err)
	assert.EqualExportedValues(t, m, found)
}

func testSavingMeal(t *testing.T, r *meal.EventSourcedMealRepository) {
	m := meal.NewMealBuilder().AddIngredient(*meal.NewIngredient("a")).WithName("a").Build()
	err := r.Save(m)
	assert.NoError(t, err)

	m.AddIngredient(*meal.NewIngredient("b"))
	err = r.Save(m)
	assert.NoError(t, err)

	found, err := r.Find(m.Id)
	assert.NoError(t, err)
	assert.EqualExportedValues(t, m, found)
	assert.Equal(t, []meal.Ingredient{*meal.NewIngredient("a"), *meal.NewIngredient("b")}, found.Ingredients)
}

func testUpdatingMealUrl(t *testing.T, r *meal.EventSourcedMealRepository) {
	m := meal.NewMealBuilder().AddIngredient(*meal.NewIngredient("a")).WithName("a").Build()
	err := r.Save(m)
	assert.NoError(t, err)

	m.UpdateUrl("https://test.localhost")

	err = r.Save(m)
	assert.NoError(t, err)

	found, err := r.Find(m.Id)
	assert.NoError(t, err)
	assert.EqualExportedValues(t, m, found)
}
