import { NextApiRequest, NextApiResponse } from "next";

export default async function handler(
  req: NextApiRequest,
  res: NextApiResponse,
) {
  const response = await fetch(
    `${process.env.API_BASE_URL}/baskets/${req.query.shopId}/items/${req.query.ingredientId}`,
    {
      method: "DELETE",
      headers: {
        "Content-Type": "application/json",
      },
    },
  );
  res.status(200).json(await response.json());
}
