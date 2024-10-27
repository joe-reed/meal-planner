package main

import (
	"database/sql"
	"github.com/joe-reed/meal-planner/apps/api/basket"
	"github.com/joe-reed/meal-planner/apps/api/categories"
	"github.com/joe-reed/meal-planner/apps/api/database"
	"github.com/joe-reed/meal-planner/apps/api/ingredients"
	"github.com/joe-reed/meal-planner/apps/api/meals"
	"github.com/joe-reed/meal-planner/apps/api/shops"
	"github.com/labstack/echo/v4"
	"strconv"
	"strings"
)

type EventSubscriber func(func(string))

func main() {
	e := echo.New()

	dbFile := "meal-planner.db"
	db, err := database.CreateDatabase(dbFile)
	if err != nil {
		e.Logger.Fatal(err)
	}

	publisher, subscribe := setupEvents()

	addMealRoutes(e, db)
	addShopRoutes(e, db, publisher)
	addIngredientRoutes(e, db)
	addCategoryRoutes(e)
	addBasketRoutes(e, db, subscribe)

	e.Debug = true
	e.Logger.Fatal(e.Start("localhost:1323"))
}

func setupEvents() (func(string), func(func(string))) {
	messageChannel := make(chan string)

	return getPublisher(messageChannel), getSubscribe(messageChannel)
}

func getPublisher(messageChannel chan string) func(string) {
	return func(message string) {
		messageChannel <- message
	}
}

func getSubscribe(messageChannel chan string) func(func(string)) {
	return func(f func(string)) {
		go func() {
			for {
				f(<-messageChannel)
			}
		}()
	}
}

func addBasketRoutes(e *echo.Echo, db *sql.DB, subscribe func(func(string))) {
	r, err := basket.NewSqliteBasketRepository(db)

	if err != nil {
		e.Logger.Fatal(e)
	}

	handler := basket.Handler{BasketRepository: r}

	subscribe(func(message string) {
		parts := strings.Split(message, ":")
		if parts[0] == "shopStarted" {
			e.Logger.Debug("Creating basket for shop " + parts[1])
			shopId, _ := strconv.Atoi(parts[1])
			b, err := basket.NewBasket(shopId)
			if err != nil {
				e.Logger.Error(err)
			}
			err = r.Save(b)
			if err != nil {
				e.Logger.Error(err)
			}
		}
	})

	e.POST("/baskets/:shopId/items", handler.AddItemToBasket)
	e.GET("/baskets/:shopId", handler.GetBasketItems)
	e.DELETE("/baskets/:shopId/items/:ingredientId", handler.RemoveItemFromBasket)
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

func addShopRoutes(e *echo.Echo, db *sql.DB, publisher func(string)) {
	r, err := shops.NewSqliteShopRepository(db)

	if err != nil {
		e.Logger.Fatal(e)
	}

	handler := shops.Handler{ShopRepository: r, Publisher: publisher}

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
