import { useQuery } from "@tanstack/react-query";
import { Ingredient } from "../types";
import { fetchIngredients } from "../actions";

export function useIngredients() {
  return useQuery<Ingredient[], Error>(["ingredients"], fetchIngredients);
}
