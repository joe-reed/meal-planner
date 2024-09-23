package main

import (
	"database/sql"
	"github.com/joe-reed/meal-planner/apps/api/categories"
	"github.com/joe-reed/meal-planner/apps/api/database"
	"github.com/joe-reed/meal-planner/apps/api/ingredients"
	"github.com/joe-reed/meal-planner/apps/api/meals"
	"github.com/joe-reed/meal-planner/apps/api/shops"
	"github.com/labstack/echo/v4"
)

func main() {
	e := echo.New()

	dbFile := "meal-planner.db"
	db, err := database.CreateDatabase(dbFile)
	if err != nil {
		e.Logger.Fatal(err)
	}

	addMealRoutes(e, db)
	addShopRoutes(e, db)
	addIngredientRoutes(e, db)
	addCategoryRoutes(e)

	e.Debug = true

	e.Logger.Fatal(e.Start("localhost:1323"))
}

func addMealRoutes(e *echo.Echo, db *sql.DB) {
	r, err := meals.NewSqliteMealRepository(db)

	if err != nil {
		e.Logger.Fatal(e)
	}

	handler := meals.Handler{MealRepository: r}

	e.GET("/", handler.GetMeals)
	e.GET("/meals/:id", handler.GetMeal)
	e.POST("/", handler.AddMeal)
	e.POST("/meals/:mealId/ingredients", handler.AddIngredientToMeal)
	e.DELETE("/meals/:mealId/ingredients/:ingredientId", handler.RemoveIngredientFromMeal)
}

func addShopRoutes(e *echo.Echo, db *sql.DB) {
	r, err := shops.NewSqliteShopRepository(db)

	if err != nil {
		e.Logger.Fatal(e)
	}

	handler := shops.Handler{ShopRepository: r}

	e.GET("/shops/current", handler.CurrentShop)
	e.POST("/shops/current/meals", handler.AddMealToCurrentShop)
	e.DELETE("/shops/current/meals/:mealId", handler.RemoveMealFromCurrentShop)
	e.POST("/shops", handler.StartShop)
}

func addIngredientRoutes(e *echo.Echo, db *sql.DB) {
	r, err := ingredients.NewSqliteIngredientRepository(db)

	if err != nil {
		e.Logger.Fatal(e)
	}

	handler := ingredients.Handler{IngredientRepository: r}

	e.GET("/ingredients", handler.GetIngredients)
	e.POST("/ingredients", handler.AddIngredient)
}

func addCategoryRoutes(e *echo.Echo) {
	handler := categories.Handler{}

	e.GET("/categories", handler.GetCategories)
}
