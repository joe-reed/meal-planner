import { useRouter } from "next/router";
import {
  useCreateIngredient,
  useIngredients,
  useMeal,
  useRemoveIngredientFromMeal,
} from "../../queries";
import BackButton from "../../components/BackButton";
import {
  Combobox,
  ComboboxButton,
  ComboboxInput,
  ComboboxOption,
  ComboboxOptions,
  Select,
} from "@headlessui/react";
import React, { useState } from "react";
import { Ingredient, Meal } from "../../types";
import { useAddIngredientToMeal } from "../../queries/useAddIngredientToMeal";
import { Modal } from "../../components/Modal";
import { useCategories } from "../../queries/useCategories";

export default function MealPage() {
  const {
    query: { id },
  } = useRouter();

  const mealQuery = useMeal(id as string);
  const ingredientsQuery = useIngredients();
  const { mutate: addIngredientToMeal } = useAddIngredientToMeal(id as string);
  const { mutate: removeIngredientFromMeal } = useRemoveIngredientFromMeal(
    id as string,
  );

  const meal = mealQuery.data as Meal;
  const ingredients = ingredientsQuery.data as Ingredient[];

  const [pendingIngredient, setPendingIngredient] = useState<{
    id: string;
    quantity: { amount: number; unit: string };
  } | null>(null);

  const [ingredientSearchQuery, setIngredientSearchQuery] = useState("");

  const [isAddIngredientModalOpen, setIsAddIngredientModalOpen] =
    useState(false);

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
          <li key={ingredient.id} className="flex w-1/2 justify-between">
            <span>{ingredients.find((i) => i.id === ingredient.id)?.name}</span>
            <span>
              <span>
                {ingredient.quantity.amount}
                {ingredient.quantity.unit !== "Number"
                  ? " " +
                    ingredient.quantity.unit +
                    (ingredient.quantity.amount > 1 ? "s" : "")
                  : "x"}
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

      <div className="w-full md:w-1/2">
        {pendingIngredient && (
          <div className="flex items-center justify-between space-x-3">
            <div className="whitespace-nowrap">
              {
                ingredients.find(
                  (ingredient) => ingredient.id === pendingIngredient.id,
                )?.name
              }
            </div>
            <div className="flex space-x-1">
              <input
                type="number"
                value={pendingIngredient.quantity.amount}
                className="px-2 py-1"
                size={2}
                onChange={(e) =>
                  setPendingIngredient({
                    ...pendingIngredient,
                    quantity: {
                      ...pendingIngredient.quantity,
                      amount: parseInt(e.target.value),
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
                <option value="Cup">Cup</option>
                <option value="Tsp">Tsp</option>
                <option value="Tbsp">Tbsp</option>
                <option value="Oz">Oz</option>
                <option value="Lb">Lb</option>
                <option value="Gram">Gram</option>
                <option value="Kg">Kg</option>
              </select>
              <button
                onClick={() => {
                  if (pendingIngredient) {
                    addIngredientToMeal(pendingIngredient);
                    setPendingIngredient(null);
                  }
                }}
                className="button"
              >
                Add
              </button>
            </div>
          </div>
        )}
        <SearchableSelect
          options={ingredients.filter(
            (ingredient) =>
              !meal.ingredients.some((i) => i.id === ingredient.id),
          )}
          onSelect={(ingredient) => {
            setPendingIngredient({
              id: ingredient.id,
              quantity: { amount: 1, unit: "Number" },
            });
          }}
          onInputChange={(query) => setIngredientSearchQuery(query)}
          emptyUi={() => (
            <button
              onClick={() => setIsAddIngredientModalOpen(true)}
              className="underline"
            >
              Add new ingredient
            </button>
          )}
        />
        <AddNewIngredientModal
          text={ingredientSearchQuery}
          isOpen={isAddIngredientModalOpen}
          setIsOpen={setIsAddIngredientModalOpen}
        />
      </div>
    </div>
  );
}

type Option = { id: string; name: string };

function SearchableSelect({
  options,
  onSelect,
  onInputChange,
  emptyUi,
}: {
  options: Option[];
  onSelect: (option: Option) => void;
  onInputChange?: (query: string) => void;
  emptyUi?: (query: string) => React.ReactNode;
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
      onChange={(option: Option) => {
        if (!option) return;
        onSelect(option);
      }}
      onClose={() => {
        setQuery("");
        onInputChange && onInputChange("");
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
              {emptyUi && emptyUi(query)}
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
}: {
  text: string;
  isOpen: boolean;
  setIsOpen: (value: boolean) => void;
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
          <div className="flex justify-between px-6">
            <form
              onSubmit={async (e) => {
                e.preventDefault();

                const formData = new FormData(e.target as HTMLFormElement);
                const name = formData.get("name") as string;
                const category = formData.get("category") as string;

                await mutateAsync({
                  name,
                  category,
                });

                close();
              }}
            >
              <label className="mb-3 block">
                <span className="mr-2">Name</span>
                <input
                  type="text"
                  name="name"
                  required
                  className="rounded-md border py-1 px-2 leading-none"
                  defaultValue={text}
                />
              </label>

              <label className="mb-3 block">
                <span className="mr-2">Category</span>
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
