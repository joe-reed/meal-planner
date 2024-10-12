import { useQuery } from "@tanstack/react-query";
import { Basket } from "../types";

export async function fetchBasket(shopId: string | undefined) {
  const response = await fetch(`/api/baskets/${shopId}`);
  if (!response.ok) {
    throw new Error("Error fetching basket");
  }
  return response.json();
}

export function useBasket(shopId: string | undefined, enabled: boolean) {
  return useQuery<Basket, Error>(
    [`baskets/${shopId}`],
    () => fetchBasket(shopId),
    { enabled },
  );
}
