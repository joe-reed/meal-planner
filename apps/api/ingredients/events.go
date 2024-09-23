package ingredients

import "github.com/joe-reed/meal-planner/apps/api/categories"

type Created struct {
	Id       string
	Name     string
	Category categories.Category
}
