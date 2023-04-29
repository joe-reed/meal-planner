import { useQuery } from 'react-query';

type Meal = {
  id: string;
  name: string;
};

async function fetchMeals() {
  const response = await fetch('/api/meals');
  if (!response.ok) {
    throw new Error('Error fetching meals');
  }
  return response.json();
}

function useMeals() {
  return useQuery<Meal[], Error>('meals', fetchMeals);
}

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
