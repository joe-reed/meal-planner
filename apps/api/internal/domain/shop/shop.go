package shop

import (
	"github.com/hallgren/eventsourcing"
	"github.com/hallgren/eventsourcing/aggregate"
	"strconv"
)

type Shop struct {
	aggregate.Root
	Id    int         `json:"id"`
	Meals []*ShopMeal `json:"meals"`
}

type ShopMeal struct {
	MealId string `json:"id"`
}

func (s *Shop) Transition(event eventsourcing.Event) {
	switch e := event.Data().(type) {
	case *Created:
		s.Id = e.Id
		s.Meals = []*ShopMeal{}
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
	}
}

func (s *Shop) Register(r aggregate.RegisterFunc) {
	r(&Created{}, &MealAdded{}, &MealRemoved{}, &MealsSet{})
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
