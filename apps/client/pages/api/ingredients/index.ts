import { NextApiRequest, NextApiResponse } from "next";

export default async function handler(
  req: NextApiRequest,
  res: NextApiResponse,
) {
  if (req.method === "POST") {
    const response = await fetch(`${process.env.API_BASE_URL}/ingredients`, {
      method: "POST",
      headers: {
        "Content-Type": "application/json",
      },
      body: req.body,
    });
    res.status(201).json(await response.json());
  } else {
    const response = await fetch(`${process.env.API_BASE_URL}/ingredients`);
    res.status(200).json(await response.json());
  }
}
