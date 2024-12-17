import { useMutation } from "@tanstack/react-query";
import { uuid } from "uuidv4";
import { createMeal } from "../actions";

export function useCreateMeal() {
  return useMutation({
    mutationFn: (meal: { name: string }) =>
      createMeal(JSON.stringify({ ...meal, id: uuid() })),
  });
}
