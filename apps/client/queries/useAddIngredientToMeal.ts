import { useMutation, useQueryClient } from "@tanstack/react-query";

export function useAddIngredientToMeal(mealId: string) {
  const queryClient = useQueryClient();

  return useMutation({
    mutationFn: (ingredientId: string) => {
      return fetch(`/api/meals/${mealId}/ingredients`, {
        method: "POST",
        body: JSON.stringify({ id: ingredientId }),
      });
    },
    onSuccess: async () => {
      await queryClient.invalidateQueries({ queryKey: [`meals/${mealId}`] });
    },
  });
}
