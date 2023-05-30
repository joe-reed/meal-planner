import { NextApiRequest, NextApiResponse } from "next";

export default async function handler(
  req: NextApiRequest,
  res: NextApiResponse
) {
  const response = await fetch(`http://localhost:1323/meals/${req.query.id}`);
  res.status(200).json(await response.json());
}
