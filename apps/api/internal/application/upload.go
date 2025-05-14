package application

import (
	"bytes"
	"encoding/csv"
	"errors"
	"fmt"
	"github.com/joe-reed/meal-planner/apps/api/internal/domain/meal"
	"github.com/joe-reed/meal-planner/apps/api/internal/domain/product"
	"github.com/joe-reed/meal-planner/apps/api/internal/domain/quantity"
	"io"
	"log/slog"
	"strconv"
)

type UploadMealsApplication struct {
	ProductRepository *product.EventSourcedProductRepository
	MealRepository    *meal.EventSourcedMealRepository
}

func NewUploadMealsApplication(productRepository *product.EventSourcedProductRepository, mealRepository *meal.EventSourcedMealRepository) *UploadMealsApplication {
	return &UploadMealsApplication{
		ProductRepository: productRepository,
		MealRepository:    mealRepository,
	}
}

type ProductsNotFound struct {
	NotFoundProducts []product.ProductName
}

func (*ProductsNotFound) Error() string {
	return "products not found"
}

func (a *UploadMealsApplication) UploadMeals(src io.Reader) error {
	ms, notFoundProducts, err := a.parseMeals(src)

	if err != nil {
		return err
	}

	if len(notFoundProducts) > 0 {
		return &ProductsNotFound{
			NotFoundProducts: notFoundProducts,
		}
	}

	for _, m := range ms {
		m, _ := a.MealRepository.FindByName(m.Name)

		if m != nil {
			return &MealAlreadyExists{
				MealName: m.Name,
			}
		}
	}

	slog.Info("uploading meals", "meals", ms)

	for _, m := range ms {
		if err := a.MealRepository.Save(m); err != nil {
			return err
		}
	}

	return nil
}

func (a *UploadMealsApplication) parseMeals(src io.Reader) (meals []*meal.Meal, notFoundProducts []product.ProductName, err error) {
	var buf bytes.Buffer
	_, err = buf.ReadFrom(src)

	if err != nil {
		return nil, nil, err
	}

	var m *meal.Meal

	csvReader := csv.NewReader(&buf)
	records, err := csvReader.ReadAll()

	if err != nil {
		return nil, nil, err
	}

	for i, record := range records {
		if i == 0 {
			if record[0] != "name" || record[1] != "product" || record[2] != "amount" || record[3] != "unit" {
				return nil, nil, errors.New("invalid csv header")
			}

			continue
		}

		if len(record) != 4 {
			return nil, nil, errors.New("invalid csv row")
		}

		mealName := record[0]

		if m == nil || mealName != m.Name {
			if m != nil {
				meals = append(meals, m)
			}

			m = meal.NewMealBuilder().WithName(mealName).Build()
		}

		productName, err := product.NewProductName(record[1])

		if err != nil {
			return nil, nil, err
		}

		amount, err := strconv.Atoi(record[2])

		if err != nil {
			return nil, nil, err
		}

		unit, ok := quantity.UnitFromString(record[3])

		if !ok {
			return nil, nil, fmt.Errorf("invalid unit: %s", record[3])
		}

		i, err := a.ProductRepository.GetByName(productName)

		if err != nil {
			notFoundProducts = append(notFoundProducts, productName)
			continue
		}

		m.AddIngredient(*meal.NewIngredient(i.Id).WithQuantity(amount, unit))
	}

	if m != nil {
		meals = append(meals, m)
	}

	return meals, notFoundProducts, nil
}
