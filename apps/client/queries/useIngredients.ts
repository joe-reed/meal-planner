import { useQuery } from "@tanstack/react-query";
import { Ingredient } from "../types";

export async function fetchIngredients() {
  const response = await fetch("/api/ingredients");
  if (!response.ok) {
    throw new Error("Error fetching ingredients");
  }
  return response.json();
}

export function useIngredients() {
  return useQuery<Ingredient[], Error>(["ingredients"], fetchIngredients);
}
