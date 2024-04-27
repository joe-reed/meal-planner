import { useMutation, useQueryClient } from "react-query";

export default function useAddMealToCurrentShop(mealId: string) {
  const queryClient = useQueryClient();

  return useMutation({
    mutationFn: () => {
      return fetch("/api/shops/current/meals", {
        method: "POST",
        body: JSON.stringify({ id: mealId }),
      });
    },
    onSuccess: async () => {
      await queryClient.invalidateQueries({ queryKey: ["shops/current"] });
    },
  });
}
