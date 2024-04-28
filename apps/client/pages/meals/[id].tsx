import { useRouter } from "next/router";
import { useMeal } from "../../queries";

export default function Meal() {
  const {
    query: { id },
  } = useRouter();

  const {
    isInitialLoading,
    isError,
    data: meal,
    error,
  } = useMeal(id as string);

  if (isInitialLoading) {
    return <p>Loading...</p>;
  }

  if (isError) {
    return <p>Error: {error.message}</p>;
  }

  return <h1>{meal?.name}</h1>;
}
