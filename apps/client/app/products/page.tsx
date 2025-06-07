"use client";

import BackButton from "../../components/BackButton";
import React from "react";
import { useGroupedProducts } from "../../queries/useGroupedProducts";

export default function ProductsPage() {
  const groupedProductsQuery = useGroupedProducts();

  if (groupedProductsQuery.isInitialLoading) {
    return <p>Loading...</p>;
  }

  if (groupedProductsQuery.isError) {
    return <p>Error: {groupedProductsQuery.error.message}</p>;
  }
  const groupedProducts = groupedProductsQuery.data;

  return (
    <div className="flex flex-col">
      <div className="mb-2 flex items-center">
        <BackButton className="mr-3" destination="/" />
        <h1 className="text-lg font-bold">Products</h1>
      </div>

      <div>
        {groupedProducts ? (
          <div className="mt-4">
            {Object.entries(groupedProducts).map(([category, products]) => (
              <div key={category} className="mb-6">
                <h2 className="text-lg font-semibold">{category}</h2>
                <ul>
                  {products.map((product) => (
                    <li key={product.id} className="mb-1">
                      {product.name}
                    </li>
                  ))}
                </ul>
              </div>
            ))}
          </div>
        ) : (
          <p>No products available.</p>
        )}
      </div>
    </div>
  );
}
