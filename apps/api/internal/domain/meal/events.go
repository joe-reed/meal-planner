package meal

type Created struct {
	Id          string
	Name        string
	Url         string
	Ingredients []Ingredient
}

type IngredientAdded struct {
	Ingredient Ingredient
}

type IngredientRemoved struct {
	Id string
}

type NameUpdated struct {
	Name string
}

type UrlUpdated struct {
	Url string
}
