import { useMutation, useQueryClient } from "@tanstack/react-query";

export function useRemoveItemFromBasket(shopId: string) {
  const queryClient = useQueryClient();

  return useMutation({
    mutationFn: (ingredientId: string) => {
      return fetch(`/api/baskets/${shopId}/items/${ingredientId}`, {
        method: "DELETE",
      });
    },
    onSuccess: async () => {
      await queryClient.invalidateQueries({ queryKey: [`baskets/${shopId}`] });
    },
  });
}
