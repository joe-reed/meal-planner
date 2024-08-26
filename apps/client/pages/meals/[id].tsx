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
  Transition,
} from "@headlessui/react";
import React, { Fragment, useState } from "react";
import { Ingredient, Meal } from "../../types";
import { useAddIngredientToMeal } from "../../queries/useAddIngredientToMeal";
import { Modal } from "../../components/Modal";

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

  if ([mealQuery, ingredientsQuery].some((query) => query.isInitialLoading)) {
    return <p>Loading...</p>;
  }

  const queryWithError = [mealQuery, ingredientsQuery].find(
    (query) => query.isError,
  );

  if (queryWithError && queryWithError.error) {
    return <p>Error: {queryWithError.error.message}</p>;
  }

  const meal = mealQuery.data as Meal;
  const ingredients = ingredientsQuery.data as Ingredient[];

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
      <ul>
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
        <SearchableSelect
          options={ingredients.filter(
            (ingredient) =>
              !meal.ingredients.some((i) => i.id === ingredient.id),
          )}
          onChange={(ingredient) => {
            addIngredientToMeal(ingredient.id);
          }}
          emptyUi={(query) => <AddNewIngredientModal text={query} />}
        />
      </div>
    </div>
  );
}

type Option = { id: string; name: string };

function SearchableSelect({
  options,
  onChange,
  emptyUi,
}: {
  options: Option[];
  onChange: (option: Option) => void;
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
        setQuery("");
        onChange(option);
      }}
      immediate
    >
      <div className="relative w-full">
        <div className="relative w-full cursor-default overflow-hidden rounded-lg bg-white text-left shadow-md focus:outline-none focus-visible:ring-2 focus-visible:ring-white/75 focus-visible:ring-offset-2 focus-visible:ring-offset-teal-300 sm:text-sm">
          <ComboboxInput
            className="w-full border-none py-2 pl-3 pr-10 text-sm leading-5 text-gray-900 focus:ring-0"
            onChange={(event) => setQuery(event.target.value)}
            value={query}
            autoFocus
          />
          <ComboboxButton className="absolute inset-y-0 right-0 flex items-center pr-2">
            <span className="h-5 w-5 text-gray-400" aria-hidden="true">
              ↕️
            </span>
          </ComboboxButton>
        </div>
        <Transition
          as={Fragment}
          leave="transition ease-in duration-100"
          leaveFrom="opacity-100"
          leaveTo="opacity-0"
          afterLeave={() => setQuery("")}
        >
          <ComboboxOptions className="absolute mt-1 max-h-60 w-full overflow-auto rounded-md bg-white py-1 text-base shadow-lg ring-1 ring-black/5 focus:outline-none sm:text-sm">
            {filteredOptions.length === 0 ? (
              <div className="relative cursor-default select-none px-4 py-2 text-gray-700">
                <span className="mr-3">No options</span>
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
        </Transition>
      </div>
    </Combobox>
  );
}

function AddNewIngredientModal({ text }: { text: string }) {
  const { mutateAsync } = useCreateIngredient();

  return (
    <Modal
      trigger={(onClick) => (
        <button onClick={onClick} className="underline">
          Add new ingredient
        </button>
      )}
      title="Add new ingredient"
      body={(close) => (
        <div className="flex justify-between px-6">
          <form
            onSubmit={async (e) => {
              e.preventDefault();

              const formData = new FormData(e.target as HTMLFormElement);
              const name = formData.get("name") as string;

              await mutateAsync({
                name,
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
                autoFocus
              />
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
  );
}
