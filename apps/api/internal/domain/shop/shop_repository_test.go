package shop_test

import (
	"github.com/joe-reed/meal-planner/apps/api/internal/database"
	"github.com/joe-reed/meal-planner/apps/api/internal/domain/shop"
	"github.com/stretchr/testify/assert"
	"os"
	"testing"
)

func TestFakeShopRepository(t *testing.T) {
	runSuite(t, func() *shop.ShopRepository {
		return shop.NewFakeShopRepository()
	}, func() {})
}

func TestSqliteMealRepository(t *testing.T) {
	runSuite(t, func() *shop.ShopRepository {
		db, err := database.CreateDatabase("test.db")
		assert.NoError(t, err)
		r, err := shop.NewSqliteShopRepository(db)
		assert.NoError(t, err)
		return r
	}, func() {
		err := os.Remove("test.db")
		assert.NoError(t, err)
	})
}

func runSuite(t *testing.T, factory func() *shop.ShopRepository, teardown func()) {
	tests := []struct {
		title string
		run   func(t *testing.T, r *shop.ShopRepository)
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

func testGettingCurrentShop(t *testing.T, r *shop.ShopRepository) {
	s1, err := shop.NewShop(1)
	assert.NoError(t, err)
	s2, err := shop.NewShop(2)
	assert.NoError(t, err)
	s3, err := shop.NewShop(3)
	assert.NoError(t, err)

	s3.AddMeal(&shop.ShopMeal{MealId: "123"}).AddMeal(&shop.ShopMeal{MealId: "456"})

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

func testGettingCurrentShopIfNoShopsExist(t *testing.T, r *shop.ShopRepository) {
	s, err := r.Current()
	assert.NoError(t, err)
	assert.Nil(t, s)
}

func testGettingCurrentShopWithNoMeals(t *testing.T, r *shop.ShopRepository) {
	s1, err := shop.NewShop(1)
	assert.NoError(t, err)

	err = r.Save(s1)
	assert.NoError(t, err)

	s, err := r.Current()
	assert.NoError(t, err)
	assert.Len(t, s.Meals, 0)
}

func testFindingShop(t *testing.T, r *shop.ShopRepository) {
	s, err := shop.NewShop(1)
	assert.NoError(t, err)
	s.AddMeal(&shop.ShopMeal{MealId: "abc"})

	err = r.Save(s)
	assert.NoError(t, err)

	found, err := r.Find(1)
	assert.NoError(t, err)
	assert.EqualExportedValues(t, s, found)
}

func testFindingShopWithNoMeals(t *testing.T, r *shop.ShopRepository) {
	s, err := shop.NewShop(1)
	assert.NoError(t, err)

	err = r.Save(s)
	assert.NoError(t, err)

	found, err := r.Find(1)
	assert.NoError(t, err)
	assert.EqualExportedValues(t, s, found)
	assert.Len(t, found.Meals, 0)
}

func testSavingShop(t *testing.T, r *shop.ShopRepository) {
	s, err := shop.NewShop(1)
	assert.NoError(t, err)

	err = r.Save(s)
	assert.NoError(t, err)

	s = s.AddMeal(&shop.ShopMeal{MealId: "abc"})

	err = r.Save(s)
	assert.NoError(t, err)

	found, err := r.Find(1)
	assert.NoError(t, err)
	assert.EqualExportedValues(t, s, found)
}
