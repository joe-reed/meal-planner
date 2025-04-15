"use client";

import { useRouter } from "next/navigation";
import { useCreateIngredient, useUploadMeals } from "../../../queries";
import BackButton from "../../../components/BackButton";
import React, { useState } from "react";
import { Select } from "@headlessui/react";
import { useCategories } from "../../../queries/useCategories";
import clsx from "clsx";

export default function UploadMealsPage() {
  const { push } = useRouter();
  const { mutateAsync } = useUploadMeals();
  const [notFoundIngredients, setNotFoundIngredients] = useState<string[]>([]);

  const [file, setFile] = useState<File | null>(null);

  function handleFileChange(e: React.ChangeEvent<HTMLInputElement>) {
    if (e.target.files) {
      setFile(e.target.files[0]);
    }
  }

  async function handleUpload() {
    if (file) {
      const { status, error } = await mutateAsync(file);

      if (!error) {
        return push("/");
      }

      if (status === 400) {
        setNotFoundIngredients(error.data.notFoundIngredients);
      }
    }
  }

  return (
    <div>
      <div className="mb-4 flex items-center">
        <BackButton className="mr-3" destination="/" />
        <h1 className="text-lg font-bold">Upload meals</h1>
      </div>

      {notFoundIngredients.length > 0 && (
        <div className="mb-4 bg-red-100 p-4">
          <h2 className="font-bold text-red-800">Ingredients not found</h2>
          <p className="mb-2">
            The following ingredients were not found. Create them below, or fix
            any issues in your file.
          </p>
          <ul>
            {notFoundIngredients.map((ingredient) => (
              <AddIngredientForm
                key={ingredient}
                name={ingredient}
                className="mb-1"
              />
            ))}
          </ul>
        </div>
      )}

      <input
        type="file"
        onChange={handleFileChange}
        className="file-button file:mr-2"
      />

      {file && (
        <button onClick={handleUpload} className="button">
          Upload
        </button>
      )}
    </div>
  );
}

function AddIngredientForm({
  name,
  className,
}: {
  name: string;
  className: string;
}) {
  const { mutateAsync } = useCreateIngredient();
  const { data: categories } = useCategories();

  const [isSubmitted, setIsSubmitted] = useState(false);
  const [hasError, setHasError] = useState(false);

  return (
    <form
      className={clsx("flex items-center", className)}
      onSubmit={async (e) => {
        e.preventDefault();
        setHasError(false);

        const formData = new FormData(e.target as HTMLFormElement);
        const name = formData.get("name") as string;
        const category = formData.get("category") as string;

        const { error } = await mutateAsync({
          name,
          category,
        });

        if (error) {
          setHasError(true);
          return;
        }

        setIsSubmitted(true);
      }}
    >
      <span className="mr-auto">{name}</span>
      <input type="hidden" name="name" value={name} />

      {isSubmitted ? (
        <span className="text-green-500">Created</span>
      ) : (
        <>
          <Select
            name="category"
            aria-label="Ingredient category"
            className={clsx(
              "mr-3 rounded-md border bg-white py-1 px-2 leading-none",
              {
                "border-red-500 text-red-900": hasError,
              },
            )}
            onChange={() => setHasError(false)}
          >
            <option value="">Select a category</option>
            {categories?.map((category) => (
              <option key={category.name} value={category.name}>
                {category.name}
              </option>
            ))}
          </Select>

          <button type="submit" className="button">
            Create
          </button>
        </>
      )}
    </form>
  );
}
