import { useMutation, useQueryClient } from "react-query";

export default function useCreateMeal() {
  const queryClient = useQueryClient();

  return useMutation({
    mutationFn: () => {
      return fetch("/api/shops", {
        method: "POST",
      });
    },
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ["shops/current"] });
    },
  });
}
