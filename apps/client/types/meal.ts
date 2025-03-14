export type Meal = {
  id: string;
  name: string;
  ingredients: MealIngredient[];
};

export type MealIngredient = {
  id: string;
  quantity: {
    amount: number;
    unit: string;
  };
};
