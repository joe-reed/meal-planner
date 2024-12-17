import { useQuery } from "@tanstack/react-query";
import { Basket } from "../types";
import { fetchBasket } from "../actions";

export function useBasket(shopId: string | undefined, enabled: boolean) {
  return useQuery<Basket, Error>(
    [`baskets/${shopId}`],
    () => fetchBasket(shopId),
    { enabled },
  );
}
