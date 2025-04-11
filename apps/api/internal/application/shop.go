package application

import (
	"fmt"
	"github.com/joe-reed/meal-planner/apps/api/internal/domain/shops"
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
