import { useRouter } from "next/router";
import { useIngredients, useMeal } from "../../queries";
import BackButton from "../../components/BackButton";
import { Combobox, Transition } from "@headlessui/react";
import { Fragment, useState } from "react";
import { Ingredient, Meal } from "../../types";
import { useAddIngredientToMeal } from "../../queries/useAddIngredientToMeal";

export default function MealPage() {
  const {
    query: { id },
  } = useRouter();

  const mealQuery = useMeal(id as string);
  const ingredientsQuery = useIngredients();
  const { mutate: addIngredientToMeal } = useAddIngredientToMeal(id as string);

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
      <div className="mb-2 flex items-center">
        <BackButton className="mr-3" />
        <h1 className="font-bold">{meal.name}</h1>
      </div>
      <h2 className="text-lg font-bold">Ingredients</h2>
      <ul>
        {meal.ingredients.map((ingredient) => (
          <li key={ingredient.id}>
            {ingredients.find((i) => i.id === ingredient.id)?.name}
          </li>
        ))}
      </ul>
      <SearchableSelect
        options={ingredients.filter(
          (ingredient) => !meal.ingredients.some((i) => i.id === ingredient.id),
        )}
        onChange={(ingredient) => addIngredientToMeal(ingredient.id)}
      />
    </div>
  );
}

type Option = { id: string; name: string };

function SearchableSelect({
  options,
  onChange,
}: {
  options: Option[];
  onChange: (option: Option) => void;
}) {
  const [query, setQuery] = useState("");

  const filteredOptions =
    query === ""
      ? options
      : options.filter((option) => {
          return option.name.toLowerCase().includes(query.toLowerCase());
        });

  return (
    <Combobox onChange={onChange}>
      <div className="relative mt-1">
        <div className="relative w-full cursor-default overflow-hidden rounded-lg bg-white text-left shadow-md focus:outline-none focus-visible:ring-2 focus-visible:ring-white/75 focus-visible:ring-offset-2 focus-visible:ring-offset-teal-300 sm:text-sm">
          <Combobox.Input
            className="w-full border-none py-2 pl-3 pr-10 text-sm leading-5 text-gray-900 focus:ring-0"
            onChange={(event) => setQuery(event.target.value)}
            displayValue={(option: Option) => option.name}
          />
          <Combobox.Button className="absolute inset-y-0 right-0 flex items-center pr-2">
            <span className="h-5 w-5 text-gray-400" aria-hidden="true">
              ↕️
            </span>
          </Combobox.Button>
        </div>
        <Transition
          as={Fragment}
          leave="transition ease-in duration-100"
          leaveFrom="opacity-100"
          leaveTo="opacity-0"
          afterLeave={() => setQuery("")}
        >
          <Combobox.Options className="absolute mt-1 max-h-60 w-full overflow-auto rounded-md bg-white py-1 text-base shadow-lg ring-1 ring-black/5 focus:outline-none sm:text-sm">
            {filteredOptions.length === 0 ? (
              <div className="relative cursor-default select-none px-4 py-2 text-gray-700">
                No options
              </div>
            ) : (
              filteredOptions.map((option) => (
                <Combobox.Option
                  key={option.id}
                  value={option}
                  className="ui-active:bg-teal-600 ui-active:text-white ui-not-active:text-gray-900 relative cursor-default select-none py-2 pl-10 pr-4"
                >
                  {option.name}
                </Combobox.Option>
              ))
            )}
          </Combobox.Options>
        </Transition>
      </div>
    </Combobox>
  );
}

/*
import {
  Fragment, useState;
} from 'react'
import {
  Combobox, Transition;
} from '@headlessui/react'
import {
  CheckIcon, ChevronUpDownIcon;
} from '@heroicons/react/20/solid'

const people = [
{
  id: 1, name;
:
  "Wade Cooper";
},
{
  id: 2, name;
:
  "Arlene Mccoy";
},
{
  id: 3, name;
:
  "Devon Webb" },
  { id: 4, name: 'Tom Cook' },
  { id: 5, name: 'Tanya Fox' },
  { id: 6, name: 'Hellen Schmidt' },
]

export default function Example() {
  const [selected, setSelected] = useState(people[0])
  const [query, setQuery] = useState('')

  const filteredPeople =
    query === ''
      ? people
      : people.filter((person) =>
          person.name
            .toLowerCase()
            .replace(/\s+/g, '')
            .includes(query.toLowerCase().replace(/\s+/g, ''))
        )

  return (
    <div className="fixed top-16 w-72">
      <Combobox value={selected} onChange={setSelected}>
        <div className="relative mt-1">
          <div className="relative w-full cursor-default overflow-hidden rounded-lg bg-white text-left shadow-md focus:outline-none focus-visible:ring-2 focus-visible:ring-white/75 focus-visible:ring-offset-2 focus-visible:ring-offset-teal-300 sm:text-sm">
            <Combobox.Input
              className="w-full border-none py-2 pl-3 pr-10 text-sm leading-5 text-gray-900 focus:ring-0"
              displayValue={(person) => person.name}
              onChange={(event) => setQuery(event.target.value)}
            />
            <Combobox.Button className="absolute inset-y-0 right-0 flex items-center pr-2">
              <ChevronUpDownIcon
                className="h-5 w-5 text-gray-400"
                aria-hidden="true"
              />
            </Combobox.Button>
          </div>
          <Transition
            as={Fragment}
            leave="transition ease-in duration-100"
            leaveFrom="opacity-100"
            leaveTo="opacity-0"
            afterLeave={() => setQuery('')}
          >
            <Combobox.Options className="absolute mt-1 max-h-60 w-full overflow-auto rounded-md bg-white py-1 text-base shadow-lg ring-1 ring-black/5 focus:outline-none sm:text-sm">
              {filteredPeople.length === 0 && query !== '' ? (
                <div className="relative cursor-default select-none px-4 py-2 text-gray-700">
                  Nothing found.
                </div>
              ) : (
                filteredPeople.map((person) => (
                  <Combobox.Option
                    key={person.id}
                    className={({ active }) =>
                      `relative cursor-default select-none py-2 pl-10 pr-4 ${
                        active ? 'bg-teal-600 text-white' : 'text-gray-900'
                      }`
                    }
                    value={person}
                  >
                    {({ selected, active }) => (
                      <>
                        <span
                          className={`block truncate ${
                            selected ? 'font-medium' : 'font-normal'
                          }`}
                        >
                          {person.name}
                        </span>
                        {selected ? (
                          <span
                            className={`absolute inset-y-0 left-0 flex items-center pl-3 ${
                              active ? 'text-white' : 'text-teal-600'
                            }`}
                          >
                            <CheckIcon className="h-5 w-5" aria-hidden="true" />
                          </span>
                        ) : null}
                      </>
                    )}
                  </Combobox.Option>
                ))
              )}
            </Combobox.Options>
          </Transition>
        </div>
      </Combobox>
    </div>
  )
}

 */
