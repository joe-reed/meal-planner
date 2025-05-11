package application

import (
	"github.com/joe-reed/meal-planner/apps/api/internal/domain/category"
	"github.com/joe-reed/meal-planner/apps/api/internal/domain/product"
	"log/slog"
)

type ProductApplication struct {
	r product.ProductRepository
}

func NewProductApplication(r product.ProductRepository) *ProductApplication {
	return &ProductApplication{r: r}
}

type ValidationError struct {
	Field   string
	Message string
}

func (e *ValidationError) Error() string {
	return e.Message
}

func (a *ProductApplication) AddProduct(id string, name product.ProductName, category category.CategoryName) (*product.Product, error) {
	err := validateId(id)
	if err != nil {
		return nil, err
	}

	err = validateName(name.String())
	if err != nil {
		return nil, err
	}

	i, err := product.NewProduct(id, name, category)
	if err != nil {
		return nil, err
	}

	slog.Debug("Adding product", "product", i)

	if err := a.r.Add(i); err != nil {
		return i, err
	}

	return i, nil
}

// todo: reduce duplication, standardise validation or use library
func validateId(id string) error {
	return validateNotEmpty("id", id)
}

func validateName(name string) error {
	return validateNotEmpty("name", name)
}

func validateNotEmpty(field string, value string) error {
	if value == "" {
		return &ValidationError{
			Field:   field,
			Message: field + "cannot be empty",
		}
	}
	return nil
}

func (a *ProductApplication) GetProducts() ([]*product.Product, error) {
	ps, err := a.r.Get()

	if err != nil {
		return nil, err
	}

	return ps, nil
}
