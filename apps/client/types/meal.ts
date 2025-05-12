export type Meal = {
  id: string;
  name: string;
  url: string;
  ingredients: Ingredient[];
};

export type Ingredient = {
  id: string;
  quantity: {
    amount: number;
    unit: string;
  };
};
