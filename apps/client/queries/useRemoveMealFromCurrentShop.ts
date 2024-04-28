import { useMutation, useQueryClient } from "@tanstack/react-query";

export function useRemoveMealFromCurrentShop(mealId: string) {
  const queryClient = useQueryClient();

  return useMutation({
    mutationFn: () => {
      return fetch(`/api/shops/current/meals/${mealId}`, {
        method: "DELETE",
      });
    },
    onSuccess: async () => {
      await queryClient.invalidateQueries({ queryKey: ["shops/current"] });
    },
  });
}
