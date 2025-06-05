import { useQuery } from "@tanstack/react-query";
import { Product } from "../types";
import { fetchGroupedProducts } from "../actions";

export function useGroupedProducts() {
  return useQuery<{ string: Product[] }, Error>(
    ["grouped-products"],
    fetchGroupedProducts,
  );
}
