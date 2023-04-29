package main

import (
	"github.com/joe-reed/meal-planner/apps/api/meals"
	"github.com/labstack/echo/v4"
)

func main() {
	e := echo.New()

	mealRepository, err := meals.NewSqliteMealRepository("meals.db")

	if err != nil {
		e.Logger.Fatal(e)
	}

	handler := meals.Handler{MealRepository: mealRepository}
	e.GET("/", handler.GetMeals)
	e.POST("/", handler.AddMeal)

	e.Debug = true

	e.Logger.Fatal(e.Start(":1323"))
}
