import { useMutation, useQueryClient } from "@tanstack/react-query";

export function useRemoveIngredientFromMeal(mealId: string) {
  const queryClient = useQueryClient();

  return useMutation({
    mutationFn: (ingredientId: string) => {
      return fetch(`/api/meals/${mealId}/ingredients/${ingredientId}`, {
        method: "DELETE",
      });
    },
    onSuccess: async () => {
      await queryClient.invalidateQueries({ queryKey: [`meals/${mealId}`] });
    },
  });
}
