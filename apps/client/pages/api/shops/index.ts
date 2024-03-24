import { NextApiRequest, NextApiResponse } from "next";

export default async function handler(
  req: NextApiRequest,
  res: NextApiResponse
) {
  if (req.method === "POST") {
    const response = await fetch("http://localhost:1323/shops", {
      method: "POST",
      headers: {
        "Content-Type": "application/json",
      },
      body: req.body,
    });
    res.status(201).json(await response.json());
  }
}