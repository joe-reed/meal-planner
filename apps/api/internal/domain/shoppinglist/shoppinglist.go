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

type ShopIngredient struct {
	product.Product
	MealCount  int                 `json:"mealCount"`
	IsInBasket bool                `json:"isInBasket"`
	Quantities []quantity.Quantity `json:"quantities"`
}

type ShoppingListProjectionOutput struct {
	ShopId       *int                       `json:"shopId"`
	ShoppingList *map[string]ShopIngredient `json:"shoppingList"`
}

func CreateShoppingListProjection(es *sqlStore.SQL) (*eventsourcing.Projection, ShoppingListProjectionOutput) {
	shoppingList := map[string]ShopIngredient{}
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
			shoppingList = map[string]ShopIngredient{}
			s = map[string]string{}
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
		case *shop.MealAdded:
			s[event.Meal.MealId] = event.Meal.MealId
			for _, i := range ms[event.Meal.MealId].Ingredients {
				shopIngredient, ok := shoppingList[i.ProductId]
				if ok {
					shopIngredient.MealCount++
					shopIngredient.Quantities = append(shopIngredient.Quantities, i.Quantity)
					shoppingList[i.ProductId] = shopIngredient
				} else {
					shoppingList[i.ProductId] = ShopIngredient{Product: prods[i.ProductId], MealCount: 1, Quantities: []quantity.Quantity{i.Quantity}}
				}
			}
		case *shop.MealRemoved:
			delete(s, event.Id)
			for _, i := range ms[event.Id].Ingredients {
				shopIngredient, ok := shoppingList[i.ProductId]
				if ok {
					shopIngredient.MealCount--
					shopIngredient.Quantities = removeQuantity(shopIngredient.Quantities, i.Quantity)
					shoppingList[i.ProductId] = shopIngredient
				}

				if shopIngredient.MealCount == 0 {
					delete(shoppingList, i.ProductId)
				}
			}
		case *meal.IngredientAdded:
			m := ms[ev.AggregateID()]
			m.Transition(ev)
			if _, ok := s[ev.AggregateID()]; !ok {
				break
			}
			shopIngredient, ok := shoppingList[event.Ingredient.ProductId]
			if ok {
				shopIngredient.MealCount++
				shopIngredient.Quantities = append(shopIngredient.Quantities, event.Ingredient.Quantity)
				shoppingList[event.Ingredient.ProductId] = shopIngredient
			} else {
				shoppingList[event.Ingredient.ProductId] = ShopIngredient{Product: prods[event.Ingredient.ProductId], MealCount: 1, Quantities: []quantity.Quantity{event.Ingredient.Quantity}}
			}
		case *meal.IngredientRemoved:
			m := ms[ev.AggregateID()]

			quantity, err := findQuantity(m.Ingredients, event.Id)

			if err != nil {
				return err
			}

			m.Transition(ev)

			if _, ok := s[ev.AggregateID()]; !ok {
				break
			}

			shopIngredient, ok := shoppingList[event.Id]

			if ok {
				shopIngredient.MealCount--
				shopIngredient.Quantities = removeQuantity(shopIngredient.Quantities, *quantity)
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
