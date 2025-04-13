package application

import (
	"fmt"
	"github.com/joe-reed/meal-planner/apps/api/internal/domain/shops"
	"log/slog"
)

type ShopApplication struct {
	r         *shops.ShopRepository
	Publisher func(string)
}

func NewShopApplication(r *shops.ShopRepository, p func(string)) *ShopApplication {
	return &ShopApplication{r: r, Publisher: p}
}

func (a *ShopApplication) GetCurrentShop() (*shops.Shop, error) {
	return a.r.Current()
}

func (a *ShopApplication) StartShop() (*shops.Shop, error) {
	shop, err := a.r.Current()
	if err != nil {
		return nil, err
	}

	var newShop *shops.Shop

	if shop == nil {
		newShop, err = shops.NewShop(1)
	} else {
		newShop, err = shops.NewShop(shop.Id + 1)
	}

	if err != nil {
		return nil, err
	}

	err = a.r.Save(newShop)

	if err != nil {
		return nil, err
	}

	a.Publisher(fmt.Sprintf("shopStarted:%d", newShop.Id))

	return newShop, nil
}

func (a *ShopApplication) AddMealToCurrentShop(shopMeal *shops.ShopMeal) (*shops.Shop, error) {
	shop, err := a.r.Current()

	if err != nil {
		return nil, err
	}

	if shop == nil {
		return nil, fmt.Errorf("no current shop")
	}

	slog.Debug("Adding meal to shop", "shopId", shop.Id, "mealId", shopMeal.MealId)
	shop.AddMeal(shopMeal)

	if err := a.r.Save(shop); err != nil {
		return nil, err
	}

	return shop, nil
}

func (a *ShopApplication) RemoveMealFromCurrentShop(mealId string) (*shops.Shop, error) {
	shop, err := a.r.Current()

	if err != nil {
		return nil, err
	}

	if shop == nil {
		return nil, fmt.Errorf("no current shop")
	}

	slog.Debug("Removing meal from shop", "shopId", shop.Id, "mealId", mealId)
	shop.RemoveMeal(mealId)

	if err := a.r.Save(shop); err != nil {
		return nil, err
	}

	return shop, nil
}
