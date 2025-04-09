package application

import (
	"github.com/joe-reed/meal-planner/apps/api/internal/domain/basket"
	"log/slog"
)

type BasketApplication struct {
	r *basket.BasketRepository
}

func NewBasketApplication(r *basket.BasketRepository) *BasketApplication {
	return &BasketApplication{r: r}
}

func (a *BasketApplication) AddItemToBasket(shopId int, basketItem *basket.BasketItem) (*basket.Basket, error) {
	b, err := a.r.FindByShopId(shopId)

	if err != nil {
		return nil, err
	}

	slog.Debug("Adding item to basket", "shopId", shopId, "basketItem", basketItem)

	b.AddItem(basketItem)

	if err := a.r.Save(b); err != nil {
		return nil, err
	}

	return b, nil
}

func (a *BasketApplication) RemoveItemFromBasket(shopId int, ingredientId string) (*basket.Basket, error) {
	slog.Debug("Removing item from basket", "shopId", shopId, "ingredientId", ingredientId)

	b, err := a.r.FindByShopId(shopId)

	if err != nil {
		return nil, err
	}

	b.RemoveItem(ingredientId)

	if err := a.r.Save(b); err != nil {
		return nil, err
	}

	return b, nil
}

func (a *BasketApplication) GetBasket(shopId int) (*basket.Basket, error) {
	b, err := a.r.FindByShopId(shopId)

	if err != nil {
		return nil, err
	}

	return b, nil
}
