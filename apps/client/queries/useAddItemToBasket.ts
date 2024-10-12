import { useMutation, useQueryClient } from "@tanstack/react-query";

export function useAddItemToBasket(shopId: string) {
  const queryClient = useQueryClient();

  return useMutation({
    mutationFn: (body: { ingredientId: string }) => {
      return fetch(`/api/baskets/${shopId}/items`, {
        method: "POST",
        body: JSON.stringify(body),
      });
    },
    onSuccess: async () => {
      await queryClient.invalidateQueries({ queryKey: [`baskets/${shopId}`] });
    },
  });
}
