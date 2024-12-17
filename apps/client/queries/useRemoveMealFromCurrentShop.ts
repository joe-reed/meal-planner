import { useMutation, useQueryClient } from "@tanstack/react-query";
import { removeMealFromCurrentShop } from "../actions";

export function useRemoveMealFromCurrentShop(mealId: string) {
  const queryClient = useQueryClient();

  return useMutation({
    mutationFn: () => removeMealFromCurrentShop(mealId),
    onSuccess: async () => {
      await queryClient.invalidateQueries({ queryKey: ["shops/current"] });
    },
  });
}
