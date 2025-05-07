package product_test

import (
	"github.com/joe-reed/meal-planner/apps/api/internal/database"
	"github.com/joe-reed/meal-planner/apps/api/internal/domain/category"
	"github.com/joe-reed/meal-planner/apps/api/internal/domain/product"
	"github.com/stretchr/testify/assert"
	"os"
	"testing"
)

func TestFakeProductRepository(t *testing.T) {
	runSuite(t, func() *product.EventSourcedProductRepository {
		return product.NewFakeProductRepository()
	}, func() {})
}

func TestSqliteProductRepository(t *testing.T) {
	runSuite(t, func() *product.EventSourcedProductRepository {
		db, err := database.CreateDatabase("test.db")
		assert.NoError(t, err)
		r, err := product.NewSqliteProductRepository(db)
		assert.NoError(t, err)
		return r
	}, func() {
		err := os.Remove("test.db")
		assert.NoError(t, err)
	})
}

func runSuite(t *testing.T, factory func() *product.EventSourcedProductRepository, teardown func()) {
	tests := []struct {
		title string
		run   func(t *testing.T, r *product.EventSourcedProductRepository)
	}{
		{"adding an product", testAddingProduct},
		{"getting all products", testGettingProducts},
		{"getting empty list of products", testGettingZeroProducts},
		{"getting product by name", testGettingProductByName},
	}

	for _, test := range tests {
		t.Run(test.title, func(t *testing.T) {
			test.run(t, factory())
			teardown()
		})
	}
}

func testAddingProduct(t *testing.T, r *product.EventSourcedProductRepository) {
	expected := product.NewProductBuilder().Build()
	err := r.Add(expected)
	assert.NoError(t, err)

	actual, err := r.Get()

	assert.NoError(t, err)
	assert.Len(t, actual, 1)
	assert.EqualExportedValues(t, expected, actual[0])
}

func testGettingProducts(t *testing.T, r *product.EventSourcedProductRepository) {
	i1 := product.NewProductBuilder().WithName("c").WithCategory(category.Frozen).Build()
	i2 := product.NewProductBuilder().WithName("b").WithCategory(category.Vegetables).Build()
	i3 := product.NewProductBuilder().WithName("a").WithCategory(category.SeedsNutsAndDriedFruits).Build()

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

func testGettingZeroProducts(t *testing.T, r *product.EventSourcedProductRepository) {
	i, err := r.Get()
	assert.NoError(t, err)
	assert.Len(t, i, 0)
}

func testGettingProductByName(t *testing.T, r *product.EventSourcedProductRepository) {
	i := product.NewProductBuilder().WithName("test name").WithCategory(category.Frozen).Build()

	err := r.Add(i)
	assert.NoError(t, err)

	found, err := r.GetByName("test name")
	assert.NoError(t, err)
	assert.EqualExportedValues(t, found, i)
}
