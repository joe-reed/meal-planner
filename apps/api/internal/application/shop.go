package application

import (
	"fmt"
	"github.com/joe-reed/meal-planner/apps/api/internal/domain/shop"
	"log/slog"
)

type ShopApplication struct {
	r         *shop.ShopRepository
	Publisher func(string)
}

func NewShopApplication(r *shop.ShopRepository, p func(string)) *ShopApplication {
	return &ShopApplication{r: r, Publisher: p}
}

func (a *ShopApplication) GetCurrentShop() (*shop.Shop, error) {
	return a.r.Current()
}

func (a *ShopApplication) StartShop() (*shop.Shop, error) {
	s, err := a.r.Current()
	if err != nil {
		return nil, err
	}

	var newShop *shop.Shop

	if s == nil {
		newShop, err = shop.NewShop(1)
	} else {
		newShop, err = shop.NewShop(s.Id + 1)
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

func (a *ShopApplication) AddMealToCurrentShop(shopMeal *shop.ShopMeal) (*shop.Shop, error) {
	s, err := a.r.Current()

	if err != nil {
		return nil, err
	}

	if s == nil {
		return nil, fmt.Errorf("no current shop")
	}

	slog.Debug("Adding meal to shop", "shopId", s.Id, "mealId", shopMeal.MealId)
	s.AddMeal(shopMeal)

	if err := a.r.Save(s); err != nil {
		return nil, err
	}

	return s, nil
}

func (a *ShopApplication) RemoveMealFromCurrentShop(mealId string) (*shop.Shop, error) {
	s, err := a.r.Current()

	if err != nil {
		return nil, err
	}

	if s == nil {
		return nil, fmt.Errorf("no current shop")
	}

	slog.Debug("Removing meal from shop", "shopId", s.Id, "mealId", mealId)
	s.RemoveMeal(mealId)

	if err := a.r.Save(s); err != nil {
		return nil, err
	}

	return s, nil
}
