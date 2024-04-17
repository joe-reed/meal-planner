import Link from "next/link";
import useMeals from "../queries/useMeals";
import useCurrentShop from "../queries/useCurrentShop";

export default function Index() {
  return (
    <>
      <nav className="flex justify-end mb-2">
        <Link href="/meals/create" className="button">
          Create meal
        </Link>
      </nav>
      <section>
        <Meals />
      </section>
      <section>
        <CurrentShop />
      </section>
    </>
  );
}

function Meals() {
  const { isLoading, isError, data: meals, error } = useMeals();

  if (isLoading) {
    return <p>Loading...</p>;
  }

  if (isError) {
    return <p>Error: {error.message}</p>;
  }

  return (
    <>
      <h2>Meals</h2>
      <ul className="flex space-x-2">
        {meals?.map((meal) => (
          <li key={meal.id} className="border px-3 py-1 rounded-lg">
            <Link href={`/meals/${meal.id}`}>{meal.name}</Link>
          </li>
        ))}
      </ul>
    </>
  );
}

function CurrentShop() {
  const { isLoading, isError, data: meals, error } = useCurrentShop();

  if (isLoading) {
    return <p>Loading...</p>;
  }

  if (isError) {
    return <p>Error: {error.message}</p>;
  }

  return (
    <>
      <section>
        <h2>Current shop</h2>
        <ul>
          {meals?.map((meal) => (
            <li key={meal.id}>{meal.name}</li>
          ))}
        </ul>
      </section>
    </>
  );
}
