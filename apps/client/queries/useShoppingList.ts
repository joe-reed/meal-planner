import { useQuery } from "@tanstack/react-query";
import { fetchShoppingList } from "../actions";
import { Product } from "../types/product";

type ShoppingListItem = Product & {
  mealCount: number;
  isInBasket: boolean;
  quantities: { unit: string; amount: number }[];
};

export function useShoppingList() {
  return useQuery<{ shopId: string; shoppingList: Record<string, ShoppingListItem> }, Error>(
    ["shopping-list"],
    fetchShoppingList,
  );
}
