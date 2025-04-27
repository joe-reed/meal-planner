"use client";

import { Ingredient } from "../../types";
import React from "react";
import BackButton from "../../components/BackButton";
import clsx from "clsx";
import { useAddItemToBasket } from "../../queries/useAddItemToBasket";
import { useRemoveItemFromBasket } from "../../queries/useRemoveItemFromBasket";
import { useShoppingList } from "../../queries/useShoppingList";
import { Popover, PopoverButton, PopoverPanel } from "@headlessui/react";
import { Unit } from "../../components/Unit";

export default function ShopPage() {
  const shoppingListQuery = useShoppingList();

  const [showItemsInBasket, setShowItemsInBasket] = React.useState(false);

  const [canUndo, setCanUndo] = React.useState(false);
  const [undoTimeout, setUndoTimeout] = React.useState<NodeJS.Timeout | null>(
    null,
  );
  // todo: keep stack of actions
  const [lastAction, setLastAction] = React.useState<{
    action: string;
    ingredientId: string;
  } | null>(null);

  function toggleShowItemsInBasket() {
    setShowItemsInBasket(!showItemsInBasket);
  }

  const { shopId, shoppingList: shoppingListData } = shoppingListQuery.data || {
    shopId: "0",
    shoppingList: {},
  };

  function useUndo() {
    const { mutate: addItemToBasket } = useAddItemToBasket(shopId);
    const { mutate: removeItemFromBasket } = useRemoveItemFromBasket(shopId);

    return function undo() {
      if (!lastAction) return;

      if (lastAction.action === "add") {
        removeItemFromBasket(lastAction.ingredientId);
      } else if (lastAction.action === "remove") {
        addItemToBasket({ ingredientId: lastAction.ingredientId });
      }
      setCanUndo(false);
      setLastAction(null);
    };
  }

  const undo = useUndo();

  function onAddToBasket(ingredient: Ingredient) {
    setLastAction({
      action: "add",
      ingredientId: ingredient.id,
    });
    setCanUndo(true);
    if (undoTimeout) {
      clearTimeout(undoTimeout);
    }
    const timeout = setTimeout(() => {
      setCanUndo(false);
      setLastAction(null);
    }, 5000);

    setUndoTimeout(timeout);
  }

  function onRemoveFromBasket(ingredient: Ingredient) {
    setLastAction({
      action: "remove",
      ingredientId: ingredient.id,
    });
    setCanUndo(true);
    if (undoTimeout) {
      clearTimeout(undoTimeout);
    }
    const timeout = setTimeout(() => {
      setCanUndo(false);
      setLastAction(null);
    }, 5000);

    setUndoTimeout(timeout);
  }

  if (shoppingListQuery.isInitialLoading) {
    return <p>Loading...</p>;
  }

  if (shoppingListQuery.error) {
    return <p>Error: {shoppingListQuery.error.message}</p>;
  }

  const shoppingList = Object.values(shoppingListData);

  const filteredIngredients = shoppingList.filter(
    (ingredient) => showItemsInBasket || !ingredient.isInBasket,
  );

  const categorisedIngredients = Object.groupBy<
    string,
    Ingredient & {
      mealCount: number;
      isInBasket: boolean;
      quantities: { unit: string; amount: number }[];
    }
  >(filteredIngredients, ({ category }) => category);

  return (
    <>
      <ShowIngredientsButton
        showItemsInBasket={showItemsInBasket}
        toggleShowItemsInBasket={toggleShowItemsInBasket}
        className="fixed bottom-4 right-4 z-10 shadow-md sm:hidden"
      />
      <div className="flex w-full flex-col">
        <div className="mb-4 flex items-center justify-between">
          <div className="flex items-center">
            <BackButton className="mr-3" destination="/" />
            <h1 className="text-lg font-bold">Current shop</h1>
          </div>
          <div className="flex items-center space-x-2">
            {canUndo && (
              <button className="button" onClick={undo}>
                Undo
              </button>
            )}
            <ShowIngredientsButton
              showItemsInBasket={showItemsInBasket}
              toggleShowItemsInBasket={toggleShowItemsInBasket}
              className="hidden sm:block"
            />
          </div>
        </div>

        {filteredIngredients.length === 0 ? (
          shoppingList.length === 0 ? (
            <p className="text-center">
              No ingredients in this shop yet. Go back and add some meals!
            </p>
          ) : (
            <p className="text-center">
              All ingredients are in basket. Use the button to show all
              ingredients.
            </p>
          )
        ) : null}

        {Object.keys(categorisedIngredients)
          .sort()
          .map((category) => (
            <div className="mb-4" key={category}>
              <h2 className="mb-2 text-xl font-bold">{category}</h2>
              <ul>
                {(categorisedIngredients[category] ?? []).map((ingredient) => (
                  <IngredientListItem
                    key={ingredient.id}
                    ingredient={ingredient}
                    shopId={shopId}
                    onAddToBasket={onAddToBasket}
                    onRemoveFromBasket={onRemoveFromBasket}
                  />
                ))}
              </ul>
            </div>
          ))}
      </div>
    </>
  );
}

function IngredientListItem({
  ingredient,
  shopId,
  onAddToBasket,
  onRemoveFromBasket,
}: {
  ingredient: Ingredient & {
    mealCount: number;
    isInBasket: boolean;
    quantities: { unit: string; amount: number }[];
  };
  shopId: string;
  onAddToBasket: (ingredient: Ingredient) => void;
  onRemoveFromBasket: (ingredient: Ingredient) => void;
}) {
  const { mutate: addItemToBasket } = useAddItemToBasket(shopId);
  const { mutate: removeItemFromBasket } = useRemoveItemFromBasket(shopId);

  return (
    <li className="mb-3 flex items-center justify-between leading-4">
      <label
        className={clsx("flex w-full justify-between break-words pr-6", {
          "line-through": ingredient.isInBasket,
        })}
      >
        {ingredient.name}

        <input
          type="checkbox"
          checked={ingredient.isInBasket}
          onChange={() => {
            if (ingredient.isInBasket) {
              removeItemFromBasket(ingredient.id);
              onRemoveFromBasket(ingredient);
            } else {
              addItemToBasket({ ingredientId: ingredient.id });
              onAddToBasket(ingredient);
            }
          }}
        />
      </label>
      <Popover className="relative">
        <PopoverButton className="flex-no-shrink whitespace-nowrap">
          {ingredient.mealCount} <span className="text-xs">meals</span>
        </PopoverButton>
        <PopoverPanel
          anchor="bottom"
          className="flex flex-col space-y-2 rounded-lg border bg-white px-3 py-2 text-xs shadow-lg"
        >
          {ingredient.quantities.map(({ unit, amount }, index) => (
            <span key={index}>
              {amount}
              <Unit quantity={{ unit, amount }} />
            </span>
          ))}
        </PopoverPanel>
      </Popover>
    </li>
  );
}
function ShowIngredientsButton({
  showItemsInBasket,
  toggleShowItemsInBasket,
  className,
}: {
  showItemsInBasket: boolean;
  toggleShowItemsInBasket: () => void;
  className?: string;
}) {
  return (
    <button
      onClick={toggleShowItemsInBasket}
      className={clsx("button", className)}
    >
      {showItemsInBasket
        ? "Hide ingredients in basket"
        : "Show all ingredients"}
    </button>
  );
}
