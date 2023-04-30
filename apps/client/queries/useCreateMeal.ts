import { useMutation } from "react-query";
import { uuid } from "uuidv4";

export default function useCreateMeal() {
  return useMutation({
    mutationFn: (meal: { name: string }) => {
      return fetch("/api/meals", {
        method: "POST",
        body: JSON.stringify({ ...meal, id: uuid() }),
      });
    },
  });
}
