import { useMutation, useQueryClient } from "@tanstack/react-query";

export function useAddIngredientToMeal(mealId: string) {
  const queryClient = useQueryClient();

  return useMutation({
    mutationFn: (body: {
      id: string;
      quantity: { amount: number; unit: string };
    }) => {
      return fetch(`/api/meals/${mealId}/ingredients`, {
        method: "POST",
        body: JSON.stringify(body),
      });
    },
    onSuccess: async () => {
      await queryClient.invalidateQueries({ queryKey: [`meals/${mealId}`] });
    },
  });
}
