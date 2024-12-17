import { useMutation, useQueryClient } from "@tanstack/react-query";
import { addItemToBasket } from "../actions";

export function useAddItemToBasket(shopId: string) {
  const queryClient = useQueryClient();

  return useMutation({
    mutationFn: (body: { ingredientId: string }) =>
      addItemToBasket(shopId, JSON.stringify(body)),
    onSuccess: async () => {
      await queryClient.invalidateQueries({ queryKey: [`baskets/${shopId}`] });
    },
  });
}
