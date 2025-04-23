import { useMutation, useQueryClient } from "@tanstack/react-query";
import { updateMeal } from "../actions";

export function useUpdateMeal(mealId: string) {
  const queryClient = useQueryClient();

  return useMutation({
    mutationFn: (body: { url: string }) =>
      updateMeal(mealId, JSON.stringify(body)),
    onSuccess: async () => {
      await queryClient.invalidateQueries({ queryKey: [`meals/${mealId}`] });
    },
  });
}
