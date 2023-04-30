import React from "react";
import { render, screen } from "@testing-library/react";
import Index from "./index";

jest.mock("../queries/useMeals", () =>
  jest.fn(() => ({
    isLoading: false,
    isError: false,
    data: [
      { id: "1", name: "foo" },
      { id: "2", name: "bar" },
    ],
  }))
);

it("renders a heading", async () => {
  render(<Index />);

  expect(screen.getByRole("heading")).toHaveTextContent("Meal planner");
});

it("renders meals", async () => {
  render(<Index />);

  expect(screen.getByText("foo")).toBeInTheDocument();
  expect(screen.getByText("bar")).toBeInTheDocument();
});
