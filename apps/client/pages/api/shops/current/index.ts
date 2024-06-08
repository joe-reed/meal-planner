import { NextApiRequest, NextApiResponse } from "next";

export default async function handler(
  req: NextApiRequest,
  res: NextApiResponse,
) {
  const response = await fetch("http://127.0.0.1:1323/shops/current");
  res.status(200).json(await response.json());
}
