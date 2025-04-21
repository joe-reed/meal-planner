import { useMutation } from "@tanstack/react-query";
import { v4 as uuid } from "uuid";
import { createMeal } from "../actions";

export function useCreateMeal() {
  return useMutation({
    mutationFn: (meal: { name: string; url: string }) =>
      createMeal(JSON.stringify({ ...meal, id: uuid() })),
  });
}
