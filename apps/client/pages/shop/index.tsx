import { useCurrentShop, useIngredients, useMeals } from "../../queries";
import { Ingredient, Meal, Shop } from "../../types";
import React from "react";
import BackButton from "../../components/BackButton";

export default function ShopPage() {
  const mealsQuery = useMeals();
  const currentShopQuery = useCurrentShop();
  const ingredientsQuery = useIngredients();

  if (
    [mealsQuery, currentShopQuery, ingredientsQuery].some(
      (query) => query.isInitialLoading,
    )
  ) {
    return <p>Loading...</p>;
  }

  const queryWithError = [mealsQuery, currentShopQuery, ingredientsQuery].find(
    (query) => query.isError,
  );

  if (queryWithError && queryWithError.error) {
    return <p>Error: {queryWithError.error.message}</p>;
  }

  const meals = mealsQuery.data as Meal[];
  const currentShop = currentShopQuery.data as Shop | null;
  const ingredients = ingredientsQuery.data as Ingredient[];

  const shopIngredients = Object.values(
    (currentShop?.meals ?? [])
      .flatMap((shopMeal) => {
        const meal = meals.find((m) => m.id === shopMeal.id) as Meal;

        return meal.ingredients.map((ingredient) => {
          return ingredients.find((i) => i.id === ingredient.id) as Ingredient;
        });
      })
      .reduce<{ [ingredientId: string]: Ingredient & { mealCount: number } }>(
        (acc, ingredient) => {
          if (!acc[ingredient.id]) {
            acc[ingredient.id] = {
              ...ingredient,
              mealCount: 0,
            };
          }

          acc[ingredient.id].mealCount += 1;

          return acc;
        },
        {},
      ),
  ).reduce<{ [category: string]: (Ingredient & { mealCount: number })[] }>(
    (acc, ingredient) => {
      const { category } = ingredient;

      acc[category] = acc[category] || [];

      acc[category].push(ingredient);

      return acc;
    },
    {},
  );

  return (
    <div className="flex w-full flex-col">
      <div className="mb-4 flex items-center">
        <BackButton className="mr-3" destination="/" />
        <h1 className="text-lg font-bold">Current shop</h1>
      </div>

      {Object.entries(shopIngredients).map(
        ([category, categoryIngredients]) => (
          <div className="mb-4" key={category}>
            <h2 className="mb-2 text-xl font-bold">{category}</h2>
            <ul>
              {categoryIngredients.map((ingredient) => (
                <li
                  key={ingredient.id}
                  className="mb-3 flex items-center justify-between leading-4"
                >
                  <span className="w-4/6 break-words">{ingredient.name}</span>
                  <span>
                    {ingredient.mealCount}{" "}
                    <span className="text-xs">meals</span>
                  </span>
                </li>
              ))}
            </ul>
          </div>
        ),
      )}
    </div>
  );
}
