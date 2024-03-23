import { useQuery } from "react-query";
import { Shop } from "../types/shop";

export async function fetchCurrentShop() {
  const response = await fetch("/api/shops/current");
  if (!response.ok) {
    throw new Error("Error fetching shop");
  }
  return response.json();
}

export default function useCurrentShop() {
  return useQuery<Shop, Error>("shops/current", fetchCurrentShop);
}
