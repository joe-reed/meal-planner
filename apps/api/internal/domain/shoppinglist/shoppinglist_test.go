package shoppinglist_test

import (
	"context"
	"database/sql"
	"github.com/brianvoe/gofakeit/v7"
	sqlStore "github.com/hallgren/eventsourcing/eventstore/sql"
	"github.com/joe-reed/meal-planner/apps/api/internal/database"
	"github.com/joe-reed/meal-planner/apps/api/internal/domain/basket"
	"github.com/joe-reed/meal-planner/apps/api/internal/domain/category"
	"github.com/joe-reed/meal-planner/apps/api/internal/domain/ingredient"
	"github.com/joe-reed/meal-planner/apps/api/internal/domain/meal"
	"github.com/joe-reed/meal-planner/apps/api/internal/domain/shop"
	"github.com/joe-reed/meal-planner/apps/api/internal/domain/shoppinglist"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	"strconv"
	"testing"
)

type ShoppingListSuite struct {
	suite.Suite
	ingredientRepository *ingredient.EventSourcedIngredientRepository
	shopRepository       *shop.ShopRepository
	mealRepository       *meal.EventSourcedMealRepository
	basketRepository     *basket.BasketRepository
	db                   *sql.DB
	es                   *sqlStore.SQL
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
	ingredientA := suite.addIngredient("ing-a", "Ing A", category.AlcoholicDrinks)

	m := suite.addMeal([]meal.MealIngredient{*meal.NewMealIngredient(ingredientA.Id)})

	s, _ := suite.addShop()

	suite.addMealToShop(s, m)

	output := suite.runProjection()

	assert.EqualExportedValues(suite.T(),
		map[string]shoppinglist.ShopIngredient{
			ingredientA.Id: {Ingredient: *ingredientA, MealCount: 1, IsInBasket: false, Quantities: []meal.Quantity{{Amount: 1, Unit: meal.Number}}},
		},
		*output.ShoppingList,
	)
}

func (suite *ShoppingListSuite) TestAddingTwoMealsWithSameIngredient() {
	ingredientA := suite.addIngredient("ing-a", "Ing A", category.AlcoholicDrinks)
	ingredientB := suite.addIngredient("ing-b", "Ing B", category.Bakery)
	ingredientC := suite.addIngredient("ing-c", "Ing C", category.Dairy)

	meal1 := suite.addMeal([]meal.MealIngredient{
		*meal.NewMealIngredient(ingredientA.Id).WithQuantity(2, meal.Tbsp),
		*meal.NewMealIngredient(ingredientB.Id).WithQuantity(100, meal.Ml),
	})
	meal2 := suite.addMeal([]meal.MealIngredient{
		*meal.NewMealIngredient(ingredientB.Id).WithQuantity(50, meal.Gram),
		*meal.NewMealIngredient(ingredientC.Id).WithQuantity(1, meal.Litre),
	})

	s, _ := suite.addShop()

	suite.addMealToShop(s, meal1)
	suite.addMealToShop(s, meal2)

	output := suite.runProjection()

	assert.EqualExportedValues(suite.T(),
		map[string]shoppinglist.ShopIngredient{
			ingredientA.Id: {Ingredient: *ingredientA, MealCount: 1, IsInBasket: false, Quantities: []meal.Quantity{{Amount: 2, Unit: meal.Tbsp}}},
			ingredientB.Id: {Ingredient: *ingredientB, MealCount: 2, IsInBasket: false, Quantities: []meal.Quantity{{Amount: 100, Unit: meal.Ml}, {Amount: 50, Unit: meal.Gram}}},
			ingredientC.Id: {Ingredient: *ingredientC, MealCount: 1, IsInBasket: false, Quantities: []meal.Quantity{{Amount: 1, Unit: meal.Litre}}},
		},
		*output.ShoppingList,
	)
}

func (suite *ShoppingListSuite) TestRemovingMeal() {
	ingredientA := suite.addIngredient("ing-a", "Ing A", category.AlcoholicDrinks)
	ingredientB := suite.addIngredient("ing-b", "Ing B", category.Bakery)
	ingredientC := suite.addIngredient("ing-c", "Ing C", category.Dairy)

	meal1 := suite.addMeal([]meal.MealIngredient{*meal.NewMealIngredient(ingredientA.Id), *meal.NewMealIngredient(ingredientB.Id)})
	meal2 := suite.addMeal([]meal.MealIngredient{*meal.NewMealIngredient(ingredientB.Id), *meal.NewMealIngredient(ingredientC.Id)})

	s, _ := suite.addShop()

	suite.addMealToShop(s, meal1)
	suite.addMealToShop(s, meal2)
	suite.removeMealFromShop(s, meal2)

	output := suite.runProjection()

	assert.EqualExportedValues(suite.T(),
		map[string]shoppinglist.ShopIngredient{
			ingredientA.Id: {Ingredient: *ingredientA, MealCount: 1, IsInBasket: false, Quantities: []meal.Quantity{{Amount: 1, Unit: meal.Number}}},
			ingredientB.Id: {Ingredient: *ingredientB, MealCount: 1, IsInBasket: false, Quantities: []meal.Quantity{{Amount: 1, Unit: meal.Number}}},
		},
		*output.ShoppingList,
	)
}

func (suite *ShoppingListSuite) TestAddingIngredientToMealBeforeAddingToShop() {
	ingredientA := suite.addIngredient("ing-a", "Ing A", category.AlcoholicDrinks)

	m := suite.addMeal([]meal.MealIngredient{*meal.NewMealIngredient(ingredientA.Id)})

	s, _ := suite.addShop()

	ingredientB := suite.addIngredient("ing-b", "Ing B", category.Bakery)

	suite.addIngredientToMeal(m, ingredientB)

	suite.addMealToShop(s, m)

	output := suite.runProjection()

	assert.EqualExportedValues(suite.T(),
		map[string]shoppinglist.ShopIngredient{
			ingredientA.Id: {Ingredient: *ingredientA, MealCount: 1, IsInBasket: false, Quantities: []meal.Quantity{{Amount: 1, Unit: meal.Number}}},
			ingredientB.Id: {Ingredient: *ingredientB, MealCount: 1, IsInBasket: false, Quantities: []meal.Quantity{{Amount: 1, Unit: meal.Number}}},
		},
		*output.ShoppingList,
	)
}

func (suite *ShoppingListSuite) TestAddingIngredientToMealInShop() {
	ingredientA := suite.addIngredient("ing-a", "Ing A", category.AlcoholicDrinks)

	m := suite.addMeal([]meal.MealIngredient{*meal.NewMealIngredient(ingredientA.Id)})

	s, _ := suite.addShop()

	suite.addMealToShop(s, m)

	ingredientB := suite.addIngredient("ing-b", "Ing B", category.Bakery)

	suite.addIngredientToMeal(m, ingredientB)

	output := suite.runProjection()

	assert.EqualExportedValues(suite.T(),
		map[string]shoppinglist.ShopIngredient{
			ingredientA.Id: {Ingredient: *ingredientA, MealCount: 1, IsInBasket: false, Quantities: []meal.Quantity{{Amount: 1, Unit: meal.Number}}},
			ingredientB.Id: {Ingredient: *ingredientB, MealCount: 1, IsInBasket: false, Quantities: []meal.Quantity{{Amount: 1, Unit: meal.Number}}},
		},
		*output.ShoppingList,
	)
}

func (suite *ShoppingListSuite) TestAddingIngredientToMoreThanOneMealInShop() {
	ingredientA := suite.addIngredient("ing-a", "Ing A", category.AlcoholicDrinks)
	ingredientB := suite.addIngredient("ing-b", "Ing B", category.Bakery)
	ingredientC := suite.addIngredient("ing-c", "Ing C", category.Dairy)

	meal1 := suite.addMeal([]meal.MealIngredient{*meal.NewMealIngredient(ingredientA.Id)})
	meal2 := suite.addMeal([]meal.MealIngredient{*meal.NewMealIngredient(ingredientB.Id)})

	s, _ := suite.addShop()

	suite.addMealToShop(s, meal1)
	suite.addMealToShop(s, meal2)

	suite.addIngredientToMeal(meal1, ingredientC)
	suite.addIngredientToMeal(meal2, ingredientC)

	output := suite.runProjection()

	assert.EqualExportedValues(suite.T(),
		map[string]shoppinglist.ShopIngredient{
			ingredientA.Id: {Ingredient: *ingredientA, MealCount: 1, IsInBasket: false, Quantities: []meal.Quantity{{Amount: 1, Unit: meal.Number}}},
			ingredientB.Id: {Ingredient: *ingredientB, MealCount: 1, IsInBasket: false, Quantities: []meal.Quantity{{Amount: 1, Unit: meal.Number}}},
			ingredientC.Id: {Ingredient: *ingredientC, MealCount: 2, IsInBasket: false, Quantities: []meal.Quantity{{Amount: 1, Unit: meal.Number}, {Amount: 1, Unit: meal.Number}}},
		},
		*output.ShoppingList,
	)
}

func (suite *ShoppingListSuite) TestAddingIngredientToMealNotInShop() {
	ingredientA := suite.addIngredient("ing-a", "Ing A", category.AlcoholicDrinks)

	meal1 := suite.addMeal([]meal.MealIngredient{*meal.NewMealIngredient(ingredientA.Id)})

	suite.addShop()

	suite.addIngredientToMeal(meal1, ingredientA)

	output := suite.runProjection()

	assert.EqualExportedValues(suite.T(),
		map[string]shoppinglist.ShopIngredient{},
		*output.ShoppingList,
	)
}

func (suite *ShoppingListSuite) TestAddingIngredientToMealRemovedFromShop() {
	ingredientA := suite.addIngredient("ing-a", "Ing A", category.AlcoholicDrinks)

	meal1 := suite.addMeal([]meal.MealIngredient{*meal.NewMealIngredient(ingredientA.Id)})

	s, _ := suite.addShop()

	suite.addMealToShop(s, meal1)
	suite.removeMealFromShop(s, meal1)

	suite.addIngredientToMeal(meal1, ingredientA)

	output := suite.runProjection()

	assert.EqualExportedValues(suite.T(),
		map[string]shoppinglist.ShopIngredient{},
		*output.ShoppingList,
	)
}

func (suite *ShoppingListSuite) TestRemovingIngredientFromMealInShop() {
	ingA := suite.addIngredient("ing-a", "Ing A", category.AlcoholicDrinks)
	ingB := suite.addIngredient("ing-b", "Ing B", category.Bakery)
	ingC := suite.addIngredient("ing-c", "Ing C", category.Dairy)

	tests := map[string]struct {
		initialMeals [][]meal.MealIngredient
		removals     []struct {
			mealIndex  int
			ingredient *ingredient.Ingredient
		}
		expectedOutput map[string]shoppinglist.ShopIngredient
	}{
		"removing ingredients from a meal": {
			initialMeals: [][]meal.MealIngredient{
				{
					*meal.NewMealIngredient(ingA.Id).WithQuantity(10, meal.Lb),
					*meal.NewMealIngredient(ingB.Id).WithQuantity(5, meal.Kg),
				},
				{
					*meal.NewMealIngredient(ingB.Id).WithQuantity(15, meal.Bunch),
					*meal.NewMealIngredient(ingC.Id).WithQuantity(1, meal.Pack),
				},
			},
			removals: []struct {
				mealIndex  int
				ingredient *ingredient.Ingredient
			}{
				{mealIndex: 0, ingredient: ingA},
				{mealIndex: 0, ingredient: ingB},
			},
			expectedOutput: map[string]shoppinglist.ShopIngredient{
				ingB.Id: {Ingredient: *ingB, MealCount: 1, IsInBasket: false, Quantities: []meal.Quantity{{Amount: 15, Unit: meal.Bunch}}},
				ingC.Id: {Ingredient: *ingC, MealCount: 1, IsInBasket: false, Quantities: []meal.Quantity{{Amount: 1, Unit: meal.Pack}}},
			},
		},
		"removing ingredient with same quantity from one of two meals": {
			initialMeals: [][]meal.MealIngredient{
				{
					*meal.NewMealIngredient(ingA.Id).WithQuantity(10, meal.Lb),
				},
				{
					*meal.NewMealIngredient(ingA.Id).WithQuantity(10, meal.Lb),
				},
			},
			removals: []struct {
				mealIndex  int
				ingredient *ingredient.Ingredient
			}{
				{mealIndex: 0, ingredient: ingA},
			},
			expectedOutput: map[string]shoppinglist.ShopIngredient{
				ingA.Id: {Ingredient: *ingA, MealCount: 1, IsInBasket: false, Quantities: []meal.Quantity{{Amount: 10, Unit: meal.Lb}}},
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
	ingredientA := suite.addIngredient("ing-a", "Ing A", category.AlcoholicDrinks)

	m := suite.addMeal([]meal.MealIngredient{*meal.NewMealIngredient(ingredientA.Id)})

	s, _ := suite.addShop()

	suite.addIngredientToMeal(m, ingredientA)
	suite.removeIngredientFromMeal(m, ingredientA)

	suite.addMealToShop(s, m)

	output := suite.runProjection()

	assert.EqualExportedValues(suite.T(),
		map[string]shoppinglist.ShopIngredient{},
		*output.ShoppingList,
	)
}

func (suite *ShoppingListSuite) TestRemovingIngredientFromMealNotInShop() {
	ingredientA := suite.addIngredient("ing-a", "Ing A", category.AlcoholicDrinks)

	meal1 := suite.addMeal([]meal.MealIngredient{*meal.NewMealIngredient(ingredientA.Id)})
	meal2 := suite.addMeal([]meal.MealIngredient{*meal.NewMealIngredient(ingredientA.Id)})

	s, _ := suite.addShop()

	suite.addMealToShop(s, meal1)

	suite.removeIngredientFromMeal(meal2, ingredientA)

	output := suite.runProjection()

	assert.EqualExportedValues(suite.T(),
		map[string]shoppinglist.ShopIngredient{
			ingredientA.Id: {Ingredient: *ingredientA, MealCount: 1, IsInBasket: false, Quantities: []meal.Quantity{{Amount: 1, Unit: meal.Number}}},
		},
		*output.ShoppingList,
	)
}

func (suite *ShoppingListSuite) TestStartingNewShop() {
	ingredientA := suite.addIngredient("ing-a", "Ing A", category.AlcoholicDrinks)
	ingredientB := suite.addIngredient("ing-b", "Ing B", category.AlcoholicDrinks)

	meal1 := suite.addMeal([]meal.MealIngredient{*meal.NewMealIngredient(ingredientA.Id)})
	meal2 := suite.addMeal([]meal.MealIngredient{*meal.NewMealIngredient(ingredientB.Id)})

	shop1, _ := suite.addShop()

	suite.addMealToShop(shop1, meal1)

	shop2, _ := suite.addShop()

	suite.addMealToShop(shop2, meal2)

	output := suite.runProjection()

	assert.EqualExportedValues(suite.T(),
		map[string]shoppinglist.ShopIngredient{
			ingredientB.Id: {Ingredient: *ingredientB, MealCount: 1, IsInBasket: false, Quantities: []meal.Quantity{{Amount: 1, Unit: meal.Number}}},
		},
		*output.ShoppingList,
	)
}

func (suite *ShoppingListSuite) TestAddingIngredientToBasket() {
	ingredientA := suite.addIngredient("ing-a", "Ing A", category.AlcoholicDrinks)

	m := suite.addMeal([]meal.MealIngredient{*meal.NewMealIngredient(ingredientA.Id)})

	s, b := suite.addShop()

	suite.addMealToShop(s, m)

	suite.addIngredientToBasket(b, ingredientA)

	output := suite.runProjection()

	assert.EqualExportedValues(suite.T(),
		map[string]shoppinglist.ShopIngredient{
			ingredientA.Id: {Ingredient: *ingredientA, MealCount: 1, IsInBasket: true, Quantities: []meal.Quantity{{Amount: 1, Unit: meal.Number}}},
		},
		*output.ShoppingList,
	)
}

func (suite *ShoppingListSuite) TestRemovingIngredientFromBasket() {
	ingredientA := suite.addIngredient("ing-a", "Ing A", category.AlcoholicDrinks)

	m := suite.addMeal([]meal.MealIngredient{*meal.NewMealIngredient(ingredientA.Id)})

	s, b := suite.addShop()

	suite.addMealToShop(s, m)

	suite.addIngredientToBasket(b, ingredientA)

	suite.removeIngredientFromBasket(b, ingredientA)

	output := suite.runProjection()

	assert.EqualExportedValues(suite.T(),
		map[string]shoppinglist.ShopIngredient{
			ingredientA.Id: {Ingredient: *ingredientA, MealCount: 1, IsInBasket: false, Quantities: []meal.Quantity{{Amount: 1, Unit: meal.Number}}},
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

func (suite *ShoppingListSuite) runProjection() shoppinglist.ShoppingListProjectionOutput {
	projection, output := shoppinglist.CreateShoppingListProjection(suite.es)

	projection.RunToEnd(context.Background())

	return output
}

func (suite *ShoppingListSuite) addIngredientToMeal(m *meal.Meal, ingredient *ingredient.Ingredient) {
	m.AddIngredient(*meal.NewMealIngredient(ingredient.Id))

	err := suite.mealRepository.Save(m)
	assert.NoError(suite.T(), err)
}

func (suite *ShoppingListSuite) addMeal(mealIngredients []meal.MealIngredient) *meal.Meal {
	id := strconv.Itoa(gofakeit.Number(100, 999))
	m := meal.NewMealBuilder().WithName("Meal " + id).WithId(id).AddIngredients(mealIngredients).Build()

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

func (suite *ShoppingListSuite) createRepositories(db *sql.DB) {
	ingredientRepository, err := ingredient.NewSqliteIngredientRepository(db)
	assert.NoError(suite.T(), err)
	suite.ingredientRepository = ingredientRepository

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

func (suite *ShoppingListSuite) addIngredient(id string, name ingredient.IngredientName, category category.CategoryName) *ingredient.Ingredient {
	i, err := ingredient.NewIngredient(id, name, category)
	assert.NoError(suite.T(), err)

	err = suite.ingredientRepository.Add(i)
	assert.NoError(suite.T(), err)

	return i
}

func (suite *ShoppingListSuite) removeIngredientFromMeal(meal *meal.Meal, ingredient *ingredient.Ingredient) {
	meal.RemoveIngredient(ingredient.Id)

	err := suite.mealRepository.Save(meal)
	assert.NoError(suite.T(), err)
}

func (suite *ShoppingListSuite) addIngredientToBasket(b *basket.Basket, ingredient *ingredient.Ingredient) {
	b.AddItem(basket.NewBasketItem(ingredient.Id))

	err := suite.basketRepository.Save(b)

	assert.NoError(suite.T(), err)
}

func (suite *ShoppingListSuite) removeIngredientFromBasket(b *basket.Basket, ingredient *ingredient.Ingredient) {
	b.RemoveItem(ingredient.Id)

	err := suite.basketRepository.Save(b)

	assert.NoError(suite.T(), err)
}

func TestShoppingListSuite(t *testing.T) {
	suite.Run(t, new(ShoppingListSuite))
}
