/* eslint-disable @typescript-eslint/no-non-null-assertion */

import React, { ReactElement } from "react";
import {
  render as baseRender,
  RenderOptions,
  screen,
  within,
} from "@testing-library/react";
import Index from "./index";
import userEvent from "@testing-library/user-event";
import { QueryClient, QueryClientProvider } from "@tanstack/react-query";

function wrapper({ children }: { children: React.ReactNode }) {
  return (
    <QueryClientProvider client={new QueryClient()}>
      {children}
    </QueryClientProvider>
  );
}

const render = (ui: ReactElement, options?: Omit<RenderOptions, "wrapper">) =>
  baseRender(ui, { wrapper, ...options });

const mockUseStartShopMutate = jest.fn();
jest.mock("../queries", () => ({
  useMeals: () => ({
    isInitialLoading: false,
    isError: false,
    data: [
      { id: "1", name: "foo" },
      { id: "2", name: "bar" },
      { id: "3", name: "baz" },
    ],
  }),
  useCurrentShop: () => ({
    isInitialLoading: false,
    isError: false,
    data: { id: 5, meals: [{ id: 3 }, { id: 2 }] },
  }),
  useStartShop: () => ({
    mutate: mockUseStartShopMutate,
  }),
  useAddMealToCurrentShop: () => ({
    mutate: jest.fn(),
  }),
  useRemoveMealFromCurrentShop: () => ({
    mutate: jest.fn(),
  }),
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

  expect(mockUseStartShopMutate).toHaveBeenCalled();
});
