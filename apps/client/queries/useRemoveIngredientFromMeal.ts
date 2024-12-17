import { useMutation, useQueryClient } from "@tanstack/react-query";
import { removeIngredientFromMeal } from "../actions";

export function useRemoveIngredientFromMeal(mealId: string) {
  const queryClient = useQueryClient();

  return useMutation({
    mutationFn: (ingredientId: string) =>
      removeIngredientFromMeal(mealId, ingredientId),
    onSuccess: async () => {
      await queryClient.invalidateQueries({ queryKey: [`meals/${mealId}`] });
    },
  });
}
