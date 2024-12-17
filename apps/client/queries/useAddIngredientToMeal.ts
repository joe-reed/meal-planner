import { useMutation, useQueryClient } from "@tanstack/react-query";
import { addIngredientToMeal } from "../actions";

export function useAddIngredientToMeal(mealId: string) {
  const queryClient = useQueryClient();

  return useMutation({
    mutationFn: (body: {
      id: string;
      quantity: { amount: number; unit: string };
    }) => addIngredientToMeal(mealId, JSON.stringify(body)),
    onSuccess: async () => {
      await queryClient.invalidateQueries({ queryKey: [`meals/${mealId}`] });
    },
  });
}
