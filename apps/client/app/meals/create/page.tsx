"use client";

import { useRouter } from "next/navigation";
import { useCreateMeal } from "../../../queries";
import BackButton from "../../../components/BackButton";

export default function CreateMealPage() {
  const { mutateAsync } = useCreateMeal();
  const { push } = useRouter();

  return (
    <div>
      <div className="mb-4">
        <BackButton destination="/" />
      </div>

      <form
        className="flex flex-col items-start"
        onSubmit={async (e) => {
          e.preventDefault();

          const formData = new FormData(e.target as HTMLFormElement);
          const meal = await mutateAsync({
            name: formData.get("name") as string,
          });

          await push(`/meals/${meal.id}`);
        }}
      >
        <label className="mb-3 flex flex-col">
          <span>Name</span>
          <input
            type="text"
            name="name"
            required
            className="rounded-md border py-1 px-2 leading-none"
          />
        </label>

        <button type="submit" className="button">
          Create
        </button>
      </form>
    </div>
  );
}
