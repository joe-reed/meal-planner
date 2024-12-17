import { useMutation, useQueryClient } from "@tanstack/react-query";
import { uuid } from "uuidv4";
import { createIngredient } from "../actions";

export function useCreateIngredient() {
  const queryClient = useQueryClient();

  return useMutation({
    mutationFn: (ingredient: { name: string; category: string }) =>
      createIngredient(JSON.stringify({ ...ingredient, id: uuid() })),
    onSuccess: async () => {
      await queryClient.invalidateQueries({ queryKey: ["ingredients"] });
    },
  });
}
