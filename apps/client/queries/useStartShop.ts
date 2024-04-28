import { useMutation, useQueryClient } from "@tanstack/react-query";

export default function useCreateMeal() {
  const queryClient = useQueryClient();

  return useMutation({
    mutationFn: () => {
      return fetch("/api/shops", {
        method: "POST",
      });
    },
    onSuccess: async () => {
      await queryClient.invalidateQueries({ queryKey: ["shops/current"] });
    },
  });
}
