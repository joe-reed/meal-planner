package ingredient

import "github.com/joe-reed/meal-planner/apps/api/internal/domain/category"

type Created struct {
	Id       string
	Name     string
	Category category.CategoryName
}
