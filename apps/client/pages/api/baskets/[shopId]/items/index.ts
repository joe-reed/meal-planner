import { NextApiRequest, NextApiResponse } from "next";

export default async function handler(
  req: NextApiRequest,
  res: NextApiResponse,
) {
  const response = await fetch(
    `${process.env.API_BASE_URL}/baskets/${req.query.shopId}/items`,
    {
      method: "POST",
      headers: {
        "Content-Type": "application/json",
      },
      body: req.body,
    },
  );
  res.status(response.status).json(await response.json());
}
