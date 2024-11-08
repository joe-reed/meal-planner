import { useRouter } from "next/router";
import { useCreateMeal } from "../../queries";
import BackButton from "../../components/BackButton";

export default function CreateMeal() {
  const { mutateAsync } = useCreateMeal();
  const { push } = useRouter();

  return (
    <div>
      <BackButton destination="/" />
      <form
        className="flex flex-col items-start"
        onSubmit={async (e) => {
          e.preventDefault();

          const formData = new FormData(e.target as HTMLFormElement);
          const response = await mutateAsync({
            name: formData.get("name") as string,
          });
          const meal = await response.json();

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
