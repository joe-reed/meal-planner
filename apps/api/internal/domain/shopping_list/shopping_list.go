package shopping_list

import (
	"github.com/hallgren/eventsourcing"
	"github.com/hallgren/eventsourcing/core"
	sqlStore "github.com/hallgren/eventsourcing/eventstore/sql"
	"github.com/joe-reed/meal-planner/apps/api/internal/domain/basket"
	"github.com/joe-reed/meal-planner/apps/api/internal/domain/ingredients"
	"github.com/joe-reed/meal-planner/apps/api/internal/domain/meals"
	"github.com/joe-reed/meal-planner/apps/api/internal/domain/shops"
)

type ShopIngredient struct {
	ingredients.Ingredient
	MealCount  int  `json:"mealCount"`
	IsInBasket bool `json:"isInBasket"`
}

type ShoppingListProjectionOutput struct {
	ShopId       *int                       `json:"shopId"`
	ShoppingList *map[string]ShopIngredient `json:"shoppingList"`
}

func CreateShoppingListProjection(es *sqlStore.SQL) (*eventsourcing.Projection, ShoppingListProjectionOutput) {
	shoppingList := map[string]ShopIngredient{}
	shop := map[string]string{}
	ings := map[string]ingredients.Ingredient{}
	ms := map[string]*meals.Meal{}
	shopId := new(int)

	start := core.Version(0)

	p := eventsourcing.NewProjection(func() (core.Iterator, error) {
		return es.All(start, 10)
	}, func(ev eventsourcing.Event) error {
		switch event := ev.Data().(type) {
		case *ingredients.Created:
			ing := ingredients.Ingredient{}
			ing.Transition(ev)
			ings[event.Id] = ing
		case *meals.Created:
			meal := meals.Meal{}
			meal.Transition(ev)
			ms[event.Id] = &meal
		case *shops.Created:
			shoppingList = map[string]ShopIngredient{}
			shop = map[string]string{}
			*shopId = event.Id
		case *basket.ItemAdded:
			shopIngredient, ok := shoppingList[event.Item.IngredientId]
			if ok {
				shopIngredient.IsInBasket = true
				shoppingList[event.Item.IngredientId] = shopIngredient
			}
		case *basket.ItemRemoved:
			shopIngredient, ok := shoppingList[event.IngredientId]
			if ok {
				shopIngredient.IsInBasket = false
				shoppingList[event.IngredientId] = shopIngredient
			}
		case *shops.MealAdded:
			shop[event.Meal.MealId] = event.Meal.MealId
			for _, ingredient := range ms[event.Meal.MealId].MealIngredients {
				shopIngredient, ok := shoppingList[ingredient.IngredientId]
				if ok {
					shopIngredient.MealCount++
					shoppingList[ingredient.IngredientId] = shopIngredient
				} else {
					shoppingList[ingredient.IngredientId] = ShopIngredient{Ingredient: ings[ingredient.IngredientId], MealCount: 1}
				}
			}
		case *shops.MealRemoved:
			delete(shop, event.Id)
			for _, ingredient := range ms[event.Id].MealIngredients {
				shopIngredient, ok := shoppingList[ingredient.IngredientId]
				if ok {
					shopIngredient.MealCount--
					shoppingList[ingredient.IngredientId] = shopIngredient
				}

				if shopIngredient.MealCount == 0 {
					delete(shoppingList, ingredient.IngredientId)
				}
			}
		case *meals.IngredientAdded:
			meal := ms[ev.AggregateID()]
			meal.Transition(ev)
			if _, ok := shop[ev.AggregateID()]; !ok {
				break
			}
			shopIngredient, ok := shoppingList[event.Ingredient.IngredientId]
			if ok {
				shopIngredient.MealCount++
				shoppingList[event.Ingredient.IngredientId] = shopIngredient
			} else {
				shoppingList[event.Ingredient.IngredientId] = ShopIngredient{Ingredient: ings[event.Ingredient.IngredientId], MealCount: 1}
			}
		case *meals.IngredientRemoved:
			meal := ms[ev.AggregateID()]
			meal.Transition(ev)
			if _, ok := shop[ev.AggregateID()]; !ok {
				break
			}

			shopIngredient, ok := shoppingList[event.Id]
			if ok {
				shopIngredient.MealCount--
				shoppingList[event.Id] = shopIngredient
			}

			if shopIngredient.MealCount == 0 {
				delete(shoppingList, event.Id)
			}
		}

		start = core.Version(ev.GlobalVersion() + 1)

		return nil
	})

	return p, ShoppingListProjectionOutput{shopId, &shoppingList}
}
