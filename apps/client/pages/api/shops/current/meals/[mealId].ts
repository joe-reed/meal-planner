import { NextApiRequest, NextApiResponse } from "next";

export default async function handler(
  req: NextApiRequest,
  res: NextApiResponse
) {
  const response = await fetch(
    `http://localhost:1323/shops/current/meals/${req.query.mealId}`,
    {
      method: "DELETE",
      headers: {
        "Content-Type": "application/json",
      },
    }
  );
  res.status(200).json(await response.json());
}
