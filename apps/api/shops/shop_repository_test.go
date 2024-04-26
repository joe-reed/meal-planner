package shops_test

import (
	"os"
	"testing"

	"github.com/joe-reed/meal-planner/apps/api/shops"
	"github.com/stretchr/testify/assert"
)

func TestFakeShopRepository(t *testing.T) {
	runSuite(t, func() shops.ShopRepository {
		return shops.NewFakeShopRepository()
	}, func() {})
}

func TestSqliteMealRepository(t *testing.T) {
	runSuite(t, func() shops.ShopRepository {
		r, err := shops.NewSqliteShopRepository("test.db")
		assert.NoError(t, err)
		return r
	}, func() {
		err := os.Remove("test.db")
		assert.NoError(t, err)
	})
}

func runSuite(t *testing.T, factory func() shops.ShopRepository, teardown func()) {
	tests := []struct {
		title string
		run   func(t *testing.T, r shops.ShopRepository)
	}{
		{"getting current shop", testGettingCurrentShop},
		{"getting current shop when no shops exist", testGettingCurrentShopIfNoShopsExist},
		{"finding shop", testFindingShop},
		{"saving shop", testSavingShop},
	}

	for _, test := range tests {
		t.Run(test.title, func(t *testing.T) {
			test.run(t, factory())
			teardown()
		})
	}
}

func testGettingCurrentShop(t *testing.T, r shops.ShopRepository) {
	s1 := shops.Shop{Id: 1}
	s2 := shops.Shop{Id: 2}
	s3 := shops.Shop{Id: 3}

	err := r.Add(&s1)
	assert.NoError(t, err)
	err = r.Add(&s2)
	assert.NoError(t, err)
	err = r.Add(&s3)
	assert.NoError(t, err)

	s, err := r.Current()
	assert.NoError(t, err)
	assert.Equal(t, &shops.Shop{Id: 3}, s)
}

func testGettingCurrentShopIfNoShopsExist(t *testing.T, r shops.ShopRepository) {
	s, err := r.Current()
	assert.NoError(t, err)
	assert.Nil(t, s)
}

func testFindingShop(t *testing.T, r shops.ShopRepository) {
	s := shops.NewShop(1).AddMeal(&shops.ShopMeal{"abc"})

	err := r.Add(s)
	assert.NoError(t, err)

	found, err := r.Find(1)
	assert.NoError(t, err)
	assert.Equal(t, s, found)
}

func testSavingShop(t *testing.T, r shops.ShopRepository) {
	s := shops.NewShop(1)

	err := r.Add(s)
	assert.NoError(t, err)

	s = s.AddMeal(&shops.ShopMeal{"abc"})

	err = r.Save(s)
	assert.NoError(t, err)

	found, err := r.Find(1)
	assert.NoError(t, err)
	assert.Equal(t, s, found)
}
