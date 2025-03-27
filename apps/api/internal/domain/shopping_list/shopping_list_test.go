package shopping_list_test

import (
	"context"
	"database/sql"
	"github.com/brianvoe/gofakeit/v7"
	sqlStore "github.com/hallgren/eventsourcing/eventstore/sql"
	"github.com/joe-reed/meal-planner/apps/api/internal/database"
	"github.com/joe-reed/meal-planner/apps/api/internal/domain/basket"
	"github.com/joe-reed/meal-planner/apps/api/internal/domain/categories"
	"github.com/joe-reed/meal-planner/apps/api/internal/domain/ingredients"
	"github.com/joe-reed/meal-planner/apps/api/internal/domain/meals"
	"github.com/joe-reed/meal-planner/apps/api/internal/domain/shopping_list"
	"github.com/joe-reed/meal-planner/apps/api/internal/domain/shops"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	"strconv"
	"testing"
)

type ShoppingListSuite struct {
	suite.Suite
	ingredientRepository *ingredients.IngredientRepository
	shopRepository       *shops.ShopRepository
	mealRepository       *meals.MealRepository
	basketRepository     *basket.BasketRepository
	db                   *sql.DB
	es                   *sqlStore.SQL
}

func (s *ShoppingListSuite) SetupTest() {
	db, err := database.CreateDatabase(":memory:")
	assert.NoError(s.T(), err)

	s.db = db
	s.es = sqlStore.Open(db)

	s.createRepositories(db)
}

func (s *ShoppingListSuite) TearDownTest() {
	err := s.db.Close()

	if err != nil {
		s.T().Fatal(err)
	}
}

func (s *ShoppingListSuite) TestAddingMealToShop() {
	ingredientA := s.addIngredient("ing-a", "Ing A", categories.AlcoholicDrinks)

	meal := s.addMeal([]meals.MealIngredient{*meals.NewMealIngredient(ingredientA.Id)})

	shop, _ := s.addShop()

	s.addMealToShop(shop, meal)

	output := s.runProjection()

	assert.EqualExportedValues(s.T(),
		map[string]shopping_list.ShopIngredient{
			ingredientA.Id: {Ingredient: *ingredientA, MealCount: 1, IsInBasket: false, Quantities: []meals.Quantity{{1, meals.Number}}},
		},
		*output.ShoppingList,
	)
}

func (s *ShoppingListSuite) TestAddingTwoMealsWithSameIngredient() {
	ingredientA := s.addIngredient("ing-a", "Ing A", categories.AlcoholicDrinks)
	ingredientB := s.addIngredient("ing-b", "Ing B", categories.Bakery)
	ingredientC := s.addIngredient("ing-c", "Ing C", categories.Dairy)

	meal1 := s.addMeal([]meals.MealIngredient{
		*meals.NewMealIngredient(ingredientA.Id).WithQuantity(2, meals.Tbsp),
		*meals.NewMealIngredient(ingredientB.Id).WithQuantity(100, meals.Ml),
	})
	meal2 := s.addMeal([]meals.MealIngredient{
		*meals.NewMealIngredient(ingredientB.Id).WithQuantity(50, meals.Gram),
		*meals.NewMealIngredient(ingredientC.Id).WithQuantity(1, meals.Litre),
	})

	shop, _ := s.addShop()

	s.addMealToShop(shop, meal1)
	s.addMealToShop(shop, meal2)

	output := s.runProjection()

	assert.EqualExportedValues(s.T(),
		map[string]shopping_list.ShopIngredient{
			ingredientA.Id: {Ingredient: *ingredientA, MealCount: 1, IsInBasket: false, Quantities: []meals.Quantity{{2, meals.Tbsp}}},
			ingredientB.Id: {Ingredient: *ingredientB, MealCount: 2, IsInBasket: false, Quantities: []meals.Quantity{{100, meals.Ml}, {50, meals.Gram}}},
			ingredientC.Id: {Ingredient: *ingredientC, MealCount: 1, IsInBasket: false, Quantities: []meals.Quantity{{1, meals.Litre}}},
		},
		*output.ShoppingList,
	)
}

func (s *ShoppingListSuite) TestRemovingMeal() {
	ingredientA := s.addIngredient("ing-a", "Ing A", categories.AlcoholicDrinks)
	ingredientB := s.addIngredient("ing-b", "Ing B", categories.Bakery)
	ingredientC := s.addIngredient("ing-c", "Ing C", categories.Dairy)

	meal1 := s.addMeal([]meals.MealIngredient{*meals.NewMealIngredient(ingredientA.Id), *meals.NewMealIngredient(ingredientB.Id)})
	meal2 := s.addMeal([]meals.MealIngredient{*meals.NewMealIngredient(ingredientB.Id), *meals.NewMealIngredient(ingredientC.Id)})

	shop, _ := s.addShop()

	s.addMealToShop(shop, meal1)
	s.addMealToShop(shop, meal2)
	s.removeMealFromShop(shop, meal2)

	output := s.runProjection()

	assert.EqualExportedValues(s.T(),
		map[string]shopping_list.ShopIngredient{
			ingredientA.Id: {Ingredient: *ingredientA, MealCount: 1, IsInBasket: false, Quantities: []meals.Quantity{{1, meals.Number}}},
			ingredientB.Id: {Ingredient: *ingredientB, MealCount: 1, IsInBasket: false, Quantities: []meals.Quantity{{1, meals.Number}}},
		},
		*output.ShoppingList,
	)
}

func (s *ShoppingListSuite) TestAddingIngredientToMealBeforeAddingToShop() {
	ingredientA := s.addIngredient("ing-a", "Ing A", categories.AlcoholicDrinks)

	meal := s.addMeal([]meals.MealIngredient{*meals.NewMealIngredient(ingredientA.Id)})

	shop, _ := s.addShop()

	ingredientB := s.addIngredient("ing-b", "Ing B", categories.Bakery)

	s.addIngredientToMeal(meal, ingredientB)

	s.addMealToShop(shop, meal)

	output := s.runProjection()

	assert.EqualExportedValues(s.T(),
		map[string]shopping_list.ShopIngredient{
			ingredientA.Id: {Ingredient: *ingredientA, MealCount: 1, IsInBasket: false, Quantities: []meals.Quantity{{1, meals.Number}}},
			ingredientB.Id: {Ingredient: *ingredientB, MealCount: 1, IsInBasket: false, Quantities: []meals.Quantity{{1, meals.Number}}},
		},
		*output.ShoppingList,
	)
}

func (s *ShoppingListSuite) TestAddingIngredientToMealInShop() {
	ingredientA := s.addIngredient("ing-a", "Ing A", categories.AlcoholicDrinks)

	meal := s.addMeal([]meals.MealIngredient{*meals.NewMealIngredient(ingredientA.Id)})

	shop, _ := s.addShop()

	s.addMealToShop(shop, meal)

	ingredientB := s.addIngredient("ing-b", "Ing B", categories.Bakery)

	s.addIngredientToMeal(meal, ingredientB)

	output := s.runProjection()

	assert.EqualExportedValues(s.T(),
		map[string]shopping_list.ShopIngredient{
			ingredientA.Id: {Ingredient: *ingredientA, MealCount: 1, IsInBasket: false, Quantities: []meals.Quantity{{1, meals.Number}}},
			ingredientB.Id: {Ingredient: *ingredientB, MealCount: 1, IsInBasket: false, Quantities: []meals.Quantity{{1, meals.Number}}},
		},
		*output.ShoppingList,
	)
}

func (s *ShoppingListSuite) TestAddingIngredientToMoreThanOneMealInShop() {
	ingredientA := s.addIngredient("ing-a", "Ing A", categories.AlcoholicDrinks)
	ingredientB := s.addIngredient("ing-b", "Ing B", categories.Bakery)
	ingredientC := s.addIngredient("ing-c", "Ing C", categories.Dairy)

	meal1 := s.addMeal([]meals.MealIngredient{*meals.NewMealIngredient(ingredientA.Id)})
	meal2 := s.addMeal([]meals.MealIngredient{*meals.NewMealIngredient(ingredientB.Id)})

	shop, _ := s.addShop()

	s.addMealToShop(shop, meal1)
	s.addMealToShop(shop, meal2)

	s.addIngredientToMeal(meal1, ingredientC)
	s.addIngredientToMeal(meal2, ingredientC)

	output := s.runProjection()

	assert.EqualExportedValues(s.T(),
		map[string]shopping_list.ShopIngredient{
			ingredientA.Id: {Ingredient: *ingredientA, MealCount: 1, IsInBasket: false, Quantities: []meals.Quantity{{1, meals.Number}}},
			ingredientB.Id: {Ingredient: *ingredientB, MealCount: 1, IsInBasket: false, Quantities: []meals.Quantity{{1, meals.Number}}},
			ingredientC.Id: {Ingredient: *ingredientC, MealCount: 2, IsInBasket: false, Quantities: []meals.Quantity{{1, meals.Number}, {1, meals.Number}}},
		},
		*output.ShoppingList,
	)
}

func (s *ShoppingListSuite) TestAddingIngredientToMealNotInShop() {
	ingredientA := s.addIngredient("ing-a", "Ing A", categories.AlcoholicDrinks)

	meal1 := s.addMeal([]meals.MealIngredient{*meals.NewMealIngredient(ingredientA.Id)})

	s.addShop()

	s.addIngredientToMeal(meal1, ingredientA)

	output := s.runProjection()

	assert.EqualExportedValues(s.T(),
		map[string]shopping_list.ShopIngredient{},
		*output.ShoppingList,
	)
}

func (s *ShoppingListSuite) TestAddingIngredientToMealRemovedFromShop() {
	ingredientA := s.addIngredient("ing-a", "Ing A", categories.AlcoholicDrinks)

	meal1 := s.addMeal([]meals.MealIngredient{*meals.NewMealIngredient(ingredientA.Id)})

	shop, _ := s.addShop()

	s.addMealToShop(shop, meal1)
	s.removeMealFromShop(shop, meal1)

	s.addIngredientToMeal(meal1, ingredientA)

	output := s.runProjection()

	assert.EqualExportedValues(s.T(),
		map[string]shopping_list.ShopIngredient{},
		*output.ShoppingList,
	)
}

func (s *ShoppingListSuite) TestRemovingIngredientFromMealInShop() {
	ingA := s.addIngredient("ing-a", "Ing A", categories.AlcoholicDrinks)
	ingB := s.addIngredient("ing-b", "Ing B", categories.Bakery)
	ingC := s.addIngredient("ing-c", "Ing C", categories.Dairy)

	tests := map[string]struct {
		initialMeals [][]meals.MealIngredient
		removals     []struct {
			mealIndex  int
			ingredient *ingredients.Ingredient
		}
		expectedOutput map[string]shopping_list.ShopIngredient
	}{
		"removing ingredients from a meal": {
			initialMeals: [][]meals.MealIngredient{
				{
					*meals.NewMealIngredient(ingA.Id).WithQuantity(10, meals.Lb),
					*meals.NewMealIngredient(ingB.Id).WithQuantity(5, meals.Kg),
				},
				{
					*meals.NewMealIngredient(ingB.Id).WithQuantity(15, meals.Bunch),
					*meals.NewMealIngredient(ingC.Id).WithQuantity(1, meals.Pack),
				},
			},
			removals: []struct {
				mealIndex  int
				ingredient *ingredients.Ingredient
			}{
				{mealIndex: 0, ingredient: ingA},
				{mealIndex: 0, ingredient: ingB},
			},
			expectedOutput: map[string]shopping_list.ShopIngredient{
				ingB.Id: {Ingredient: *ingB, MealCount: 1, IsInBasket: false, Quantities: []meals.Quantity{{15, meals.Bunch}}},
				ingC.Id: {Ingredient: *ingC, MealCount: 1, IsInBasket: false, Quantities: []meals.Quantity{{1, meals.Pack}}},
			},
		},
		"removing ingredient with same quantity from one of twon meals": {
			initialMeals: [][]meals.MealIngredient{
				{
					*meals.NewMealIngredient(ingA.Id).WithQuantity(10, meals.Lb),
				},
				{
					*meals.NewMealIngredient(ingA.Id).WithQuantity(10, meals.Lb),
				},
			},
			removals: []struct {
				mealIndex  int
				ingredient *ingredients.Ingredient
			}{
				{mealIndex: 0, ingredient: ingA},
			},
			expectedOutput: map[string]shopping_list.ShopIngredient{
				ingA.Id: {Ingredient: *ingA, MealCount: 1, IsInBasket: false, Quantities: []meals.Quantity{{10, meals.Lb}}},
			},
		},
	}

	for name, t := range tests {
		s.Run(name, func() {
			var mealsList []*meals.Meal

			for _, mealIng := range t.initialMeals {
				meal := s.addMeal(mealIng)
				mealsList = append(mealsList, meal)
			}

			shop, _ := s.addShop()

			for _, meal := range mealsList {
				s.addMealToShop(shop, meal)
			}

			for _, removal := range t.removals {
				s.removeIngredientFromMeal(mealsList[removal.mealIndex], removal.ingredient)
			}

			output := s.runProjection()

			assert.EqualExportedValues(s.T(), t.expectedOutput, *output.ShoppingList)
		})
	}
}

func (s *ShoppingListSuite) TestRemovingIngredientFromMealBeforeAddingToShop() {
	ingredientA := s.addIngredient("ing-a", "Ing A", categories.AlcoholicDrinks)

	meal := s.addMeal([]meals.MealIngredient{*meals.NewMealIngredient(ingredientA.Id)})

	shop, _ := s.addShop()

	s.addIngredientToMeal(meal, ingredientA)
	s.removeIngredientFromMeal(meal, ingredientA)

	s.addMealToShop(shop, meal)

	output := s.runProjection()

	assert.EqualExportedValues(s.T(),
		map[string]shopping_list.ShopIngredient{},
		*output.ShoppingList,
	)
}

func (s *ShoppingListSuite) TestRemovingIngredientFromMealNotInShop() {
	ingredientA := s.addIngredient("ing-a", "Ing A", categories.AlcoholicDrinks)

	meal1 := s.addMeal([]meals.MealIngredient{*meals.NewMealIngredient(ingredientA.Id)})
	meal2 := s.addMeal([]meals.MealIngredient{*meals.NewMealIngredient(ingredientA.Id)})

	shop, _ := s.addShop()

	s.addMealToShop(shop, meal1)

	s.removeIngredientFromMeal(meal2, ingredientA)

	output := s.runProjection()

	assert.EqualExportedValues(s.T(),
		map[string]shopping_list.ShopIngredient{
			ingredientA.Id: {Ingredient: *ingredientA, MealCount: 1, IsInBasket: false, Quantities: []meals.Quantity{{1, meals.Number}}},
		},
		*output.ShoppingList,
	)
}

func (s *ShoppingListSuite) TestStartingNewShop() {
	ingredientA := s.addIngredient("ing-a", "Ing A", categories.AlcoholicDrinks)
	ingredientB := s.addIngredient("ing-b", "Ing B", categories.AlcoholicDrinks)

	meal1 := s.addMeal([]meals.MealIngredient{*meals.NewMealIngredient(ingredientA.Id)})
	meal2 := s.addMeal([]meals.MealIngredient{*meals.NewMealIngredient(ingredientB.Id)})

	shop1, _ := s.addShop()

	s.addMealToShop(shop1, meal1)

	shop2, _ := s.addShop()

	s.addMealToShop(shop2, meal2)

	output := s.runProjection()

	assert.EqualExportedValues(s.T(),
		map[string]shopping_list.ShopIngredient{
			ingredientB.Id: {Ingredient: *ingredientB, MealCount: 1, IsInBasket: false, Quantities: []meals.Quantity{{1, meals.Number}}},
		},
		*output.ShoppingList,
	)
}

func (s *ShoppingListSuite) TestAddingIngredientToBasket() {
	ingredientA := s.addIngredient("ing-a", "Ing A", categories.AlcoholicDrinks)

	meal := s.addMeal([]meals.MealIngredient{*meals.NewMealIngredient(ingredientA.Id)})

	shop, b := s.addShop()

	s.addMealToShop(shop, meal)

	s.addIngredientToBasket(b, ingredientA)

	output := s.runProjection()

	assert.EqualExportedValues(s.T(),
		map[string]shopping_list.ShopIngredient{
			ingredientA.Id: {Ingredient: *ingredientA, MealCount: 1, IsInBasket: true, Quantities: []meals.Quantity{{1, meals.Number}}},
		},
		*output.ShoppingList,
	)
}

func (s *ShoppingListSuite) TestRemovingIngredientFromBasket() {
	ingredientA := s.addIngredient("ing-a", "Ing A", categories.AlcoholicDrinks)

	meal := s.addMeal([]meals.MealIngredient{*meals.NewMealIngredient(ingredientA.Id)})

	shop, b := s.addShop()

	s.addMealToShop(shop, meal)

	s.addIngredientToBasket(b, ingredientA)

	s.removeIngredientFromBasket(b, ingredientA)

	output := s.runProjection()

	assert.EqualExportedValues(s.T(),
		map[string]shopping_list.ShopIngredient{
			ingredientA.Id: {Ingredient: *ingredientA, MealCount: 1, IsInBasket: false, Quantities: []meals.Quantity{{1, meals.Number}}},
		},
		*output.ShoppingList,
	)
}

func (s *ShoppingListSuite) TestReturningShopId() {
	s.addShop()
	shop2, _ := s.addShop()

	output := s.runProjection()

	s.Assert().Equal(
		shop2.Id,
		*output.ShopId,
	)
}

func (s *ShoppingListSuite) runProjection() shopping_list.ShoppingListProjectionOutput {
	projection, output := shopping_list.CreateShoppingListProjection(s.es)

	projection.RunToEnd(context.Background())

	return output
}

func (s *ShoppingListSuite) addIngredientToMeal(meal *meals.Meal, ingredient *ingredients.Ingredient) {
	meal.AddIngredient(*meals.NewMealIngredient(ingredient.Id))

	err := s.mealRepository.Save(meal)
	assert.NoError(s.T(), err)
}

func (s *ShoppingListSuite) addMeal(mealIngredients []meals.MealIngredient) *meals.Meal {
	id := strconv.Itoa(gofakeit.Number(100, 999))
	meal := meals.NewMealBuilder().WithName("Meal " + id).WithId(id).AddIngredients(mealIngredients).Build()

	err := s.mealRepository.Save(meal)
	assert.NoError(s.T(), err)

	return meal
}

func (s *ShoppingListSuite) addShop() (*shops.Shop, *basket.Basket) {
	shopId := gofakeit.Number(1, 9999999)
	shop, err := shops.NewShop(shopId)
	assert.NoError(s.T(), err)

	err = s.shopRepository.Save(shop)
	assert.NoError(s.T(), err)

	b, err := basket.NewBasket(shopId)

	err = s.basketRepository.Save(b)
	assert.NoError(s.T(), err)

	return shop, b
}

func (s *ShoppingListSuite) addMealToShop(shop *shops.Shop, meal *meals.Meal) {
	shop.AddMeal(&shops.ShopMeal{MealId: meal.Id})
	err := s.shopRepository.Save(shop)
	assert.NoError(s.T(), err)
}

func (s *ShoppingListSuite) removeMealFromShop(shop *shops.Shop, meal *meals.Meal) {
	shop.RemoveMeal(meal.Id)
	err := s.shopRepository.Save(shop)
	assert.NoError(s.T(), err)
}

func (s *ShoppingListSuite) createRepositories(db *sql.DB) {
	ingredientRepository, err := ingredients.NewSqliteIngredientRepository(db)
	assert.NoError(s.T(), err)
	s.ingredientRepository = ingredientRepository

	shopRepository, err := shops.NewSqliteShopRepository(db)
	assert.NoError(s.T(), err)
	s.shopRepository = shopRepository

	mealRepository, err := meals.NewSqliteMealRepository(db)
	assert.NoError(s.T(), err)
	s.mealRepository = mealRepository

	basketRepository, err := basket.NewSqliteBasketRepository(db)
	assert.NoError(s.T(), err)
	s.basketRepository = basketRepository
}

func (s *ShoppingListSuite) addIngredient(id, name string, category categories.Category) *ingredients.Ingredient {
	ingredient, err := ingredients.NewIngredient(id, name, category)
	assert.NoError(s.T(), err)

	err = s.ingredientRepository.Add(ingredient)
	assert.NoError(s.T(), err)

	return ingredient
}

func (s *ShoppingListSuite) removeIngredientFromMeal(meal *meals.Meal, ingredient *ingredients.Ingredient) {
	meal.RemoveIngredient(ingredient.Id)

	err := s.mealRepository.Save(meal)
	assert.NoError(s.T(), err)
}

func (s *ShoppingListSuite) addIngredientToBasket(b *basket.Basket, ingredient *ingredients.Ingredient) {
	b.AddItem(basket.NewBasketItem(ingredient.Id))

	err := s.basketRepository.Save(b)

	assert.NoError(s.T(), err)
}

func (s *ShoppingListSuite) removeIngredientFromBasket(b *basket.Basket, ingredient *ingredients.Ingredient) {
	b.RemoveItem(ingredient.Id)

	err := s.basketRepository.Save(b)

	assert.NoError(s.T(), err)
}

func TestShoppingListSuite(t *testing.T) {
	suite.Run(t, new(ShoppingListSuite))
}
