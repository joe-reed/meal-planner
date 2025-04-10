package categories

import "sort"

type CategoryName int

//go:generate go run github.com/campoy/jsonenums -type=CategoryName
const (
	Fruit CategoryName = iota
	Meat
	FishAndSeafood
	FoodCupboard
	Drinks
	Chilled
	Frozen
	Bakery
	Vegetables
	TeaAndCoffee
	AlcoholicDrinks
	SaucesOilsAndDressings
	PastaRiceAndNoodles
	SeedsNutsAndDriedFruits
	ChocolateAndSweets
	TinsCansAndPackets
	Desserts
	Dairy
	Eggs
)

type Category struct {
	Name string `json:"name"`
}

func Categories() []Category {
	v := make([]Category, 0, len(_CategoryNameValueToName))

	for _, value := range _CategoryNameValueToName {
		v = append(v, Category{Name: value})
	}

	sort.Slice(v, func(i, j int) bool {
		return v[i].Name < v[j].Name
	})

	return v
}
