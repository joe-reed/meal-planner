import { useMutation, useQueryClient } from "@tanstack/react-query";
import { v4 as uuid } from "uuid";
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
