import Link from "next/link";
import { Meal } from "../types/meal";
import {
  useAddMealToCurrentShop,
  useCurrentShop,
  useMeals,
  useRemoveMealFromCurrentShop,
  useStartShop,
} from "../queries";
import { Shop } from "../types/shop";

export default function Index() {
  const mealsQuery = useMeals();
  const currentShopQuery = useCurrentShop();

  if ([mealsQuery, currentShopQuery].some((query) => query.isInitialLoading)) {
    return <p>Loading...</p>;
  }

  const queryWithError = [mealsQuery, currentShopQuery].find(
    (query) => query.isError
  );

  if (queryWithError && queryWithError.error) {
    return <p>Error: {queryWithError.error.message}</p>;
  }

  const meals = mealsQuery.data as Meal[];
  const currentShop = currentShopQuery.data as Shop;

  return (
    <>
      <nav className="flex justify-end mb-2">
        <Link href="/meals/create" className="button">
          Create meal
        </Link>
        <StartShopButton />
      </nav>
      <section className="mb-4">
        <Meals meals={meals} currentShop={currentShop} />
      </section>
      <section>
        <CurrentShop meals={meals} currentShop={currentShop} />
      </section>
    </>
  );
}

function Meals({ meals, currentShop }: { meals: Meal[]; currentShop: Shop }) {
  return (
    <>
      <h2 className="font-bold mb-2">Meals</h2>
      <ul className="flex space-x-2">
        {meals?.map((meal) => (
          <li
            key={meal.id}
            className="border px-3 py-1 rounded-lg flex items-center"
          >
            <Link href={`/meals/${meal.id}`}>{meal.name}</Link>
            <span className="ml-2">
              {currentShop.meals.some((m) => m.id == meal.id) ? (
                <div>✅</div>
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
  currentShop: Shop;
}) {
  return (
    <>
      {currentShop ? (
        <>
          <h2 className="font-bold mb-2">Shop #{currentShop.id}</h2>
          <ul className="flex flex-col space-y-1">
            {currentShop.meals.map((meal) => (
              <li key={meal.id} className="flex w-full justify-between">
                <p>{meals.find((m) => m.id == meal.id)?.name}</p>
                <RemoveMealFromShopButton mealId={meal.id} />
              </li>
            ))}
          </ul>
        </>
      ) : null}
    </>
  );
}

function StartShopButton() {
  const { mutate } = useStartShop();

  return (
    <div>
      <form
        onSubmit={(e) => {
          e.preventDefault();

          mutate();
        }}
      >
        <button type="submit" className="button">
          Start Shop
        </button>
      </form>
    </div>
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

function RemoveMealFromShopButton({ mealId }: { mealId: string }) {
  const { mutate } = useRemoveMealFromCurrentShop(mealId);

  return (
    <form
      onSubmit={(e) => {
        e.preventDefault();

        mutate();
      }}
    >
      <button type="submit">➖</button>
    </form>
  );
}
