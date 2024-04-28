import { useRouter } from "next/router";
import { useMeal } from "../../queries";
import BackButton from "../../components/BackButton";

export default function Meal() {
  const {
    query: { id },
    back,
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

  return (
    <div className="flex">
      <BackButton className="mr-3" />
      <h1>{meal?.name}</h1>
    </div>
  );
}
