export type Shop = {
  id: string;
  meals: { id: string }[];
  items: {
    productId: string;
    quantity: {
      amount: number;
      unit: string;
    };
  }[];
};
