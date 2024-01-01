import { useQuery } from "react-query";
import { Meal } from "../types/meal";

export async function fetchCurrentShop() {
  const response = await fetch("/api/shops/current");
  if (!response.ok) {
    throw new Error("Error fetching shop");
  }
  return response.json();
}

export default function useCurrentShop() {
  return useQuery<Meal[], Error>("shops/current", fetchCurrentShop);
}
