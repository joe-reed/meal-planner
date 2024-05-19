package meals

type Created struct {
	Id   string
	Name string
}

type IngredientAdded struct {
	Ingredient MealIngredient
}

type IngredientRemoved struct {
	Id string
}
