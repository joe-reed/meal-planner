import { useMutation, useQueryClient } from "@tanstack/react-query";
import { addItemToCurrentShop } from "../actions";

export function useAddItemToCurrentShop() {
  const queryClient = useQueryClient();

  return useMutation({
    mutationFn: (body: {
      productId: string;
      quantity: { amount: number; unit: string };
    }) => addItemToCurrentShop(JSON.stringify(body)),
    onSuccess: async () => {
      await queryClient.invalidateQueries({ queryKey: ["shops/current"] });
    },
  });
}
