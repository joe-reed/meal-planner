import { useMutation, useQueryClient } from "@tanstack/react-query";
import { startShop } from "../actions";

export function useStartShop() {
  const queryClient = useQueryClient();

  return useMutation({
    mutationFn: startShop,
    onSuccess: async () => {
      await queryClient.invalidateQueries({ queryKey: ["shops/current"] });
    },
  });
}
