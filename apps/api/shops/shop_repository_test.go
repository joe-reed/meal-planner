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
		os.Remove("test.db")
	})
}

func runSuite(t *testing.T, factory func() shops.ShopRepository, teardown func()) {
	tests := []struct {
		title string
		run   func(t *testing.T, r shops.ShopRepository)
	}{
		{"getting current shop", testGettingCurrentShop},
		{"getting current shop when no shops exist", testGettingCurrentShopIfNoShopsExist},
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

	r.Add(&s1)
	r.Add(&s2)
	r.Add(&s3)

	s, err := r.Current()
	assert.NoError(t, err)
	assert.Equal(t, &shops.Shop{Id: 3}, s)
}

func testGettingCurrentShopIfNoShopsExist(t *testing.T, r shops.ShopRepository) {
	s, err := r.Current()
	assert.NoError(t, err)
	assert.Nil(t, s)
}
