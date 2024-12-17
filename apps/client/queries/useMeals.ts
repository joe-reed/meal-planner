import { useQuery } from "@tanstack/react-query";
import { Meal } from "../types";
import { fetchMeals } from "../actions";

export function useMeals() {
  return useQuery<Meal[], Error>(["meals"], fetchMeals);
}
