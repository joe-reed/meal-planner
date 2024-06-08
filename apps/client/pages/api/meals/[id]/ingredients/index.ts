import { NextApiRequest, NextApiResponse } from "next";

export default async function handler(
  req: NextApiRequest,
  res: NextApiResponse,
) {
  const response = await fetch(
    `http://127.0.0.1:1323/meals/${req.query.id}/ingredients`,
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
