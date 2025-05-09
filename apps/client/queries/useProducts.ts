import { useQuery } from "@tanstack/react-query";
import { Product } from "../types";
import { fetchProducts } from "../actions";

export function useProducts() {
  return useQuery<Product[], Error>(["products"], fetchProducts);
}
