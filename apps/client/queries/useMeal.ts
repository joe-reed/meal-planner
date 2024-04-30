import { useQuery } from "@tanstack/react-query";
import { Meal } from "../types/meal";

export async function fetchMeal(id: string) {
  const response = await fetch(`/api/meals/${id}`);
  if (!response.ok) {
    throw new Error("Error fetching meal");
  }
  return response.json();
}

export function useMeal(id: string) {
  return useQuery<Meal, Error>([`meals/${id}`], () => fetchMeal(id));
}
