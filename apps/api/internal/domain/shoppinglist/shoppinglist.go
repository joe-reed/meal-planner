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

type ShoppingListItem struct {
	product.Product
	MealCount  int                 `json:"mealCount"`
	IsInBasket bool                `json:"isInBasket"`
	Quantities []quantity.Quantity `json:"quantities"`
}

type ShoppingListProjectionOutput struct {
	ShopId       *int                         `json:"shopId"`
	ShoppingList *map[string]ShoppingListItem `json:"shoppingList"`
}

func CreateShoppingListProjection(es *sqlStore.SQLite) (*eventsourcing.Projection, ShoppingListProjectionOutput) {
	shoppingList := map[string]ShoppingListItem{}
	s := map[string]string{}
	prods := map[string]product.Product{}
	ms := map[string]*meal.Meal{}
	shopId := new(int)
	items := make(map[string]*shop.Item)

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
			shoppingList = map[string]ShoppingListItem{}
			s = map[string]string{}
			*shopId = event.Id
		case *basket.ItemAdded:
			shoppingListItem, ok := shoppingList[event.Item.IngredientId]
			if ok {
				shoppingListItem.IsInBasket = true
				shoppingList[event.Item.IngredientId] = shoppingListItem
			}
		case *basket.ItemRemoved:
			shoppingListItem, ok := shoppingList[event.IngredientId]
			if ok {
				shoppingListItem.IsInBasket = false
				shoppingList[event.IngredientId] = shoppingListItem
			}
		case *shop.MealAdded:
			s[event.Meal.MealId] = event.Meal.MealId
			for _, i := range ms[event.Meal.MealId].Ingredients {
				shoppingListItem, ok := shoppingList[i.ProductId]
				if ok {
					shoppingListItem.MealCount++
					shoppingListItem.Quantities = append(shoppingListItem.Quantities, i.Quantity)
					shoppingList[i.ProductId] = shoppingListItem
				} else {
					shoppingList[i.ProductId] = ShoppingListItem{Product: prods[i.ProductId], MealCount: 1, Quantities: []quantity.Quantity{i.Quantity}}
				}
			}
		case *shop.MealRemoved:
			delete(s, event.Id)
			for _, i := range ms[event.Id].Ingredients {
				shoppingListItem, ok := shoppingList[i.ProductId]
				if ok {
					shoppingListItem.MealCount--
					shoppingListItem.Quantities = removeQuantity(shoppingListItem.Quantities, i.Quantity)
					shoppingList[i.ProductId] = shoppingListItem
				}

				if shoppingListItem.MealCount == 0 {
					delete(shoppingList, i.ProductId)
				}
			}
		case *meal.IngredientAdded:
			m := ms[ev.AggregateID()]
			m.Transition(ev)
			if _, ok := s[ev.AggregateID()]; !ok {
				break
			}
			shoppingListItem, ok := shoppingList[event.Ingredient.ProductId]
			if ok {
				shoppingListItem.MealCount++
				shoppingListItem.Quantities = append(shoppingListItem.Quantities, event.Ingredient.Quantity)
				shoppingList[event.Ingredient.ProductId] = shoppingListItem
			} else {
				shoppingList[event.Ingredient.ProductId] = ShoppingListItem{Product: prods[event.Ingredient.ProductId], MealCount: 1, Quantities: []quantity.Quantity{event.Ingredient.Quantity}}
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

			shoppingListItem, ok := shoppingList[event.Id]

			if ok {
				shoppingListItem.MealCount--
				shoppingListItem.Quantities = removeQuantity(shoppingListItem.Quantities, *q)
				shoppingList[event.Id] = shoppingListItem
			}

			if shoppingListItem.MealCount == 0 {
				delete(shoppingList, event.Id)
			}
		case *shop.ItemAdded:
			shoppingListItem, ok := shoppingList[event.Item.ProductId]
			if ok {
				shoppingListItem.MealCount++
				shoppingListItem.Quantities = append(shoppingListItem.Quantities, event.Item.Quantity)
				shoppingList[event.Item.ProductId] = shoppingListItem
				items[event.Item.ProductId] = event.Item
			} else {
				shoppingList[event.Item.ProductId] = ShoppingListItem{Product: prods[event.Item.ProductId], MealCount: 1, Quantities: []quantity.Quantity{event.Item.Quantity}}
				items[event.Item.ProductId] = event.Item
			}
		case *shop.ItemRemoved:
			shoppingListItem, ok := shoppingList[event.ProductId]
			if ok {
				shoppingListItem.MealCount--
				shoppingListItem.Quantities = removeQuantity(shoppingListItem.Quantities, items[event.ProductId].Quantity)
				shoppingList[event.ProductId] = shoppingListItem
				delete(items, event.ProductId)
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
