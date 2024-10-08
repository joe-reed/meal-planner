package basket

import (
	"github.com/hallgren/eventsourcing"
	"strconv"
)

type Basket struct {
	eventsourcing.AggregateRoot
	ShopId int           `json:"shopId"`
	Items  []*BasketItem `json:"items"`
}

type BasketItem struct {
	IngredientId string `json:"ingredientId"`
}

func (b *Basket) Transition(event eventsourcing.Event) {
	switch e := event.Data().(type) {
	case *Created:
		b.ShopId = e.ShopId
		b.Items = []*BasketItem{}
	case *ItemAdded:
		b.Items = append(b.Items, &BasketItem{IngredientId: e.Item.IngredientId})
	case *ItemRemoved:
		Items := []*BasketItem{}
		for _, Item := range b.Items {
			if Item.IngredientId != e.IngredientId {
				Items = append(Items, Item)
			}
		}
		b.Items = Items
	case *ItemsSet:
		var Items []*BasketItem
		for _, Item := range e.Items {
			Items = append(Items, &BasketItem{IngredientId: Item.IngredientId})
		}
		b.Items = Items
	}
}

func (b *Basket) Register(r eventsourcing.RegisterFunc) {
	r(&Created{}, &ItemAdded{}, &ItemRemoved{}, &ItemsSet{})
}

func NewBasket(shopId int) (*Basket, error) {
	b := &Basket{}

	err := b.SetID(strconv.Itoa(shopId))

	if err != nil {
		return nil, err
	}

	b.TrackChange(b, &Created{ShopId: shopId})

	return b, nil
}

func (b *Basket) AddItem(m *BasketItem) *Basket {
	b.TrackChange(b, &ItemAdded{Item: *m})
	return b
}

func (b *Basket) SetItems(m []*BasketItem) *Basket {
	b.TrackChange(b, &ItemsSet{Items: m})
	return b
}

func (b *Basket) RemoveItem(id string) {
	b.TrackChange(b, &ItemRemoved{IngredientId: id})
}
