export function Unit({
  quantity: { unit, amount },
}: {
  quantity: { unit: string; amount: number };
}) {
  return unit !== "Number" ? " " + unit + (amount > 1 ? "s" : "") : "x";
}
