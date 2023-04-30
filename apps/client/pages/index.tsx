import useMeals from '../queries/useMeals';

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
      <h1>Meal planner</h1>
      {meals?.map((meal) => (
        <p key={meal.id}>{meal.name}</p>
      ))}
    </>
  );
}
