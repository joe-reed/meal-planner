import { useQuery } from "@tanstack/react-query";
import { Shop } from "../types";

export async function fetchCurrentShop() {
  const response = await fetch("/api/shops/current");
  if (!response.ok) {
    throw new Error("Error fetching shop");
  }
  return response.json();
}

export function useCurrentShop() {
  return useQuery<Shop, Error>(["shops/current"], fetchCurrentShop);
}
