import { useMutation, useQueryClient } from "@tanstack/react-query";
import { removeItemFromBasket } from "../actions";

export function useRemoveItemFromBasket(shopId: string) {
  const queryClient = useQueryClient();

  return useMutation({
    mutationFn: (ingredientId: string) =>
      removeItemFromBasket(shopId, ingredientId),
    onSuccess: async () => {
      await queryClient.invalidateQueries({ queryKey: ["shopping-list"] });
    },
  });
}
