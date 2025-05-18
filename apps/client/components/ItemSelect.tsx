import { Product } from "../types";
import React, { useRef, useState } from "react";
import { z } from "zod";
import { SearchableSelect } from "./SearchableSelect";
import { useCreateProduct } from "../queries";
import { useCategories } from "../queries/useCategories";
import { Modal } from "./Modal";
import { Select } from "@headlessui/react";

type PendingItem = {
  productId: string;
  quantity: { amount: string; unit: string };
};

export function ItemSelect({
  onItemAdd,
  products,
  productIdsToExclude,
  className,
}: {
  onItemAdd: (body: {
    productId: string;
    quantity: { amount: number; unit: string };
  }) => void;
  products: Product[];
  productIdsToExclude: string[];
  className?: string;
}) {
  const [pendingItem, setPendingItem] = useState<PendingItem | null>(null);

  const [productSearchQuery, setProductSearchQuery] = useState("");

  const [isAddProductModalOpen, setIsAddProductModalOpen] = useState(false);

  function selectProduct(product: Product) {
    setPendingItem({
      productId: product.id,
      quantity: { amount: "", unit: "Number" },
    });
    setProductSearchQuery("");
    setTimeout(() => {
      numberInputRef.current?.focus();
    }, 10);
  }

  function addItem(pendingItem: PendingItem) {
    onItemAdd(
      z
        .object({
          productId: z.string(),
          quantity: z.object({
            amount: z.coerce.number().positive(),
            unit: z.string(),
          }),
        })
        .parse(pendingItem),
    );
    setPendingItem(null);

    ingredientSearchInputRef.current?.focus();
    ingredientSearchInputRef.current?.select();
  }

  const numberInputRef = useRef<HTMLInputElement>(null);

  const ingredientSearchInputRef = useRef<HTMLInputElement>(null);

  return (
    <div className={className}>
      {pendingItem && (
        <div className="mb-10 flex items-center justify-between space-x-3">
          <div className="whitespace-nowrap">
            {
              products.find((product) => product.id === pendingItem.productId)
                ?.name
            }
          </div>
          <div className="flex space-x-1">
            <input
              ref={numberInputRef}
              autoFocus
              type="number"
              value={pendingItem.quantity.amount}
              className="button bg-white px-2 py-1"
              size={2}
              onChange={(e) =>
                setPendingItem({
                  ...pendingItem,
                  quantity: {
                    ...pendingItem.quantity,
                    amount: e.target.value,
                  },
                })
              }
            />
            <select
              onChange={(e) => {
                setPendingItem({
                  ...pendingItem,
                  quantity: {
                    ...pendingItem.quantity,
                    unit: e.target.value,
                  },
                });
              }}
              className="button bg-white px-2 py-1"
            >
              {/*todo: fetch these from api*/}
              <option value="Number">Number</option>
              <option value="Tsp">Tsp</option>
              <option value="Tbsp">Tbsp</option>
              <option value="Cup">Cup</option>
              <option value="Oz">Oz</option>
              <option value="Lb">Lb</option>
              <option value="Gram">Gram</option>
              <option value="Kg">Kg</option>
              <option value="Ml">Ml</option>
              <option value="Litre">Litre</option>
              <option value="Pinch">Pinch</option>
              <option value="Bunch">Bunch</option>
              <option value="Pack">Pack</option>
              <option value="Tin">Tin</option>
            </select>
            <button
              onClick={() => {
                if (pendingItem) {
                  addItem(pendingItem);
                }
              }}
              className="button"
            >
              Add
            </button>
          </div>
        </div>
      )}
      <div className="flex items-center">
        <SearchableSelect<Product>
          options={products.filter(
            (product) => !productIdsToExclude.some((i) => i === product.id),
          )}
          onSelect={selectProduct}
          onInputChange={(query) => setProductSearchQuery(query)}
          inputRef={ingredientSearchInputRef}
        />
        <button
          onClick={() => setIsAddProductModalOpen(true)}
          className="ml-2 whitespace-nowrap underline"
        >
          Add new ingredient
        </button>
        <AddNewProductModal
          text={productSearchQuery}
          isOpen={isAddProductModalOpen}
          setIsOpen={setIsAddProductModalOpen}
          onAdd={selectProduct}
        />
      </div>
    </div>
  );
}

function AddNewProductModal({
  text,
  isOpen,
  setIsOpen,
  onAdd,
}: {
  text: string;
  isOpen: boolean;
  setIsOpen: (value: boolean) => void;
  onAdd: (ingredient: Product) => void;
}) {
  const { mutateAsync } = useCreateProduct();

  const { data: categories } = useCategories();

  return (
    <>
      <Modal
        isOpen={isOpen}
        setIsOpen={setIsOpen}
        title="Add new product"
        body={(close) => (
          <div className="flex justify-between">
            <form
              onSubmit={async (e) => {
                e.preventDefault();

                const formData = new FormData(e.target as HTMLFormElement);
                const name = formData.get("name") as string;
                const category = formData.get("category") as string;

                const response = await mutateAsync({
                  name,
                  category,
                });

                onAdd(response);

                close();
              }}
            >
              <label className="mb-3 flex flex-col">
                <span>Name</span>
                <input
                  type="text"
                  name="name"
                  required
                  className="rounded-md border py-1 px-2 leading-none"
                  defaultValue={text}
                  data-autofocus
                />
              </label>

              <label className="mb-3 flex flex-col">
                <span>Category</span>
                <Select
                  name="category"
                  aria-label="Product category"
                  className="rounded-md border bg-white py-1 px-2 leading-none"
                >
                  <option value="">Select a category</option>
                  {categories?.map((category) => (
                    <option key={category.name} value={category.name}>
                      {category.name}
                    </option>
                  ))}
                </Select>
              </label>

              <div>
                <button type="submit" className="button mr-3">
                  Create
                </button>

                <button onClick={close} className="underline">
                  Cancel
                </button>
              </div>
            </form>
          </div>
        )}
      />
    </>
  );
}
