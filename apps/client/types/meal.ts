export type Meal = {
  id: string;
  name: string;
  url: string;
  ingredients: MealIngredient[];
};

export type MealIngredient = {
  id: string;
  quantity: {
    amount: number;
    unit: string;
  };
};
