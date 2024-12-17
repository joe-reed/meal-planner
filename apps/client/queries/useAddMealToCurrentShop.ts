import { useMutation, useQueryClient } from "@tanstack/react-query";
import { addMealToCurrentShop } from "../actions";

export function useAddMealToCurrentShop(mealId: string) {
  const queryClient = useQueryClient();

  return useMutation({
    mutationFn: () => addMealToCurrentShop(JSON.stringify({ id: mealId })),
    onSuccess: async () => {
      await queryClient.invalidateQueries({ queryKey: ["shops/current"] });
    },
  });
}
