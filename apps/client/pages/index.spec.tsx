/* eslint-disable @typescript-eslint/no-non-null-assertion */

import React from "react";
import { render, screen, within } from "@testing-library/react";
import Index from "./index";
import userEvent from "@testing-library/user-event";

jest.mock("../queries/useMeals", () => () => ({
  isLoading: false,
  isError: false,
  data: [
    { id: "1", name: "foo" },
    { id: "2", name: "bar" },
    { id: "3", name: "baz" },
  ],
}));

jest.mock("../queries/useCurrentShop", () => () => ({
  isLoading: false,
  isError: false,
  data: { id: 5, meals: [{ id: 3 }, { id: 2 }] },
}));

const mockMutate = jest.fn();
jest.mock("../queries/useStartShop", () => () => ({
  mutate: mockMutate,
}));

it("renders meals", async () => {
  render(<Index />);

  const meals = screen.getByText("Meals").parentElement!;

  expect(meals).not.toBeNull();
  expect(within(meals).getByText("foo")).toBeInTheDocument();
  expect(within(meals).getByText("bar")).toBeInTheDocument();
  expect(within(meals).getByText("baz")).toBeInTheDocument();
});

it("renders current shop", async () => {
  render(<Index />);

  const shop = screen.getByText("Current shop").parentElement!;
  expect(shop).not.toBeNull();
  expect(within(shop).getByText("Shop number 5")).toBeInTheDocument();
  expect(within(shop).getByText("baz")).toBeInTheDocument();
  expect(within(shop).getByText("bar")).toBeInTheDocument();
});

it("starts a shop", async () => {
  render(<Index />);

  await userEvent.click(screen.getByText("Start Shop"));

  expect(mockMutate).toHaveBeenCalled();
});
