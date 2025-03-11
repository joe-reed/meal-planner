import { useMutation } from "@tanstack/react-query";

export function useUploadMeals() {
  return useMutation({
    mutationFn: async (meals: File) => {
      const formData = new FormData();
      formData.append("meals", meals);

      const response = await fetch("/api/meals/upload", {
        method: "POST",
        body: formData,
      });

      if (response.status === 400) {
        return {
          status: response.status,
          error: { message: "validation error", data: await response.json() },
        };
      }

      return { status: response.status, error: null };
    },
  });
}
