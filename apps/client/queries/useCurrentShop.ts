import { useQuery } from "@tanstack/react-query";
import { Shop } from "../types";
import { fetchCurrentShop } from "../actions";

export function useCurrentShop() {
  return useQuery<Shop, Error>(["shops/current"], fetchCurrentShop);
}
