package shop

type Created struct {
	Id int
}

type MealAdded struct {
	Meal ShopMeal
}

type MealRemoved struct {
	Id string
}

type MealsSet struct {
	Meals []*ShopMeal
}

type ItemAdded struct {
	Item *Item
}
