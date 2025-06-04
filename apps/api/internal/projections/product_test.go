package projections_test

import (
	"database/sql"
	sqlStore "github.com/hallgren/eventsourcing/eventstore/sql"
	"github.com/joe-reed/meal-planner/apps/api/internal/database"
	"github.com/joe-reed/meal-planner/apps/api/internal/domain/category"
	"github.com/joe-reed/meal-planner/apps/api/internal/domain/product"
	"github.com/joe-reed/meal-planner/apps/api/internal/projections"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

type ShoppingListSuite struct {
	suite.Suite
	productRepository *product.EventSourcedProductRepository
	db                *sql.DB
	es                *sqlStore.SQL
}

func (suite *ShoppingListSuite) SetupTest() {
	db, err := database.CreateDatabase(":memory:")
	assert.NoError(suite.T(), err)

	suite.db = db
	suite.es = sqlStore.Open(db)

	suite.createRepositories(db)
}

func (suite *ShoppingListSuite) TearDownTest() {
	err := suite.db.Close()

	if err != nil {
		suite.T().Fatal(err)
	}
}

func (suite *ShoppingListSuite) createRepositories(db *sql.DB) {
	productRepository, err := product.NewSqliteProductRepository(db)
	assert.NoError(suite.T(), err)
	suite.productRepository = productRepository
}

func (suite *ShoppingListSuite) addProduct(id string, name product.ProductName, category category.CategoryName) *product.Product {
	i, err := product.NewProduct(id, name, category)
	assert.NoError(suite.T(), err)

	err = suite.productRepository.Add(i)
	assert.NoError(suite.T(), err)

	return i
}

func (suite *ShoppingListSuite) TestProductProjection() {
	productA := suite.addProduct("ing-a", "Ing A", category.AlcoholicDrinks)
	productB := suite.addProduct("ing-b", "Ing B", category.AlcoholicDrinks)
	productC := suite.addProduct("ing-c", "Ing C", category.Dairy)

	p, output := projections.CreateProductProjection(suite.es)

	result := p.RunToEnd(suite.T().Context())

	assert.NoError(suite.T(), result.Error)
	assert.NotEmpty(suite.T(), output)

	assert.Equal(suite.T(), 2, len(output[category.AlcoholicDrinks]))
	assert.Equal(suite.T(), productA, output[category.AlcoholicDrinks][0])
	assert.Equal(suite.T(), productB, output[category.AlcoholicDrinks][1])
	assert.Equal(suite.T(), 1, len(output[category.Dairy]))
	assert.Equal(suite.T(), productC, output[category.Dairy][0])
}
