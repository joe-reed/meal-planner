package shoppinglist_test

import (
	"context"
	"database/sql"
	"github.com/brianvoe/gofakeit/v7"
	sqlStore "github.com/hallgren/eventsourcing/eventstore/sql"
	"github.com/joe-reed/meal-planner/apps/api/internal/database"
	"github.com/joe-reed/meal-planner/apps/api/internal/domain/basket"
	"github.com/joe-reed/meal-planner/apps/api/internal/domain/category"
	"github.com/joe-reed/meal-planner/apps/api/internal/domain/meal"
	"github.com/joe-reed/meal-planner/apps/api/internal/domain/product"
	"github.com/joe-reed/meal-planner/apps/api/internal/domain/quantity"
	"github.com/joe-reed/meal-planner/apps/api/internal/domain/shop"
	"github.com/joe-reed/meal-planner/apps/api/internal/domain/shoppinglist"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	"strconv"
	"testing"
)

type ShoppingListSuite struct {
	suite.Suite
	productRepository *product.EventSourcedProductRepository
	shopRepository    *shop.ShopRepository
	mealRepository    *meal.EventSourcedMealRepository
	basketRepository  *basket.BasketRepository
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

func (suite *ShoppingListSuite) TestAddingMealToShop() {
	productA := suite.addProduct("ing-a", "Ing A", category.AlcoholicDrinks)

	m := suite.addMeal([]meal.Ingredient{*meal.NewIngredient(productA.Id)})

	s, _ := suite.addShop()

	suite.addMealToShop(s, m)

	output := suite.runProjection()

	assert.EqualExportedValues(suite.T(),
		map[string]shoppinglist.ShopItem{
			productA.Id: {Product: *productA, MealCount: 1, IsInBasket: false, Quantities: []quantity.Quantity{{Amount: 1, Unit: quantity.Number}}},
		},
		*output.ShoppingList,
	)
}

func (suite *ShoppingListSuite) TestAddingTwoMealsWithSameIngredient() {
	productA := suite.addProduct("ing-a", "Ing A", category.AlcoholicDrinks)
	productB := suite.addProduct("ing-b", "Ing B", category.Bakery)
	productC := suite.addProduct("ing-c", "Ing C", category.Dairy)

	meal1 := suite.addMeal([]meal.Ingredient{
		*meal.NewIngredient(productA.Id).WithQuantity(2, quantity.Tbsp),
		*meal.NewIngredient(productB.Id).WithQuantity(100, quantity.Ml),
	})
	meal2 := suite.addMeal([]meal.Ingredient{
		*meal.NewIngredient(productB.Id).WithQuantity(50, quantity.Gram),
		*meal.NewIngredient(productC.Id).WithQuantity(1, quantity.Litre),
	})

	s, _ := suite.addShop()

	suite.addMealToShop(s, meal1)
	suite.addMealToShop(s, meal2)

	output := suite.runProjection()

	assert.EqualExportedValues(suite.T(),
		map[string]shoppinglist.ShopItem{
			productA.Id: {Product: *productA, MealCount: 1, IsInBasket: false, Quantities: []quantity.Quantity{{Amount: 2, Unit: quantity.Tbsp}}},
			productB.Id: {Product: *productB, MealCount: 2, IsInBasket: false, Quantities: []quantity.Quantity{{Amount: 100, Unit: quantity.Ml}, {Amount: 50, Unit: quantity.Gram}}},
			productC.Id: {Product: *productC, MealCount: 1, IsInBasket: false, Quantities: []quantity.Quantity{{Amount: 1, Unit: quantity.Litre}}},
		},
		*output.ShoppingList,
	)
}

func (suite *ShoppingListSuite) TestRemovingMeal() {
	productA := suite.addProduct("ing-a", "Ing A", category.AlcoholicDrinks)
	productB := suite.addProduct("ing-b", "Ing B", category.Bakery)
	productC := suite.addProduct("ing-c", "Ing C", category.Dairy)

	meal1 := suite.addMeal([]meal.Ingredient{*meal.NewIngredient(productA.Id), *meal.NewIngredient(productB.Id)})
	meal2 := suite.addMeal([]meal.Ingredient{*meal.NewIngredient(productB.Id), *meal.NewIngredient(productC.Id)})

	s, _ := suite.addShop()

	suite.addMealToShop(s, meal1)
	suite.addMealToShop(s, meal2)
	suite.removeMealFromShop(s, meal2)

	output := suite.runProjection()

	assert.EqualExportedValues(suite.T(),
		map[string]shoppinglist.ShopItem{
			productA.Id: {Product: *productA, MealCount: 1, IsInBasket: false, Quantities: []quantity.Quantity{{Amount: 1, Unit: quantity.Number}}},
			productB.Id: {Product: *productB, MealCount: 1, IsInBasket: false, Quantities: []quantity.Quantity{{Amount: 1, Unit: quantity.Number}}},
		},
		*output.ShoppingList,
	)
}

func (suite *ShoppingListSuite) TestAddingIngredientToMealBeforeAddingToShop() {
	productA := suite.addProduct("ing-a", "Ing A", category.AlcoholicDrinks)

	m := suite.addMeal([]meal.Ingredient{*meal.NewIngredient(productA.Id)})

	s, _ := suite.addShop()

	productB := suite.addProduct("ing-b", "Ing B", category.Bakery)

	suite.addIngredientToMeal(m, productB)

	suite.addMealToShop(s, m)

	output := suite.runProjection()

	assert.EqualExportedValues(suite.T(),
		map[string]shoppinglist.ShopItem{
			productA.Id: {Product: *productA, MealCount: 1, IsInBasket: false, Quantities: []quantity.Quantity{{Amount: 1, Unit: quantity.Number}}},
			productB.Id: {Product: *productB, MealCount: 1, IsInBasket: false, Quantities: []quantity.Quantity{{Amount: 1, Unit: quantity.Number}}},
		},
		*output.ShoppingList,
	)
}

func (suite *ShoppingListSuite) TestAddingIngredientToMealInShop() {
	productA := suite.addProduct("ing-a", "Ing A", category.AlcoholicDrinks)

	m := suite.addMeal([]meal.Ingredient{*meal.NewIngredient(productA.Id)})

	s, _ := suite.addShop()

	suite.addMealToShop(s, m)

	productB := suite.addProduct("ing-b", "Ing B", category.Bakery)

	suite.addIngredientToMeal(m, productB)

	output := suite.runProjection()

	assert.EqualExportedValues(suite.T(),
		map[string]shoppinglist.ShopItem{
			productA.Id: {Product: *productA, MealCount: 1, IsInBasket: false, Quantities: []quantity.Quantity{{Amount: 1, Unit: quantity.Number}}},
			productB.Id: {Product: *productB, MealCount: 1, IsInBasket: false, Quantities: []quantity.Quantity{{Amount: 1, Unit: quantity.Number}}},
		},
		*output.ShoppingList,
	)
}

func (suite *ShoppingListSuite) TestAddingIngredientToMoreThanOneMealInShop() {
	productA := suite.addProduct("ing-a", "Ing A", category.AlcoholicDrinks)
	productB := suite.addProduct("ing-b", "Ing B", category.Bakery)
	productC := suite.addProduct("ing-c", "Ing C", category.Dairy)

	meal1 := suite.addMeal([]meal.Ingredient{*meal.NewIngredient(productA.Id)})
	meal2 := suite.addMeal([]meal.Ingredient{*meal.NewIngredient(productB.Id)})

	s, _ := suite.addShop()

	suite.addMealToShop(s, meal1)
	suite.addMealToShop(s, meal2)

	suite.addIngredientToMeal(meal1, productC)
	suite.addIngredientToMeal(meal2, productC)

	output := suite.runProjection()

	assert.EqualExportedValues(suite.T(),
		map[string]shoppinglist.ShopItem{
			productA.Id: {Product: *productA, MealCount: 1, IsInBasket: false, Quantities: []quantity.Quantity{{Amount: 1, Unit: quantity.Number}}},
			productB.Id: {Product: *productB, MealCount: 1, IsInBasket: false, Quantities: []quantity.Quantity{{Amount: 1, Unit: quantity.Number}}},
			productC.Id: {Product: *productC, MealCount: 2, IsInBasket: false, Quantities: []quantity.Quantity{{Amount: 1, Unit: quantity.Number}, {Amount: 1, Unit: quantity.Number}}},
		},
		*output.ShoppingList,
	)
}

func (suite *ShoppingListSuite) TestAddingIngredientToMealNotInShop() {
	productA := suite.addProduct("ing-a", "Ing A", category.AlcoholicDrinks)

	meal1 := suite.addMeal([]meal.Ingredient{*meal.NewIngredient(productA.Id)})

	suite.addShop()

	suite.addIngredientToMeal(meal1, productA)

	output := suite.runProjection()

	assert.EqualExportedValues(suite.T(),
		map[string]shoppinglist.ShopItem{},
		*output.ShoppingList,
	)
}

func (suite *ShoppingListSuite) TestAddingIngredientToMealRemovedFromShop() {
	productA := suite.addProduct("ing-a", "Ing A", category.AlcoholicDrinks)

	meal1 := suite.addMeal([]meal.Ingredient{*meal.NewIngredient(productA.Id)})

	s, _ := suite.addShop()

	suite.addMealToShop(s, meal1)
	suite.removeMealFromShop(s, meal1)

	suite.addIngredientToMeal(meal1, productA)

	output := suite.runProjection()

	assert.EqualExportedValues(suite.T(),
		map[string]shoppinglist.ShopItem{},
		*output.ShoppingList,
	)
}

func (suite *ShoppingListSuite) TestRemovingIngredientFromMealInShop() {
	productA := suite.addProduct("ing-a", "Ing A", category.AlcoholicDrinks)
	productB := suite.addProduct("ing-b", "Ing B", category.Bakery)
	productC := suite.addProduct("ing-c", "Ing C", category.Dairy)

	tests := map[string]struct {
		initialMeals [][]meal.Ingredient
		removals     []struct {
			mealIndex  int
			ingredient *product.Product
		}
		expectedOutput map[string]shoppinglist.ShopItem
	}{
		"removing ingredients from a meal": {
			initialMeals: [][]meal.Ingredient{
				{
					*meal.NewIngredient(productA.Id).WithQuantity(10, quantity.Lb),
					*meal.NewIngredient(productB.Id).WithQuantity(5, quantity.Kg),
				},
				{
					*meal.NewIngredient(productB.Id).WithQuantity(15, quantity.Bunch),
					*meal.NewIngredient(productC.Id).WithQuantity(1, quantity.Pack),
				},
			},
			removals: []struct {
				mealIndex  int
				ingredient *product.Product
			}{
				{mealIndex: 0, ingredient: productA},
				{mealIndex: 0, ingredient: productB},
			},
			expectedOutput: map[string]shoppinglist.ShopItem{
				productB.Id: {Product: *productB, MealCount: 1, IsInBasket: false, Quantities: []quantity.Quantity{{Amount: 15, Unit: quantity.Bunch}}},
				productC.Id: {Product: *productC, MealCount: 1, IsInBasket: false, Quantities: []quantity.Quantity{{Amount: 1, Unit: quantity.Pack}}},
			},
		},
		"removing ingredient with same quantity from one of two meals": {
			initialMeals: [][]meal.Ingredient{
				{
					*meal.NewIngredient(productA.Id).WithQuantity(10, quantity.Lb),
				},
				{
					*meal.NewIngredient(productA.Id).WithQuantity(10, quantity.Lb),
				},
			},
			removals: []struct {
				mealIndex  int
				ingredient *product.Product
			}{
				{mealIndex: 0, ingredient: productA},
			},
			expectedOutput: map[string]shoppinglist.ShopItem{
				productA.Id: {Product: *productA, MealCount: 1, IsInBasket: false, Quantities: []quantity.Quantity{{Amount: 10, Unit: quantity.Lb}}},
			},
		},
	}

	for name, t := range tests {
		suite.Run(name, func() {
			var meals []*meal.Meal

			for _, mealIng := range t.initialMeals {
				m := suite.addMeal(mealIng)
				meals = append(meals, m)
			}

			s, _ := suite.addShop()

			for _, m := range meals {
				suite.addMealToShop(s, m)
			}

			for _, removal := range t.removals {
				suite.removeIngredientFromMeal(meals[removal.mealIndex], removal.ingredient)
			}

			output := suite.runProjection()

			assert.EqualExportedValues(suite.T(), t.expectedOutput, *output.ShoppingList)
		})
	}
}

func (suite *ShoppingListSuite) TestRemovingIngredientFromMealBeforeAddingToShop() {
	productA := suite.addProduct("ing-a", "Ing A", category.AlcoholicDrinks)

	m := suite.addMeal([]meal.Ingredient{*meal.NewIngredient(productA.Id)})

	s, _ := suite.addShop()

	suite.addIngredientToMeal(m, productA)
	suite.removeIngredientFromMeal(m, productA)

	suite.addMealToShop(s, m)

	output := suite.runProjection()

	assert.EqualExportedValues(suite.T(),
		map[string]shoppinglist.ShopItem{},
		*output.ShoppingList,
	)
}

func (suite *ShoppingListSuite) TestRemovingIngredientFromMealNotInShop() {
	productA := suite.addProduct("ing-a", "Ing A", category.AlcoholicDrinks)

	meal1 := suite.addMeal([]meal.Ingredient{*meal.NewIngredient(productA.Id)})
	meal2 := suite.addMeal([]meal.Ingredient{*meal.NewIngredient(productA.Id)})

	s, _ := suite.addShop()

	suite.addMealToShop(s, meal1)

	suite.removeIngredientFromMeal(meal2, productA)

	output := suite.runProjection()

	assert.EqualExportedValues(suite.T(),
		map[string]shoppinglist.ShopItem{
			productA.Id: {Product: *productA, MealCount: 1, IsInBasket: false, Quantities: []quantity.Quantity{{Amount: 1, Unit: quantity.Number}}},
		},
		*output.ShoppingList,
	)
}

func (suite *ShoppingListSuite) TestStartingNewShop() {
	productA := suite.addProduct("ing-a", "Ing A", category.AlcoholicDrinks)
	productB := suite.addProduct("ing-b", "Ing B", category.AlcoholicDrinks)

	meal1 := suite.addMeal([]meal.Ingredient{*meal.NewIngredient(productA.Id)})
	meal2 := suite.addMeal([]meal.Ingredient{*meal.NewIngredient(productB.Id)})

	shop1, _ := suite.addShop()

	suite.addMealToShop(shop1, meal1)

	shop2, _ := suite.addShop()

	suite.addMealToShop(shop2, meal2)

	output := suite.runProjection()

	assert.EqualExportedValues(suite.T(),
		map[string]shoppinglist.ShopItem{
			productB.Id: {Product: *productB, MealCount: 1, IsInBasket: false, Quantities: []quantity.Quantity{{Amount: 1, Unit: quantity.Number}}},
		},
		*output.ShoppingList,
	)
}

func (suite *ShoppingListSuite) TestAddingIngredientToBasket() {
	productA := suite.addProduct("ing-a", "Ing A", category.AlcoholicDrinks)

	m := suite.addMeal([]meal.Ingredient{*meal.NewIngredient(productA.Id)})

	s, b := suite.addShop()

	suite.addMealToShop(s, m)

	suite.addItemToBasket(b, productA)

	output := suite.runProjection()

	assert.EqualExportedValues(suite.T(),
		map[string]shoppinglist.ShopItem{
			productA.Id: {Product: *productA, MealCount: 1, IsInBasket: true, Quantities: []quantity.Quantity{{Amount: 1, Unit: quantity.Number}}},
		},
		*output.ShoppingList,
	)
}

func (suite *ShoppingListSuite) TestRemovingIngredientFromBasket() {
	productA := suite.addProduct("ing-a", "Ing A", category.AlcoholicDrinks)

	m := suite.addMeal([]meal.Ingredient{*meal.NewIngredient(productA.Id)})

	s, b := suite.addShop()

	suite.addMealToShop(s, m)

	suite.addItemToBasket(b, productA)

	suite.removeIngredientFromBasket(b, productA)

	output := suite.runProjection()

	assert.EqualExportedValues(suite.T(),
		map[string]shoppinglist.ShopItem{
			productA.Id: {Product: *productA, MealCount: 1, IsInBasket: false, Quantities: []quantity.Quantity{{Amount: 1, Unit: quantity.Number}}},
		},
		*output.ShoppingList,
	)
}

func (suite *ShoppingListSuite) TestReturningShopId() {
	suite.addShop()
	shop2, _ := suite.addShop()

	output := suite.runProjection()

	suite.Assert().Equal(
		shop2.Id,
		*output.ShopId,
	)
}

func (suite *ShoppingListSuite) TestAddingItemToShop() {
	productA := suite.addProduct("ing-a", "Ing A", category.AlcoholicDrinks)

	s, b := suite.addShop()

	suite.addItemToShop(s, &shop.Item{
		ProductId: productA.Id,
		Quantity: quantity.Quantity{
			Amount: 1,
			Unit:   quantity.Number,
		},
	})

	suite.addItemToBasket(b, productA)

	output := suite.runProjection()

	// todo: should items in shop increase meal count?
	assert.EqualExportedValues(suite.T(),
		map[string]shoppinglist.ShopItem{
			productA.Id: {Product: *productA, MealCount: 1, IsInBasket: true, Quantities: []quantity.Quantity{{Amount: 1, Unit: quantity.Number}}},
		},
		*output.ShoppingList,
	)
}

func (suite *ShoppingListSuite) TestAddingItemToShopDirectlyAndViaMeal() {
	productA := suite.addProduct("ing-a", "Ing A", category.AlcoholicDrinks)

	m := suite.addMeal([]meal.Ingredient{*meal.NewIngredient(productA.Id)})

	s, b := suite.addShop()

	suite.addItemToShop(s, &shop.Item{
		ProductId: productA.Id,
		Quantity: quantity.Quantity{
			Amount: 1,
			Unit:   quantity.Number,
		},
	})

	suite.addMealToShop(s, m)
	suite.addItemToBasket(b, productA)

	output := suite.runProjection()

	assert.EqualExportedValues(suite.T(),
		map[string]shoppinglist.ShopItem{
			productA.Id: {Product: *productA, MealCount: 2, IsInBasket: true, Quantities: []quantity.Quantity{{Amount: 1, Unit: quantity.Number}, {Amount: 1, Unit: quantity.Number}}},
		},
		*output.ShoppingList,
	)
}

func (suite *ShoppingListSuite) runProjection() shoppinglist.ShoppingListProjectionOutput {
	projection, output := shoppinglist.CreateShoppingListProjection(suite.es)

	projection.RunToEnd(context.Background())

	return output
}

func (suite *ShoppingListSuite) addIngredientToMeal(m *meal.Meal, product *product.Product) {
	m.AddIngredient(*meal.NewIngredient(product.Id))

	err := suite.mealRepository.Save(m)
	assert.NoError(suite.T(), err)
}

func (suite *ShoppingListSuite) addMeal(ingredients []meal.Ingredient) *meal.Meal {
	id := strconv.Itoa(gofakeit.Number(100, 999))
	m := meal.NewMealBuilder().WithName("Meal " + id).WithId(id).AddIngredients(ingredients).Build()

	err := suite.mealRepository.Save(m)
	assert.NoError(suite.T(), err)

	return m
}

func (suite *ShoppingListSuite) addShop() (*shop.Shop, *basket.Basket) {
	shopId := gofakeit.Number(1, 9999999)
	s, err := shop.NewShop(shopId)
	assert.NoError(suite.T(), err)

	err = suite.shopRepository.Save(s)
	assert.NoError(suite.T(), err)

	b, err := basket.NewBasket(shopId)

	err = suite.basketRepository.Save(b)
	assert.NoError(suite.T(), err)

	return s, b
}

func (suite *ShoppingListSuite) addMealToShop(s *shop.Shop, m *meal.Meal) {
	s.AddMeal(&shop.ShopMeal{MealId: m.Id})
	err := suite.shopRepository.Save(s)
	assert.NoError(suite.T(), err)
}

func (suite *ShoppingListSuite) removeMealFromShop(shop *shop.Shop, meal *meal.Meal) {
	shop.RemoveMeal(meal.Id)
	err := suite.shopRepository.Save(shop)
	assert.NoError(suite.T(), err)
}

func (suite *ShoppingListSuite) addItemToShop(s *shop.Shop, i *shop.Item) {
	s.AddItem(i)
	err := suite.shopRepository.Save(s)
	assert.NoError(suite.T(), err)
}

func (suite *ShoppingListSuite) createRepositories(db *sql.DB) {
	productRepository, err := product.NewSqliteProductRepository(db)
	assert.NoError(suite.T(), err)
	suite.productRepository = productRepository

	shopRepository, err := shop.NewSqliteShopRepository(db)
	assert.NoError(suite.T(), err)
	suite.shopRepository = shopRepository

	mealRepository, err := meal.NewSqliteMealRepository(db)
	assert.NoError(suite.T(), err)
	suite.mealRepository = mealRepository

	basketRepository, err := basket.NewSqliteBasketRepository(db)
	assert.NoError(suite.T(), err)
	suite.basketRepository = basketRepository
}

func (suite *ShoppingListSuite) addProduct(id string, name product.ProductName, category category.CategoryName) *product.Product {
	i, err := product.NewProduct(id, name, category)
	assert.NoError(suite.T(), err)

	err = suite.productRepository.Add(i)
	assert.NoError(suite.T(), err)

	return i
}

func (suite *ShoppingListSuite) removeIngredientFromMeal(meal *meal.Meal, product *product.Product) {
	meal.RemoveIngredient(product.Id)

	err := suite.mealRepository.Save(meal)
	assert.NoError(suite.T(), err)
}

func (suite *ShoppingListSuite) addItemToBasket(b *basket.Basket, product *product.Product) {
	b.AddItem(basket.NewBasketItem(product.Id))

	err := suite.basketRepository.Save(b)

	assert.NoError(suite.T(), err)
}

func (suite *ShoppingListSuite) removeIngredientFromBasket(b *basket.Basket, product *product.Product) {
	b.RemoveItem(product.Id)

	err := suite.basketRepository.Save(b)

	assert.NoError(suite.T(), err)
}

func TestShoppingListSuite(t *testing.T) {
	suite.Run(t, new(ShoppingListSuite))
}
