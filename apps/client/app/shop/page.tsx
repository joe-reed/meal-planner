"use client";

import { useCurrentShop, useIngredients, useMeals } from "../../queries";
import { Ingredient, Meal, Shop } from "../../types";
import React from "react";
import BackButton from "../../components/BackButton";
import { useBasket } from "../../queries/useBasket";
import clsx from "clsx";
import { useAddItemToBasket } from "../../queries/useAddItemToBasket";
import { useRemoveItemFromBasket } from "../../queries/useRemoveItemFromBasket";

export default function ShopPage() {
  const mealsQuery = useMeals();
  const currentShopQuery = useCurrentShop();
  const ingredientsQuery = useIngredients();
  const shopId = currentShopQuery.data?.id;

  const basketQuery = useBasket(shopId, !!shopId);

  const [showItemsInBasket, setShowItemsInBasket] = React.useState(false);

  function toggleShowItemsInBasket() {
    setShowItemsInBasket(!showItemsInBasket);
  }

  if (
    [mealsQuery, currentShopQuery, ingredientsQuery, basketQuery].some(
      (query) => query.isInitialLoading,
    )
  ) {
    return <p>Loading...</p>;
  }

  const queryWithError = [
    mealsQuery,
    currentShopQuery,
    ingredientsQuery,
    basketQuery,
  ].find((query) => query.isError);

  if (queryWithError && queryWithError.error) {
    return <p>Error: {queryWithError.error.message}</p>;
  }

  const meals = mealsQuery.data as Meal[];
  const currentShop = currentShopQuery.data as Shop | null;
  const ingredients = ingredientsQuery.data as Ingredient[];
  const basket = basketQuery.data;

  const shopIngredients = Object.values(
    (currentShop?.meals ?? [])
      .flatMap((shopMeal) => {
        const meal = meals.find((m) => m.id === shopMeal.id) as Meal;

        return meal.ingredients.map((ingredient) => {
          return ingredients.find((i) => i.id === ingredient.id) as Ingredient;
        });
      })
      .map((ingredient) => ({
        ...ingredient,
        isInBasket:
          basket?.items.some((item) => item.ingredientId === ingredient.id) ??
          false,
      }))
      .reduce<{
        [ingredientId: string]: Ingredient & {
          mealCount: number;
          isInBasket: boolean;
        };
      }>((acc, ingredient) => {
        if (!acc[ingredient.id]) {
          acc[ingredient.id] = {
            ...ingredient,
            mealCount: 0,
          };
        }

        acc[ingredient.id].mealCount += 1;

        return acc;
      }, {}),
  );

  const filteredIngredients = shopIngredients.filter(
    (ingredient) => showItemsInBasket || !ingredient.isInBasket,
  );

  const categorisedIngredients = filteredIngredients.reduce<{
    [category: string]: (Ingredient & {
      mealCount: number;
      isInBasket: boolean;
    })[];
  }>((acc, ingredient) => {
    const { category } = ingredient;

    acc[category] = acc[category] || [];

    acc[category].push(ingredient);

    return acc;
  }, {});

  return (
    <div className="flex w-full flex-col">
      <div className="mb-4 flex items-center justify-between">
        <div className="flex items-center">
          <BackButton className="mr-3" destination="/" />
          <h1 className="text-lg font-bold">Current shop</h1>
        </div>
        <button onClick={toggleShowItemsInBasket} className="button">
          Toggle all
        </button>
      </div>

      {filteredIngredients.length === 0 ? (
        shopIngredients.length === 0 ? (
          <p className="text-center">
            No ingredients in this shop yet. Go back and add some meals!
          </p>
        ) : (
          <p className="text-center">
            All ingredients are in basket. Use the button above to show all
            ingredients.
          </p>
        )
      ) : null}

      {Object.keys(categorisedIngredients).sort().map(
        (category) => (
          <div className="mb-4" key={category}>
            <h2 className="mb-2 text-xl font-bold">{category}</h2>
            <ul>
              {categorisedIngredients[category].map((ingredient) => (
                <IngredientListItem
                  key={ingredient.id}
                  ingredient={ingredient}
                  shopId={shopId as string}
                />
              ))}
            </ul>
          </div>
        ),
      )}
    </div>
  );
}

function IngredientListItem({
  ingredient,
  shopId,
}: {
  ingredient: Ingredient & {
    mealCount: number;
    isInBasket: boolean;
  };
  shopId: string;
}) {
  const { mutate: addItemToBasket } = useAddItemToBasket(shopId);
  const { mutate: removeItemFromBasket } = useRemoveItemFromBasket(shopId);

  return (
    <li className="mb-3 flex items-center justify-between leading-4">
      <label
        className={clsx("flex w-full justify-between break-words pr-6", {
          "line-through": ingredient.isInBasket,
        })}
      >
        {ingredient.name}

        <input
          type="checkbox"
          checked={ingredient.isInBasket}
          onChange={() => {
            if (ingredient.isInBasket) {
              removeItemFromBasket(ingredient.id);
            } else {
              addItemToBasket({ ingredientId: ingredient.id });
            }
          }}
        />
      </label>
      <span className="flex-no-shrink whitespace-nowrap">
        {ingredient.mealCount} <span className="text-xs">meals</span>
      </span>
    </li>
  );
}
