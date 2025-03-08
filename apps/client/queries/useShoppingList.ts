import { useQuery } from "@tanstack/react-query";
import { fetchShoppingList } from "../actions";

export function useShoppingList() {
  return useQuery<{ shopId: string; shoppingList: { string: any } }, Error>(
    ["shopping-list"],
    fetchShoppingList,
  );
}
