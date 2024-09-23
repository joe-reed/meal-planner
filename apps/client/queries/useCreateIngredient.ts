import { useMutation, useQueryClient } from "@tanstack/react-query";
import { uuid } from "uuidv4";

export function useCreateIngredient() {
  const queryClient = useQueryClient();

  return useMutation({
    mutationFn: (ingredient: { name: string; category: string }) => {
      return fetch("/api/ingredients", {
        method: "POST",
        body: JSON.stringify({ ...ingredient, id: uuid() }),
      });
    },
    onSuccess: async () => {
      await queryClient.invalidateQueries({ queryKey: ["ingredients"] });
    },
  });
}
