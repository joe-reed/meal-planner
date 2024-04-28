import { useRouter } from "next/router";
import { useCreateMeal } from "../../queries";
import BackButton from "../../components/BackButton";

export default function CreateMeal() {
  const { mutate } = useCreateMeal();
  const { push } = useRouter();

  return (
    <div>
      <BackButton />
      <form
        onSubmit={async (e) => {
          e.preventDefault();

          const formData = new FormData(e.target as HTMLFormElement);
          mutate({ name: formData.get("name") as string });
          await push("/");
        }}
      >
        <label className="mr-2">
          <span className="mr-2">Name</span>
          <input
            type="text"
            name="name"
            required
            className="rounded-md border py-1 px-2 leading-none"
            autoFocus
          />
        </label>

        <button type="submit" className="button">
          Create
        </button>
      </form>
    </div>
  );
}
