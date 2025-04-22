package main

import (
	"context"
	"database/sql"
	sqlStore "github.com/hallgren/eventsourcing/eventstore/sql"
	"github.com/joe-reed/meal-planner/apps/api/internal/application"
	"github.com/joe-reed/meal-planner/apps/api/internal/database"
	"github.com/joe-reed/meal-planner/apps/api/internal/domain/basket"
	"github.com/joe-reed/meal-planner/apps/api/internal/domain/ingredient"
	"github.com/joe-reed/meal-planner/apps/api/internal/domain/meal"
	"github.com/joe-reed/meal-planner/apps/api/internal/domain/shop"
	"github.com/joe-reed/meal-planner/apps/api/internal/domain/shoppinglist"
	"github.com/joe-reed/meal-planner/apps/api/internal/handlers"
	"github.com/labstack/echo/v4"
	"strconv"
	"strings"
)

type EventSubscriber func(func(string))

func main() {
	e := echo.New()

	dbFile := "sqlite/meal-planner.db"
	db, err := database.CreateDatabase(dbFile)

	if err != nil {
		e.Logger.Fatal(err)
	}

	es := sqlStore.Open(db)

	publisher, subscribe := setupEvents()

	addMealRoutes(e, db)
	addUploadRoutes(e, db)
	addShopRoutes(e, db, publisher)
	addIngredientRoutes(e, db)
	addCategoryRoutes(e)
	addBasketRoutes(e, db, subscribe)

	if err != nil {
		e.Logger.Fatal(err)
	}

	p, output := shoppinglist.CreateShoppingListProjection(es)

	result := p.RunToEnd(context.TODO())

	if result.Error != nil {
		e.Logger.Fatal(result.Error)
	}

	e.GET("/shopping-list", func(c echo.Context) error {
		result = p.RunToEnd(context.TODO())
		if result.Error != nil {
			e.Logger.Fatal(result.Error)
		}

		return c.JSON(200, output)
	})

	e.Debug = true
	e.Logger.Fatal(e.Start(":1323"))
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

	handler := handlers.BasketHandler{Application: application.NewBasketApplication(r)}

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

	e.GET("/baskets/:shopId", handler.GetBasket)
	e.POST("/baskets/:shopId/items", handler.AddItemToBasket)
	e.DELETE("/baskets/:shopId/items/:ingredientId", handler.RemoveItemFromBasket)
}

func addMealRoutes(e *echo.Echo, db *sql.DB) {
	mealRepo, err := meal.NewSqliteMealRepository(db)

	if err != nil {
		e.Logger.Fatal(e)
	}

	handler := handlers.MealsHandler{
		Application: application.NewMealApplication(mealRepo),
	}

	e.GET("/meals", handler.GetMeals)
	e.GET("/meals/:id", handler.FindMeal)
	e.POST("/meals", handler.AddMeal)
	e.POST("/meals/:mealId/ingredients", handler.AddIngredientToMeal)
	e.DELETE("/meals/:mealId/ingredients/:ingredientId", handler.RemoveIngredientFromMeal)
	e.PATCH("/meals/:mealId", handler.UpdateMeal)
}

func addUploadRoutes(e *echo.Echo, db *sql.DB) {
	mealRepo, err := meal.NewSqliteMealRepository(db)

	if err != nil {
		e.Logger.Fatal(e)
	}

	ingredientRepo, err := ingredient.NewSqliteIngredientRepository(db)

	if err != nil {
		e.Logger.Fatal(e)
	}

	handler := handlers.UploadHandler{
		Application: application.NewUploadMealsApplication(ingredientRepo, mealRepo),
	}

	e.POST("/meals/upload", handler.UploadMeals)
}

func addShopRoutes(e *echo.Echo, db *sql.DB, publisher func(string)) {
	r, err := shop.NewSqliteShopRepository(db)

	if err != nil {
		e.Logger.Fatal(e)
	}

	handler := handlers.ShopsHandler{Application: application.NewShopApplication(r, publisher)}

	e.GET("/shops/current", handler.CurrentShop)
	e.POST("/shops/current/meals", handler.AddMealToCurrentShop)
	e.DELETE("/shops/current/meals/:mealId", handler.RemoveMealFromCurrentShop)
	e.POST("/shops", handler.StartShop)
}

func addIngredientRoutes(e *echo.Echo, db *sql.DB) {
	r, err := ingredient.NewSqliteIngredientRepository(db)

	if err != nil {
		e.Logger.Fatal(e)
	}

	handler := handlers.IngredientsHandler{Application: application.NewIngredientApplication(r)}

	e.GET("/ingredients", handler.GetIngredients)
	e.POST("/ingredients", handler.AddIngredient)
}

func addCategoryRoutes(e *echo.Echo) {
	handler := handlers.CategoriesHandler{}

	e.GET("/categories", handler.GetCategories)
}
