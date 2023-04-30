import useCreateMeal from "../../queries/useCreateMeal";
import { useRouter } from "next/router";

export default function CreateMeal() {
  const { mutate } = useCreateMeal();
  const { push } = useRouter();

  return (
    <div>
      <form
        onSubmit={(e) => {
          e.preventDefault();

          const formData = new FormData(e.target as HTMLFormElement);
          mutate({ name: formData.get("name") as string });
          push("/");
        }}
      >
        <label>
          Name
          <input type="text" name="name" required />
        </label>

        <button type="submit">Create</button>
      </form>
    </div>
  );
}
