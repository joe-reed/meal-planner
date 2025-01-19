"use client";

import { useRouter } from "next/navigation";
import { useUploadMeals } from "../../../queries";
import BackButton from "../../../components/BackButton";
import React, { useState } from "react";

export default function UploadMealsPage() {
  const { mutateAsync } = useUploadMeals();
  const { push } = useRouter();

  const [file, setFile] = useState<File | null>(null);

  function handleFileChange(e: React.ChangeEvent<HTMLInputElement>) {
    if (e.target.files) {
      setFile(e.target.files[0]);
    }
  }

  async function handleUpload() {
    if (file) {
      await mutateAsync(file);

      push(`/`);
    }
  }

  return (
    <div>
      <div className="mb-4 flex items-center">
        <BackButton className="mr-3" destination="/" />
        <h1 className="text-lg font-bold">Upload meals</h1>
      </div>

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
