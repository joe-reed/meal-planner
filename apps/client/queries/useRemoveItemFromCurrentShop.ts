import { useMutation, useQueryClient } from "@tanstack/react-query";
import { removeItemFromCurrentShop } from "../actions";

export function useRemoveItemFromCurrentShop(productId: string) {
  const queryClient = useQueryClient();

  return useMutation({
    mutationFn: () => removeItemFromCurrentShop(productId),
    onSuccess: async () => {
      await queryClient.invalidateQueries({ queryKey: ["shops/current"] });
    },
  });
}
