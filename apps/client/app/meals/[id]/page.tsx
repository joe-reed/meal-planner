"use client";

import { useParams } from "next/navigation";
import {
  useProducts,
  useMeal,
  useRemoveIngredientFromMeal,
} from "../../../queries";
import BackButton from "../../../components/BackButton";
import React, { useState } from "react";
import { Product, Meal } from "../../../types";
import { useAddIngredientToMeal } from "../../../queries/useAddIngredientToMeal";
import { Unit } from "../../../components/Unit";
import { useUpdateMeal } from "../../../queries";
import clsx from "clsx";
import { ItemSelect } from "../../../components/ItemSelect";

export default function MealPage() {
  const params = useParams<{ id: string }>();
  const id = params?.id;

  const mealQuery = useMeal(id as string);
  const productsQuery = useProducts();
  const { mutate: addIngredientToMeal } = useAddIngredientToMeal(id as string);
  const { mutate: removeIngredientFromMeal } = useRemoveIngredientFromMeal(
    id as string,
  );

  const meal = mealQuery.data as Meal;
  const products = productsQuery.data as Product[];

  if ([mealQuery, productsQuery].some((query) => query.isInitialLoading)) {
    return <p>Loading...</p>;
  }

  const queryWithError = [mealQuery, productsQuery].find(
    (query) => query.isError,
  );

  if (queryWithError && queryWithError.error) {
    return <p>Error: {queryWithError.error.message}</p>;
  }

  return (
    <div className="flex flex-col">
      <div className="mb-2 flex items-center">
        <BackButton className="mr-3" destination="/" />
        <Name meal={meal} />
      </div>
      <Url meal={meal} className="mb-4 self-start" />
      <h2 className="mb-2 font-bold">Ingredients</h2>
      {(meal.ingredients === null || meal.ingredients.length === 0) && (
        <p className="mb-2">
          No ingredients yet: Add one using the search box below.
        </p>
      )}
      <ul className="mb-6">
        {meal.ingredients.map((ingredient) => (
          <li key={ingredient.id} className="flex justify-between md:w-1/2">
            <span>{products.find((i) => i.id === ingredient.id)?.name}</span>
            <span>
              <span>
                {ingredient.quantity.amount}
                <Unit quantity={ingredient.quantity} />
              </span>
              <button
                onClick={() => removeIngredientFromMeal(ingredient.id)}
                className="ml-2 text-red-500"
              >
                ❌
              </button>
            </span>
          </li>
        ))}
      </ul>
      <ItemSelect
        onItemAdd={({ productId, quantity }) => {
          addIngredientToMeal({ id: productId, quantity });
        }}
        products={products}
        productIdsToExclude={meal.ingredients.map(({ id }) => id)}
        className="w-full md:w-2/3"
      />
    </div>
  );
}

function Url({ meal, className }: { meal: Meal; className?: string }) {
  const { mutate: updateMeal } = useUpdateMeal(meal.id);

  const [url, setUrl] = useState(meal.url);
  const [isEditing, setIsEditing] = useState(false);

  function handleSubmit(e: React.FormEvent) {
    e.preventDefault();
    updateMeal({ ...meal, url });
    setIsEditing(false);
  }

  if (meal.url === "" && !isEditing) {
    return (
      <button
        onClick={() => setIsEditing(true)}
        className={clsx("flex items-center hover:underline", className)}
      >
        <span className="mr-1">➕</span>
        <span>Add URL</span>
      </button>
    );
  }

  return isEditing ? (
    <form
      onSubmit={handleSubmit}
      className={clsx("flex w-full items-center", className)}
    >
      <button onClick={handleSubmit} className="mr-2">
        💾
      </button>
      <input
        type="text"
        value={url}
        onChange={(e) => setUrl(e.target.value)}
        className="-mt-1 w-full rounded-md border py-1 px-2 leading-none"
        autoFocus
      />
    </form>
  ) : (
    <div className={clsx("flex items-center", className)}>
      <button onClick={() => setIsEditing(true)} className="mr-2 text-xs">
        ✏️
      </button>
      <a
        href={meal.url.includes("http") ? meal.url : `https://${meal.url}`}
        className="text-blue-600 hover:underline"
        target="_blank"
      >
        {meal.url}
      </a>
    </div>
  );
}

function Name({ meal, className }: { meal: Meal; className?: string }) {
  const { mutate: updateMeal } = useUpdateMeal(meal.id);

  const [name, setName] = useState(meal.name);
  const [isEditing, setIsEditing] = useState(false);

  function handleSubmit(e: React.FormEvent) {
    e.preventDefault();
    updateMeal({ ...meal, name });
    setIsEditing(false);
  }

  return isEditing ? (
    <form
      onSubmit={handleSubmit}
      className={clsx("flex w-full items-center", className)}
    >
      <button onClick={handleSubmit} className="mr-2">
        💾
      </button>
      <input
        type="text"
        value={name}
        onChange={(e) => setName(e.target.value)}
        className="-mt-1 w-full rounded-md border py-1 px-2 leading-none"
        autoFocus
      />
    </form>
  ) : (
    <div className={clsx("flex", className)}>
      <button onClick={() => setIsEditing(true)} className="mr-2 text-xs">
        ✏️
      </button>
      <h1 className="text-lg font-bold">{meal.name}</h1>
    </div>
  );
}
