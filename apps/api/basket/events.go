package basket

type Created struct {
	ShopId int
}

type ItemAdded struct {
	Item BasketItem
}

type ItemRemoved struct {
	IngredientId string
}

type ItemsSet struct {
	Items []*BasketItem
}
