"use client";

import { useParams } from "next/navigation";
import { z } from "zod";
import {
  useCreateIngredient,
  useIngredients,
  useMeal,
  useRemoveIngredientFromMeal,
} from "../../../queries";
import BackButton from "../../../components/BackButton";
import {
  Combobox,
  ComboboxButton,
  ComboboxInput,
  ComboboxOption,
  ComboboxOptions,
  Select,
} from "@headlessui/react";
import React, { Ref, useRef, useState } from "react";
import { Ingredient, Meal } from "../../../types";
import { useAddIngredientToMeal } from "../../../queries/useAddIngredientToMeal";
import { Modal } from "../../../components/Modal";
import { useCategories } from "../../../queries/useCategories";
import { Unit } from "../../../components/Unit";

type PendingIngredient = {
  id: string;
  quantity: { amount: string; unit: string };
};

export default function MealPage() {
  const params = useParams<{ id: string }>();
  const id = params?.id;

  const mealQuery = useMeal(id as string);
  const ingredientsQuery = useIngredients();
  const { mutate: addIngredientToMeal } = useAddIngredientToMeal(id as string);
  const { mutate: removeIngredientFromMeal } = useRemoveIngredientFromMeal(
    id as string,
  );

  const meal = mealQuery.data as Meal;
  const ingredients = ingredientsQuery.data as Ingredient[];

  const [pendingIngredient, setPendingIngredient] =
    useState<PendingIngredient | null>(null);

  const [ingredientSearchQuery, setIngredientSearchQuery] = useState("");

  const [isAddIngredientModalOpen, setIsAddIngredientModalOpen] =
    useState(false);

  const numberInputRef = useRef<HTMLInputElement>(null);

  const ingredientSearchInputRef = useRef<HTMLInputElement>(null);

  function selectIngredient(ingredient: Ingredient) {
    setPendingIngredient({
      id: ingredient.id,
      quantity: { amount: "", unit: "Number" },
    });
    setIngredientSearchQuery("");
    setTimeout(() => {
      numberInputRef.current?.focus();
    }, 10);
  }

  function addIngredient(pendingIngredient: PendingIngredient) {
    addIngredientToMeal(
      z
        .object({
          id: z.string(),
          quantity: z.object({
            amount: z.coerce.number().positive(),
            unit: z.string(),
          }),
        })
        .parse(pendingIngredient),
    );
    setPendingIngredient(null);

    ingredientSearchInputRef.current?.focus();
    ingredientSearchInputRef.current?.select();
  }

  if ([mealQuery, ingredientsQuery].some((query) => query.isInitialLoading)) {
    return <p>Loading...</p>;
  }

  const queryWithError = [mealQuery, ingredientsQuery].find(
    (query) => query.isError,
  );

  if (queryWithError && queryWithError.error) {
    return <p>Error: {queryWithError.error.message}</p>;
  }

  return (
    <div className="flex flex-col">
      <div className="mb-4 flex items-center">
        <BackButton className="mr-3" destination="/" />
        <h1 className="text-lg font-bold">{meal.name}</h1>
      </div>
      <h2 className="mb-2 font-bold">Ingredients</h2>
      {(meal.ingredients === null || meal.ingredients.length === 0) && (
        <p className="mb-2">
          No ingredients yet: Add one using the search box below.
        </p>
      )}
      <ul className="mb-6">
        {meal.ingredients.map((ingredient) => (
          <li key={ingredient.id} className="flex justify-between md:w-1/2">
            <span>{ingredients.find((i) => i.id === ingredient.id)?.name}</span>
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

      <div className="w-full md:w-2/3">
        {pendingIngredient && (
          <div className="mb-10 flex items-center justify-between space-x-3">
            <div className="whitespace-nowrap">
              {
                ingredients.find(
                  (ingredient) => ingredient.id === pendingIngredient.id,
                )?.name
              }
            </div>
            <div className="flex space-x-1">
              <input
                ref={numberInputRef}
                autoFocus
                type="number"
                value={pendingIngredient.quantity.amount}
                className="button bg-white px-2 py-1"
                size={2}
                onChange={(e) =>
                  setPendingIngredient({
                    ...pendingIngredient,
                    quantity: {
                      ...pendingIngredient.quantity,
                      amount: e.target.value,
                    },
                  })
                }
              />
              <select
                onChange={(e) => {
                  setPendingIngredient({
                    ...pendingIngredient,
                    quantity: {
                      ...pendingIngredient.quantity,
                      unit: e.target.value,
                    },
                  });
                }}
                className="button bg-white px-2 py-1"
              >
                {/*todo: fetch these from api*/}
                <option value="Number">Number</option>
                <option value="Tsp">Tsp</option>
                <option value="Tbsp">Tbsp</option>
                <option value="Cup">Cup</option>
                <option value="Oz">Oz</option>
                <option value="Lb">Lb</option>
                <option value="Gram">Gram</option>
                <option value="Kg">Kg</option>
                <option value="Ml">Ml</option>
                <option value="Litre">Litre</option>
                <option value="Pinch">Pinch</option>
                <option value="Bunch">Bunch</option>
                <option value="Pack">Pack</option>
                <option value="Tin">Tin</option>
              </select>
              <button
                onClick={() => {
                  if (pendingIngredient) {
                    addIngredient(pendingIngredient);
                  }
                }}
                className="button"
              >
                Add
              </button>
            </div>
          </div>
        )}
        <div className="flex items-center">
          <SearchableSelect<Ingredient>
            options={ingredients.filter(
              (ingredient) =>
                !meal.ingredients.some((i) => i.id === ingredient.id),
            )}
            onSelect={selectIngredient}
            onInputChange={(query) => setIngredientSearchQuery(query)}
            inputRef={ingredientSearchInputRef}
          />
          <button
            onClick={() => setIsAddIngredientModalOpen(true)}
            className="ml-2 whitespace-nowrap underline"
          >
            Add new ingredient
          </button>
          <AddNewIngredientModal
            text={ingredientSearchQuery}
            isOpen={isAddIngredientModalOpen}
            setIsOpen={setIsAddIngredientModalOpen}
            onAdd={selectIngredient}
          />
        </div>
      </div>
    </div>
  );
}

type Option = { id: string; name: string };

function SearchableSelect<T extends Option>({
  options,
  onSelect,
  onInputChange,
  inputRef,
}: {
  options: T[];
  onSelect: (option: T) => void;
  onInputChange?: (query: string) => void;
  inputRef?: Ref<HTMLInputElement>;
}) {
  const [query, setQuery] = useState("");

  const filteredOptions =
    query === ""
      ? options
      : options.filter((option) => {
          return option.name.toLowerCase().includes(query.toLowerCase());
        });

  return (
    <Combobox
      onChange={(option: T) => {
        if (!option) return;
        onSelect(option);
        setQuery("");
      }}
      immediate
    >
      <div className="relative w-full">
        <div className="relative w-full cursor-default overflow-hidden rounded-lg bg-white text-left shadow-md focus:outline-none focus-visible:ring-2 focus-visible:ring-white/75 focus-visible:ring-offset-2 focus-visible:ring-offset-teal-300 sm:text-sm">
          <ComboboxInput
            className="w-full border-none py-2 pl-3 pr-10 text-sm leading-5 text-gray-900 focus:ring-0"
            onChange={(event) => {
              setQuery(event.target.value);
              onInputChange && onInputChange(event.target.value);
            }}
            value={query}
            ref={inputRef}
            autoFocus
          />
          <ComboboxButton className="absolute inset-y-0 right-0 flex items-center pr-2">
            <span className="h-5 w-5 text-gray-400" aria-hidden="true">
              ↕️
            </span>
          </ComboboxButton>
        </div>
        <ComboboxOptions
          transition
          className="absolute mt-1 max-h-60 w-full origin-top overflow-auto rounded-md bg-white py-1 text-base shadow-lg ring-1 ring-black/5 transition duration-200 ease-out empty:invisible focus:outline-none data-[closed]:scale-95 data-[closed]:opacity-0 sm:text-sm"
        >
          {filteredOptions.length === 0 ? (
            <div className="relative cursor-default select-none px-4 py-2 text-gray-700">
              <span className="mr-3">No results found</span>
            </div>
          ) : (
            filteredOptions.map((option) => (
              <ComboboxOption
                key={option.id}
                value={option}
                className="ui-active:bg-teal-600 ui-active:text-white ui-not-active:text-gray-900 relative cursor-default select-none py-2 pl-10 pr-4"
              >
                {option.name}
              </ComboboxOption>
            ))
          )}
        </ComboboxOptions>
      </div>
    </Combobox>
  );
}

function AddNewIngredientModal({
  text,
  isOpen,
  setIsOpen,
  onAdd,
}: {
  text: string;
  isOpen: boolean;
  setIsOpen: (value: boolean) => void;
  onAdd: (ingredient: Ingredient) => void;
}) {
  const { mutateAsync } = useCreateIngredient();

  const { data: categories } = useCategories();

  return (
    <>
      <Modal
        isOpen={isOpen}
        setIsOpen={setIsOpen}
        title="Add new ingredient"
        body={(close) => (
          <div className="flex justify-between">
            <form
              onSubmit={async (e) => {
                e.preventDefault();

                const formData = new FormData(e.target as HTMLFormElement);
                const name = formData.get("name") as string;
                const category = formData.get("category") as string;

                const response = await mutateAsync({
                  name,
                  category,
                });

                onAdd(response);

                close();
              }}
            >
              <label className="mb-3 flex flex-col">
                <span>Name</span>
                <input
                  type="text"
                  name="name"
                  required
                  className="rounded-md border py-1 px-2 leading-none"
                  defaultValue={text}
                  data-autofocus
                />
              </label>

              <label className="mb-3 flex flex-col">
                <span>Category</span>
                <Select
                  name="category"
                  aria-label="Ingredient category"
                  className="rounded-md border bg-white py-1 px-2 leading-none"
                >
                  <option value="">Select a category</option>
                  {categories?.map((category) => (
                    <option key={category.name} value={category.name}>
                      {category.name}
                    </option>
                  ))}
                </Select>
              </label>

              <div>
                <button type="submit" className="button mr-3">
                  Create
                </button>

                <button onClick={close} className="underline">
                  Cancel
                </button>
              </div>
            </form>
          </div>
        )}
      />
    </>
  );
}
