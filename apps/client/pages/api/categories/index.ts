import { NextApiRequest, NextApiResponse } from "next";

export default async function handler(
  req: NextApiRequest,
  res: NextApiResponse,
) {
  const response = await fetch(`${process.env.API_BASE_URL}/categories`);
  res.status(200).json(await response.json());
}
