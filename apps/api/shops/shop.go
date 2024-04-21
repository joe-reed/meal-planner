package shops

type Shop struct {
	Id    int         `json:"id"`
	Meals []*ShopMeal `json:"meals"`
}

type ShopMeal struct {
	MealId string `json:"id"`
}

func NewShop(id int) *Shop {
	return &Shop{
		Id:    id,
		Meals: []*ShopMeal{},
	}
}

func (s *Shop) AddMeal(m *ShopMeal) *Shop {
	s.Meals = append(s.Meals, m)
	return s
}

func (s *Shop) SetMeals(m []*ShopMeal) *Shop {
	s.Meals = m
	return s
}
