"use client";

import Link from "next/link";
import {
  useAddMealToCurrentShop,
  useCurrentShop,
  useMeals,
  useRemoveMealFromCurrentShop,
  useStartShop,
  useAddItemToCurrentShop,
  useProducts,
} from "../queries";
import { Meal, Product, Shop } from "../types";
import React, { PropsWithChildren } from "react";
import { Modal } from "../components/Modal";
import clsx from "clsx";
import { ItemSelect } from "../components/ItemSelect";
import { Unit } from "../components/Unit";

export default function HomePage() {
  const mealsQuery = useMeals();
  const currentShopQuery = useCurrentShop();

  if ([mealsQuery, currentShopQuery].some((query) => query.isInitialLoading)) {
    return <p>Loading...</p>;
  }

  const queryWithError = [mealsQuery, currentShopQuery].find(
    (query) => query.isError,
  );

  if (queryWithError && queryWithError.error) {
    return <p>Error: {queryWithError.error.message}</p>;
  }

  const meals = mealsQuery.data;
  const currentShop = currentShopQuery.data as Shop | null;

  return (
    <>
      <div className="mb-7 flex flex-col justify-between sm:flex-row">
        <span className="mb-3 flex space-x-2 sm:mb-0 sm:inline">
          <Link href="/meals/create" className="button w-1/2 text-center">
            🍲 Create meal
          </Link>

          <Link href="/meals/upload" className="button w-1/2 text-center">
            📤 Upload meals
          </Link>
        </span>

        <span>
          <Link
            href="/shop"
            className="button block text-center sm:mb-0 sm:inline"
          >
            🛒 Go shopping
          </Link>
        </span>
      </div>
      <section className="mb-8">
        <div className="w-full">
          <CurrentShop meals={meals} currentShop={currentShop} />
        </div>
      </section>
      <section>
        <Meals meals={meals} currentShop={currentShop} />
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
  const { data: products, isError, isInitialLoading } = useProducts();

  if (isInitialLoading) {
    return <p>Loading products...</p>;
  }

  if (isError) {
    return <p>Error loading products</p>;
  }

  return (
    <>
      {currentShop ? (
        <>
          <div className="mb-2 flex items-center justify-between">
            <h2 className="font-bold">Shop #{currentShop.id}</h2>
            <h3 className="text-xs font-bold">
              {currentShop.meals.length} meals
            </h3>
            <NewShopButton className="button" />
          </div>
          <div className="flex w-full space-x-8">
            <ul className="flex w-1/2 flex-col space-y-1">
              {currentShop.meals.map((meal) => (
                <li key={meal.id} className="flex w-full justify-between">
                  <MealLink meal={meals.find((m) => m.id == meal.id) as Meal} />
                  <RemoveMealFromShopButton mealId={meal.id}>
                    <span className="text-xs">❌</span>
                  </RemoveMealFromShopButton>
                </li>
              ))}
            </ul>

            <div className="w-full">
              <ul className="mb-3 flex flex-col">
                {currentShop.items.map((item) => (
                  <li
                    key={item.productId}
                    className="flex w-full justify-between"
                  >
                    <span>
                      {products.find((i) => i.id === item.productId)?.name}
                    </span>
                    <span>
                      {item.quantity.amount}
                      <Unit quantity={item.quantity} />
                    </span>
                  </li>
                ))}
              </ul>

              <AddItemToShop
                products={products}
                productIdsToExclude={currentShop.items.map((i) => i.productId)}
              />
            </div>
          </div>
        </>
      ) : (
        <div>
          <p className="mb-2">No shop in progress. Get started! </p>
          <NewShopButton className="button" />
        </div>
      )}
    </>
  );
}

function NewShopButton({ className }: { className?: string }) {
  const { mutate } = useStartShop();

  return (
    <Modal
      trigger={(onClick) => (
        <button className={clsx("button", className)} onClick={onClick}>
          🆕 New Shop
        </button>
      )}
      title="New shop"
      body={(close) => (
        <>
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

                close();
              }}
            >
              <button type="submit" className="button">
                Go
              </button>
            </form>

            <button onClick={close} className="underline">
              Cancel
            </button>
          </div>
        </>
      )}
    />
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

function AddItemToShop({
  products,
  productIdsToExclude,
}: {
  products: Product[];
  productIdsToExclude: string[];
}) {
  const { mutate: addItemToShop } = useAddItemToCurrentShop();

  return (
    <ItemSelect
      onItemAdd={addItemToShop}
      products={products}
      productIdsToExclude={productIdsToExclude}
    />
  );
}
