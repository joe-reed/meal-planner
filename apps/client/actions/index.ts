"use server";

const headers = {
  "Content-Type": "application/json",
};

export async function addIngredientToMeal(mealId: string, body: string) {
  const response = await fetch(
    `${process.env.API_BASE_URL}/meals/${mealId}/ingredients`,
    { method: "POST", headers, body },
  );
  return response.json();
}

export async function addItemToBasket(shopId: string, body: string) {
  const response = await fetch(
    `${process.env.API_BASE_URL}/baskets/${shopId}/items`,
    { method: "POST", headers, body },
  );
  return response.json();
}

export async function addMealToCurrentShop(body: string) {
  const response = await fetch(
    `${process.env.API_BASE_URL}/shops/current/meals`,
    { method: "POST", headers, body },
  );
  return response.json();
}

export async function addItemToCurrentShop(body: string) {
  const response = await fetch(
    `${process.env.API_BASE_URL}/shops/current/items`,
    { method: "POST", headers, body },
  );
  return response.json();
}

export async function removeItemFromCurrentShop(productId: string) {
  const response = await fetch(
    `${process.env.API_BASE_URL}/shops/current/items/${productId}`,
    { method: "DELETE", headers },
  );
  return response.json();
}

export async function fetchBasket(shopId: string | undefined) {
  const response = await fetch(`${process.env.API_BASE_URL}/baskets/${shopId}`);
  if (!response.ok) {
    throw new Error("Error fetching basket");
  }
  return response.json();
}

export async function fetchCategories() {
  const response = await fetch(`${process.env.API_BASE_URL}/categories`);
  if (!response.ok) {
    throw new Error("Error fetching categories");
  }
  return response.json();
}

export async function createProduct(body: string) {
  const response = await fetch(`${process.env.API_BASE_URL}/products`, {
    method: "POST",
    headers,
    body,
  });
  if (!response.ok) {
    const message = await response.json();

    return { error: message.error, product: null };
  }
  return response.json();
}

export async function createMeal(body: string) {
  const response = await fetch(`${process.env.API_BASE_URL}/meals`, {
    method: "POST",
    headers,
    body,
  });

  if (!response.ok) {
    return { error: await response.text(), meal: null };
  }

  return { error: null, meal: await response.json() };
}

export async function updateMeal(mealId: string, body: string) {
  const response = await fetch(`${process.env.API_BASE_URL}/meals/${mealId}`, {
    method: "PATCH",
    headers,
    body,
  });

  if (!response.ok) {
    return { error: await response.text(), meal: null };
  }

  return { error: null, meal: await response.json() };
}

export async function fetchCurrentShop() {
  const response = await fetch(`${process.env.API_BASE_URL}/shops/current`);
  if (!response.ok) {
    throw new Error("Error fetching current shop");
  }
  return response.json();
}

export async function fetchProducts() {
  const response = await fetch(`${process.env.API_BASE_URL}/products`);
  if (!response.ok) {
    throw new Error("Error fetching products");
  }
  return response.json();
}

export async function fetchMeal(mealId: string) {
  const response = await fetch(`${process.env.API_BASE_URL}/meals/${mealId}`);
  if (!response.ok) {
    throw new Error("Error fetching meal");
  }
  return response.json();
}

export async function fetchMeals() {
  const response = await fetch(`${process.env.API_BASE_URL}/meals`);
  if (!response.ok) {
    throw new Error("Error fetching meals");
  }
  return response.json();
}

export async function fetchShoppingList() {
  const response = await fetch(`${process.env.API_BASE_URL}/shopping-list`);
  if (!response.ok) {
    throw new Error("Error fetching shopping-list");
  }
  return response.json();
}

export async function removeIngredientFromMeal(
  mealId: string,
  ingredientId: string,
) {
  const response = await fetch(
    `${process.env.API_BASE_URL}/meals/${mealId}/ingredients/${ingredientId}`,
    { method: "DELETE", headers },
  );
  return response.json();
}

export async function removeItemFromBasket(
  shopId: string,
  ingredientId: string,
) {
  const response = await fetch(
    `${process.env.API_BASE_URL}/baskets/${shopId}/items/${ingredientId}`,
    { method: "DELETE", headers },
  );
  return response.json();
}

export async function removeMealFromCurrentShop(mealId: string) {
  const response = await fetch(
    `${process.env.API_BASE_URL}/shops/current/meals/${mealId}`,
    { method: "DELETE", headers },
  );
  return response.json();
}

export async function startShop() {
  const response = await fetch(`${process.env.API_BASE_URL}/shops`, {
    method: "POST",
    headers,
  });
  return response.json();
}
