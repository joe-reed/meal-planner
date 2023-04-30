import Link from "next/link";
import useMeals from "../queries/useMeals";

export default function Index() {
  const { isLoading, isError, data: meals, error } = useMeals();

  if (isLoading) {
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
      </nav>
      <ul className="flex space-x-2">
        {meals?.map((meal) => (
          <li key={meal.id} className="border px-3 py-1 rounded-lg">
            {meal.name}
          </li>
        ))}
      </ul>
    </>
  );
}
