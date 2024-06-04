import Link from "next/link";
import {
  useAddMealToCurrentShop,
  useCurrentShop,
  useIngredients,
  useMeals,
  useRemoveMealFromCurrentShop,
  useStartShop,
} from "../queries";
import { Meal, Ingredient, Shop } from "../types";
import React, { PropsWithChildren, useState } from "react";
import { Dialog } from "@headlessui/react";

export default function Index() {
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

  return (
    <>
      <nav className="mb-2 flex justify-end space-x-2">
        <Link href="/meals/create" className="button">
          🍲 Create meal
        </Link>
        <StartShopButton />
      </nav>
      <section className="mb-4">
        <Meals meals={meals} currentShop={currentShop} />
      </section>
      <section className="flex justify-between">
        <div className="mr-4 w-1/2 md:mr-0">
          <CurrentShop meals={meals} currentShop={currentShop} />
        </div>
        <div className="w-1/2 md:w-1/4">
          <ShoppingList
            meals={meals}
            currentShop={currentShop}
            ingredients={ingredients}
          />
        </div>
      </section>
    </>
  );
}

function Meals({
  meals,
  currentShop,
}: {
  meals: Meal[];
  currentShop: Shop | null;
}) {
  return (
    <>
      <h2 className="mb-2 font-bold">Meals</h2>
      <ul className="flex flex-wrap">
        {meals
          ?.sort((a, b) => {
            return a.name.toLowerCase().localeCompare(b.name.toLowerCase());
          })
          .map((meal) => (
            <li
              key={meal.id}
              className="mb-2 mr-2 flex items-center rounded-lg border px-3 py-1"
            >
              <MealLink meal={meal} />
              <span className="ml-2">
                {currentShop?.meals.some((m) => m.id == meal.id) ? (
                  <RemoveMealFromShopButton mealId={meal.id}>
                    ✅
                  </RemoveMealFromShopButton>
                ) : (
                  <AddMealToShopButton mealId={meal.id} />
                )}
              </span>
            </li>
          ))}
      </ul>
    </>
  );
}

function CurrentShop({
  meals,
  currentShop,
}: {
  meals: Meal[];
  currentShop: Shop | null;
}) {
  return (
    <>
      {currentShop ? (
        <>
          <div className="mb-2 flex justify-between">
            <h2 className="font-bold">Shop #{currentShop.id}</h2>
            <h3 className="text-xs font-bold">{meals.length} meals</h3>
          </div>
          <ul className="flex flex-col space-y-1">
            {currentShop.meals.map((meal) => (
              <li key={meal.id} className="flex w-full justify-between">
                <MealLink meal={meals.find((m) => m.id == meal.id) as Meal} />
                <RemoveMealFromShopButton mealId={meal.id}>
                  <span className="text-xs">❌</span>
                </RemoveMealFromShopButton>
              </li>
            ))}
          </ul>
        </>
      ) : (
        <p>Click Start Shop to get started!</p>
      )}
    </>
  );
}

function StartShopButton() {
  const { mutate } = useStartShop();
  const [isOpen, setIsOpen] = useState(false);

  return (
    <>
      <button className="button" onClick={() => setIsOpen(true)}>
        🛒 Start Shop
      </button>
      <Dialog
        open={isOpen}
        onClose={() => setIsOpen(false)}
        className="relative z-50"
      >
        <div className="fixed inset-0 bg-black/30" aria-hidden="true" />
        <div className="fixed inset-0 flex w-screen items-center justify-center p-4">
          <Dialog.Panel className="mx-auto max-w-sm rounded bg-white px-4 py-3">
            <Dialog.Title className="mb-2 font-bold">
              Start new shop
            </Dialog.Title>
            <p className="mb-2">Are you sure you want to start a new shop?</p>
            <p className="mb-5">
              The previous shop will be finished and a new empty shop will be
              started.
            </p>
            <div className="flex justify-between px-20">
              <form
                onSubmit={(e) => {
                  e.preventDefault();

                  mutate();

                  setIsOpen(false);
                }}
              >
                <button type="submit" className="button">
                  Start
                </button>
              </form>

              <button onClick={() => setIsOpen(false)} className="underline">
                Cancel
              </button>
            </div>
          </Dialog.Panel>
        </div>
      </Dialog>
    </>
  );
}

function AddMealToShopButton({ mealId }: { mealId: string }) {
  const { mutate } = useAddMealToCurrentShop(mealId);

  return (
    <form
      onSubmit={(e) => {
        e.preventDefault();

        mutate();
      }}
    >
      <button type="submit">➕</button>
    </form>
  );
}

function RemoveMealFromShopButton({
  mealId,
  children,
}: PropsWithChildren<{ mealId: string }>) {
  const { mutate } = useRemoveMealFromCurrentShop(mealId);

  return (
    <form
      onSubmit={(e) => {
        e.preventDefault();

        mutate();
      }}
    >
      <button type="submit">{children}</button>
    </form>
  );
}

function MealLink({ meal }: { meal: Meal }) {
  return (
    <Link href={`/meals/${meal.id}`} className="hover:underline">
      {meal.name}
    </Link>
  );
}

function ShoppingList({
  currentShop,
  meals,
  ingredients,
}: {
  currentShop: Shop | null;
  meals: Meal[];
  ingredients: Ingredient[];
}) {
  const shopIngredients = Object.values(
    (currentShop?.meals ?? [])
      .flatMap((shopMeal) => {
        const meal = meals.find((m) => m.id === shopMeal.id) as Meal;

        return meal.ingredients.map((ingredient) => {
          return ingredients.find((i) => i.id === ingredient.id) as Ingredient;
        });
      })
      .reduce(
        (acc, ingredient) => {
          if (!acc[ingredient.id]) {
            acc[ingredient.id] = {
              ...ingredient,
              amount: 0,
            };
          }

          acc[ingredient.id].amount += 1;

          return acc;
        },
        {} as Record<string, Ingredient & { amount: number }>,
      ),
  );

  return (
    <div className="flex w-full flex-col">
      <h2 className="mb-2 font-bold">Shopping List</h2>
      <ul className="w-full">
        {shopIngredients.map((ingredient) => (
          <li key={ingredient.id} className="flex justify-between">
            <span>{ingredient.name}</span>
            <span>{ingredient.amount}</span>
          </li>
        ))}
      </ul>
    </div>
  );
}
