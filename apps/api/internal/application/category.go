package application

import "github.com/joe-reed/meal-planner/apps/api/internal/domain/categories"

type CategoryApplication struct{}

func NewCategoryApplication() *CategoryApplication {
	return &CategoryApplication{}
}

func (a *CategoryApplication) GetCategories() []categories.Category {
	return categories.Categories()
}
