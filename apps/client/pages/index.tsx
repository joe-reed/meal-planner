import Link from "next/link";
import useMeals from "../queries/useMeals";
import useCurrentShop from "../queries/useCurrentShop";
import useStartShop from "../queries/useStartShop";
import { Meal } from "../types/meal";
import useAddMealToCurrentShop from "../queries/useAddMealToCurrentShop";

export default function Index() {
  const { isInitialLoading, isError, data: meals, error } = useMeals();

  if (isInitialLoading) {
    return <p>Loading...</p>;
  }

  if (isError) {
    return <p>Error: {error.message}</p>;
  }

  return (
    <>
      <nav className="flex justify-end mb-2">
        <Link href="/meals/create" className="button">
          Create meal
        </Link>
        <StartShopButton />
      </nav>
      <section>
        <Meals meals={meals || []} />
      </section>
      <section>
        <CurrentShop meals={meals || []} />
      </section>
    </>
  );
}

function Meals({ meals }: { meals: Meal[] }) {
  return (
    <>
      <h2>Meals</h2>
      <ul className="flex space-x-2">
        {meals?.map((meal) => (
          <li key={meal.id} className="border px-3 py-1 rounded-lg">
            <Link href={`/meals/${meal.id}`}>{meal.name}</Link>
            <AddMealToShopButton mealId={meal.id} />
          </li>
        ))}
      </ul>
    </>
  );
}

function CurrentShop({ meals }: { meals: Meal[] }) {
  const {
    isInitialLoading,
    isError,
    data: currentShop,
    error,
  } = useCurrentShop();

  if (isInitialLoading) {
    return <p>Loading...</p>;
  }

  if (isError) {
    return <p>Error: {error.message}</p>;
  }

  return (
    <>
      <section>
        {currentShop ? (
          <>
            <h2>Current shop</h2>
            <span>Shop number {currentShop.id}</span>
            <ul>
              {currentShop.meals.map((meal) => (
                <li key={meal.id}>
                  {meals.find((m) => m.id == meal.id)?.name}
                </li>
              ))}
            </ul>
          </>
        ) : null}
      </section>
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
    <div>
      <form
        onSubmit={(e) => {
          e.preventDefault();

          mutate();
        }}
      >
        <button type="submit" className="button">
          +
        </button>
      </form>
    </div>
  );
}
