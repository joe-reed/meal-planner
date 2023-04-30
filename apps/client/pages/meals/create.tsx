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
        <label className="mr-2">
          <span className="mr-2">Name</span>
          <input
            type="text"
            name="name"
            required
            className="border rounded-md py-1 leading-none px-2"
          />
        </label>

        <button type="submit" className="button">
          Create
        </button>
      </form>
    </div>
  );
}
