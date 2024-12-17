import { useQuery } from "@tanstack/react-query";
import { Meal } from "../types";
import { fetchMeal } from "../actions";

export function useMeal(id: string) {
  return useQuery<Meal, Error>([`meals/${id}`], () => fetchMeal(id));
}
