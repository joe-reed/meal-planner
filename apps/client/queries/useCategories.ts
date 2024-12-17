import { useQuery } from "@tanstack/react-query";
import { Category } from "../types";
import { fetchCategories } from "../actions";

export function useCategories() {
  return useQuery<Category[], Error>(["categories"], fetchCategories);
}
