package shops_test

import (
	"github.com/joe-reed/meal-planner/apps/api/database"
	"os"
	"testing"

	"github.com/joe-reed/meal-planner/apps/api/shops"
	"github.com/stretchr/testify/assert"
)

func TestFakeShopRepository(t *testing.T) {
	runSuite(t, func() *shops.ShopRepository {
		return shops.NewFakeShopRepository()
	}, func() {})
}

func TestSqliteMealRepository(t *testing.T) {
	runSuite(t, func() *shops.ShopRepository {
		db, err := database.CreateDatabase("test.db")
		assert.NoError(t, err)
		r, err := shops.NewSqliteShopRepository(db)
		assert.NoError(t, err)
		return r
	}, func() {
		err := os.Remove("test.db")
		assert.NoError(t, err)
	})
}

func runSuite(t *testing.T, factory func() *shops.ShopRepository, teardown func()) {
	tests := []struct {
		title string
		run   func(t *testing.T, r *shops.ShopRepository)
	}{
		{"getting current shop", testGettingCurrentShop},
		{"getting current shop with no meals", testGettingCurrentShopWithNoMeals},
		{"getting current shop when no shops exist", testGettingCurrentShopIfNoShopsExist},
		{"finding shop", testFindingShop},
		{"finding shop with no meals", testFindingShopWithNoMeals},
		{"saving shop", testSavingShop},
	}

	for _, test := range tests {
		t.Run(test.title, func(t *testing.T) {
			test.run(t, factory())
			teardown()
		})
	}
}

func testGettingCurrentShop(t *testing.T, r *shops.ShopRepository) {
	s1, err := shops.NewShop(1)
	assert.NoError(t, err)
	s2, err := shops.NewShop(2)
	assert.NoError(t, err)
	s3, err := shops.NewShop(3)
	assert.NoError(t, err)

	s3.AddMeal(&shops.ShopMeal{MealId: "123"}).AddMeal(&shops.ShopMeal{MealId: "456"})

	err = r.Save(s1)
	assert.NoError(t, err)
	err = r.Save(s2)
	assert.NoError(t, err)
	err = r.Save(s3)
	assert.NoError(t, err)

	s, err := r.Current()
	assert.NoError(t, err)

	assert.EqualExportedValues(t, s3, s)
}

func testGettingCurrentShopIfNoShopsExist(t *testing.T, r *shops.ShopRepository) {
	s, err := r.Current()
	assert.NoError(t, err)
	assert.Nil(t, s)
}

func testGettingCurrentShopWithNoMeals(t *testing.T, r *shops.ShopRepository) {
	s1, err := shops.NewShop(1)
	assert.NoError(t, err)

	err = r.Save(s1)
	assert.NoError(t, err)

	s, err := r.Current()
	assert.NoError(t, err)
	assert.Len(t, s.Meals, 0)
}

func testFindingShop(t *testing.T, r *shops.ShopRepository) {
	s, err := shops.NewShop(1)
	assert.NoError(t, err)
	s.AddMeal(&shops.ShopMeal{MealId: "abc"})

	err = r.Save(s)
	assert.NoError(t, err)

	found, err := r.Find(1)
	assert.NoError(t, err)
	assert.EqualExportedValues(t, s, found)
}

func testFindingShopWithNoMeals(t *testing.T, r *shops.ShopRepository) {
	s, err := shops.NewShop(1)
	assert.NoError(t, err)

	err = r.Save(s)
	assert.NoError(t, err)

	found, err := r.Find(1)
	assert.NoError(t, err)
	assert.EqualExportedValues(t, s, found)
	assert.Len(t, found.Meals, 0)
}

func testSavingShop(t *testing.T, r *shops.ShopRepository) {
	s, err := shops.NewShop(1)
	assert.NoError(t, err)

	err = r.Save(s)
	assert.NoError(t, err)

	s = s.AddMeal(&shops.ShopMeal{MealId: "abc"})

	err = r.Save(s)
	assert.NoError(t, err)

	found, err := r.Find(1)
	assert.NoError(t, err)
	assert.EqualExportedValues(t, s, found)
}
