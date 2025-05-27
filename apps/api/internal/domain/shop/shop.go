package shop

import (
	"github.com/hallgren/eventsourcing"
	"github.com/hallgren/eventsourcing/aggregate"
	"github.com/joe-reed/meal-planner/apps/api/internal/domain/quantity"
	"strconv"
)

type Shop struct {
	aggregate.Root
	Id    int         `json:"id"`
	Meals []*ShopMeal `json:"meals"`
	Items []*Item     `json:"items"`
}

type ShopMeal struct {
	MealId string `json:"id"`
}

type Item struct {
	ProductId string            `json:"productId"`
	Quantity  quantity.Quantity `json:"quantity"`
}

func (s *Shop) Transition(event eventsourcing.Event) {
	switch e := event.Data().(type) {
	case *Created:
		s.Id = e.Id
		s.Meals = []*ShopMeal{}
		s.Items = []*Item{}
	case *MealAdded:
		s.Meals = append(s.Meals, &ShopMeal{MealId: e.Meal.MealId})
	case *MealRemoved:
		meals := []*ShopMeal{}
		for _, meal := range s.Meals {
			if meal.MealId != e.Id {
				meals = append(meals, meal)
			}
		}
		s.Meals = meals
	case *MealsSet:
		var meals []*ShopMeal
		for _, meal := range e.Meals {
			meals = append(meals, &ShopMeal{MealId: meal.MealId})
		}
		s.Meals = meals
	case *ItemAdded:
		s.Items = append(s.Items, e.Item)
	case *ItemRemoved:
		items := []*Item{}
		for _, item := range s.Items {
			if item.ProductId != e.ProductId {
				items = append(items, item)
			}
		}
		s.Items = items
	}
}

func (s *Shop) Register(r aggregate.RegisterFunc) {
	r(&Created{}, &MealAdded{}, &MealRemoved{}, &MealsSet{}, &ItemAdded{}, &ItemRemoved{})
}

func NewShop(id int) (*Shop, error) {
	s := &Shop{}

	err := s.SetID(strconv.Itoa(id))

	if err != nil {
		return nil, err
	}

	aggregate.TrackChange(s, &Created{Id: id})

	return s, nil
}

func (s *Shop) AddMeal(m *ShopMeal) *Shop {
	aggregate.TrackChange(s, &MealAdded{Meal: *m})
	return s
}

func (s *Shop) SetMeals(m []*ShopMeal) *Shop {
	aggregate.TrackChange(s, &MealsSet{Meals: m})
	return s
}

func (s *Shop) RemoveMeal(id string) {
	aggregate.TrackChange(s, &MealRemoved{Id: id})
}

func (s *Shop) AddItem(item *Item) {
	aggregate.TrackChange(s, &ItemAdded{Item: item})
}

func (s *Shop) RemoveItem(productId string) {
	aggregate.TrackChange(s, &ItemRemoved{ProductId: productId})
}
