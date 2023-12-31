import React from "react";
import { render, screen, within } from "@testing-library/react";
import Index from "./index";

jest.mock("../queries/useMeals", () => () => ({
  isLoading: false,
  isError: false,
  data: [
    { id: "1", name: "foo" },
    { id: "2", name: "bar" },
    { id: "3", name: "baz" },
  ],
}));

jest.mock("../queries/useWeek", () => () => ({
  isLoading: false,
  isError: false,
  data: [
    { id: "1", name: "foo" },
    { id: "2", name: "bar" },
  ],
}));

it("renders meals", async () => {
  render(<Index />);

  const meals = screen.getByText("Meals").parentElement!;

  expect(meals).not.toBeNull();
  expect(within(meals).getByText("foo")).toBeInTheDocument();
  expect(within(meals).getByText("bar")).toBeInTheDocument();
  expect(within(meals).getByText("baz")).toBeInTheDocument();
});

it("renders selected meals", async () => {
  render(<Index />);

  const week = screen.getByText("This week").parentElement!;
  expect(week).not.toBeNull();
  expect(within(week).getByText("foo")).toBeInTheDocument();
  expect(within(week).getByText("bar")).toBeInTheDocument();
});
