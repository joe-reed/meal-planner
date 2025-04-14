package application

import "github.com/joe-reed/meal-planner/apps/api/internal/domain/category"

type CategoryApplication struct{}

func NewCategoryApplication() *CategoryApplication {
	return &CategoryApplication{}
}

func (a *CategoryApplication) GetCategories() []category.Category {
	return category.Categories()
}
