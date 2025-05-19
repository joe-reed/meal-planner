package shoppinglist

import (
	"errors"
	"github.com/hallgren/eventsourcing"
	"github.com/hallgren/eventsourcing/core"
	sqlStore "github.com/hallgren/eventsourcing/eventstore/sql"
	"github.com/joe-reed/meal-planner/apps/api/internal/domain/basket"
	"github.com/joe-reed/meal-planner/apps/api/internal/domain/meal"
	"github.com/joe-reed/meal-planner/apps/api/internal/domain/product"
	"github.com/joe-reed/meal-planner/apps/api/internal/domain/quantity"
	"github.com/joe-reed/meal-planner/apps/api/internal/domain/shop"
)

type ShopItem struct {
	product.Product
	MealCount  int                 `json:"mealCount"`
	IsInBasket bool                `json:"isInBasket"`
	Quantities []quantity.Quantity `json:"quantities"`
}

type ShoppingListProjectionOutput struct {
	ShopId       *int                 `json:"shopId"`
	ShoppingList *map[string]ShopItem `json:"shoppingList"`
}

func CreateShoppingListProjection(es *sqlStore.SQL) (*eventsourcing.Projection, ShoppingListProjectionOutput) {
	shoppingList := map[string]ShopItem{}
	s := map[string]string{}
	prods := map[string]product.Product{}
	ms := map[string]*meal.Meal{}
	shopId := new(int)

	start := core.Version(0)

	p := eventsourcing.NewProjection(func() (core.Iterator, error) {
		return es.All(start, 10)
	}, func(ev eventsourcing.Event) error {
		switch event := ev.Data().(type) {
		case *product.Created:
			prod := product.Product{}
			prod.Transition(ev)
			prods[event.Id] = prod
		case *meal.Created:
			m := meal.Meal{}
			m.Transition(ev)
			ms[event.Id] = &m
		case *shop.Created:
			shoppingList = map[string]ShopItem{}
			s = map[string]string{}
			*shopId = event.Id
		case *basket.ItemAdded:
			shopItem, ok := shoppingList[event.Item.IngredientId]
			if ok {
				shopItem.IsInBasket = true
				shoppingList[event.Item.IngredientId] = shopItem
			}
		case *basket.ItemRemoved:
			shopItem, ok := shoppingList[event.IngredientId]
			if ok {
				shopItem.IsInBasket = false
				shoppingList[event.IngredientId] = shopItem
			}
		case *shop.MealAdded:
			s[event.Meal.MealId] = event.Meal.MealId
			for _, i := range ms[event.Meal.MealId].Ingredients {
				shopItem, ok := shoppingList[i.ProductId]
				if ok {
					shopItem.MealCount++
					shopItem.Quantities = append(shopItem.Quantities, i.Quantity)
					shoppingList[i.ProductId] = shopItem
				} else {
					shoppingList[i.ProductId] = ShopItem{Product: prods[i.ProductId], MealCount: 1, Quantities: []quantity.Quantity{i.Quantity}}
				}
			}
		case *shop.MealRemoved:
			delete(s, event.Id)
			for _, i := range ms[event.Id].Ingredients {
				shopItem, ok := shoppingList[i.ProductId]
				if ok {
					shopItem.MealCount--
					shopItem.Quantities = removeQuantity(shopItem.Quantities, i.Quantity)
					shoppingList[i.ProductId] = shopItem
				}

				if shopItem.MealCount == 0 {
					delete(shoppingList, i.ProductId)
				}
			}
		case *meal.IngredientAdded:
			m := ms[ev.AggregateID()]
			m.Transition(ev)
			if _, ok := s[ev.AggregateID()]; !ok {
				break
			}
			shopItem, ok := shoppingList[event.Ingredient.ProductId]
			if ok {
				shopItem.MealCount++
				shopItem.Quantities = append(shopItem.Quantities, event.Ingredient.Quantity)
				shoppingList[event.Ingredient.ProductId] = shopItem
			} else {
				shoppingList[event.Ingredient.ProductId] = ShopItem{Product: prods[event.Ingredient.ProductId], MealCount: 1, Quantities: []quantity.Quantity{event.Ingredient.Quantity}}
			}
		case *meal.IngredientRemoved:
			m := ms[ev.AggregateID()]

			q, err := findQuantity(m.Ingredients, event.Id)

			if err != nil {
				return err
			}

			m.Transition(ev)

			if _, ok := s[ev.AggregateID()]; !ok {
				break
			}

			shopItem, ok := shoppingList[event.Id]

			if ok {
				shopItem.MealCount--
				shopItem.Quantities = removeQuantity(shopItem.Quantities, *q)
				shoppingList[event.Id] = shopItem
			}

			if shopItem.MealCount == 0 {
				delete(shoppingList, event.Id)
			}
		case *shop.ItemAdded:
			shopItem, ok := shoppingList[event.Item.ProductId]
			if ok {
				shopItem.MealCount++
				shopItem.Quantities = append(shopItem.Quantities, event.Item.Quantity)
				shoppingList[event.Item.ProductId] = shopItem
			} else {
				shoppingList[event.Item.ProductId] = ShopItem{Product: prods[event.Item.ProductId], MealCount: 1, Quantities: []quantity.Quantity{event.Item.Quantity}}
			}
		}

		start = core.Version(ev.GlobalVersion() + 1)

		return nil
	})

	return p, ShoppingListProjectionOutput{shopId, &shoppingList}
}

func findQuantity(ingredients []meal.Ingredient, ingredientId string) (*quantity.Quantity, error) {
	for _, i := range ingredients {
		if i.ProductId == ingredientId {
			return &i.Quantity, nil
		}
	}

	return nil, errors.New("ingredient not found")
}

func removeQuantity(quantities []quantity.Quantity, quantity quantity.Quantity) []quantity.Quantity {
	for i, q := range quantities {
		if q == quantity {
			return append(quantities[:i], quantities[i+1:]...)
		}
	}

	return quantities
}
