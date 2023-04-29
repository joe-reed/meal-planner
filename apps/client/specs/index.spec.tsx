import React from "react";
import { render, screen } from "@testing-library/react";
import Index from "../pages/index";

  it("should render a heading", () => {
    render(<Index />);
    expect(screen.getByRole('heading')).toHaveTextContent('Meal planner')
  });
