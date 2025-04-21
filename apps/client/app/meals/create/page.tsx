"use client";

import { useRouter } from "next/navigation";
import { useCreateMeal } from "../../../queries";
import BackButton from "../../../components/BackButton";
import { useState } from "react";

export default function CreateMealPage() {
  const { mutateAsync } = useCreateMeal();
  const { push } = useRouter();
  const [errorMessage, setErrorMessage] = useState<string | null>(null);

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

          const { error, meal } = await mutateAsync({
            name: formData.get("name") as string,
            url: formData.get("url") as string,
          });

          if (error) {
            setErrorMessage(error);
            return;
          }

          await push(`/meals/${meal.id}`);
        }}
      >
        <label className="mb-3 flex w-2/3 flex-col">
          <span>Name</span>
          <input
            type="text"
            name="name"
            required
            className="rounded-md border py-1 px-2 leading-none"
          />
        </label>

        <label className="mb-3 flex w-2/3 flex-col">
          <span>URL</span>
          <input
            type="text"
            name="url"
            required
            className="rounded-md border py-1 px-2 leading-none"
          />
        </label>

        <button type="submit" className="button">
          Create
        </button>

        {errorMessage ? <p className="text-red-500">{errorMessage}</p> : null}
      </form>
    </div>
  );
}
