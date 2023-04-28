package main

import (
	"github.com/joe-reed/meal-planner/meals"

	"github.com/labstack/echo/v4"
)

func main() {
	e := echo.New()
	mealRepository := meals.NewFakeMealRepository()
	mealRepository.Add(meals.NewMealBuilder().WithName("foo").Build())
	mealRepository.Add(meals.NewMealBuilder().WithName("bar").Build())

	handler := meals.Handler{MealRepository: mealRepository}
	e.GET("/", handler.GetMeals)
	e.Logger.Fatal(e.Start(":1323"))
}
