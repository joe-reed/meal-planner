import { useQuery } from "@tanstack/react-query";
import { Category, Ingredient } from "../types";

export async function fetchCategories() {
  const response = await fetch("/api/categories");
  if (!response.ok) {
    throw new Error("Error fetching categories");
  }
  return response.json();
}

export function useCategories() {
  return useQuery<Category[], Error>(["categories"], fetchCategories);
}
