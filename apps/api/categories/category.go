package categories

import "sort"

type Category int

//go:generate go run github.com/campoy/jsonenums -type=Category
const (
	Fruit Category = iota
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

func Categories() []string {
	v := make([]string, 0, len(_CategoryValueToName))

	for _, value := range _CategoryValueToName {
		v = append(v, value)
	}

	sort.Strings(v)

	return v
}
