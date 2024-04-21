package main

import (
	"github.com/joe-reed/meal-planner/apps/api/meals"
	"github.com/joe-reed/meal-planner/apps/api/shops"
	"github.com/labstack/echo/v4"
)

func main() {
	e := echo.New()

	dbFile := "meal-planner.db"
	addMealRoutes(e, dbFile)
	addShopRoutes(e, dbFile)

	e.Debug = true

	e.Logger.Fatal(e.Start(":1323"))
}

func addMealRoutes(e *echo.Echo, dbFile string) {
	r, err := meals.NewSqliteMealRepository(dbFile)

	if err != nil {
		e.Logger.Fatal(e)
	}

	handler := meals.Handler{MealRepository: r}

	e.GET("/", handler.GetMeals)
	e.GET("/meals/:id", handler.GetMeal)
	e.POST("/", handler.AddMeal)
}

func addShopRoutes(e *echo.Echo, dbFile string) {
	r, err := shops.NewSqliteShopRepository(dbFile)

	if err != nil {
		e.Logger.Fatal(e)
	}

	handler := shops.Handler{ShopRepository: r}

	e.GET("/shops/current", handler.CurrentShop)
	e.POST("/shops", handler.StartShop)
	e.POST("/shops/:shopId/meals", handler.AddMealToShop)
}
