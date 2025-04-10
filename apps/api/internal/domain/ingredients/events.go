package ingredients

import "github.com/joe-reed/meal-planner/apps/api/internal/domain/categories"

type Created struct {
	Id       string
	Name     string
	Category categories.CategoryName
}
