import { render, screen } from "@testing-library/react";
import userEvent from "@testing-library/user-event";
import CreateMeal from "./create";

const mockMutate = jest.fn();
jest.mock("../../queries/useCreateMeal", () => () => ({
  mutate: mockMutate,
}));

const mockPush = jest.fn();
jest.mock("next/router", () => ({
  useRouter: () => ({
    push: mockPush,
  }),
}));

it("creates a meal", async () => {
  render(<CreateMeal />);

  await submitForm("my new meal");

  expect(mockMutate).toHaveBeenCalledWith({ name: "my new meal" });
});

it("redirects home", async () => {
  render(<CreateMeal />);

  submitForm();

  expect(mockPush).toHaveBeenCalledWith("/");
});

async function submitForm(name = "my new meal") {
  await userEvent.type(screen.getByLabelText(/name/i), name);
  await userEvent.click(screen.getByRole("button", { name: /create/i }));
}
