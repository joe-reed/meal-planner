package meal

type Created struct {
	Id              string
	Name            string
	Url             string
	MealIngredients []MealIngredient
}

type IngredientAdded struct {
	Ingredient MealIngredient
}

type IngredientRemoved struct {
	Id string
}

type UrlUpdated struct {
	Url string
}
