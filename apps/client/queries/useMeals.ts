import { useQuery } from "@tanstack/react-query";
import { Meal } from "../types";

export async function fetchMeals() {
  const response = await fetch("/api/meals");
  if (!response.ok) {
    throw new Error("Error fetching meals");
  }
  return response.json();
}

export function useMeals() {
  return useQuery<Meal[], Error>(["meals"], fetchMeals);
}
