import { useMutation } from "@tanstack/react-query";

export function useUploadMeals() {
  return useMutation({
    mutationFn: (meals: File) => {
      const formData = new FormData();
      formData.append("meals", meals);

      return fetch("/api/meals/upload", {
        method: "POST",
        body: formData,
      });
    },
  });
}
