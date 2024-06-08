import { NextApiRequest, NextApiResponse } from "next";

export default async function handler(
  req: NextApiRequest,
  res: NextApiResponse,
) {
  if (req.method === "POST") {
    const response = await fetch("http://127.0.0.1:1323", {
      method: "POST",
      headers: {
        "Content-Type": "application/json",
      },
      body: req.body,
    });
    res.status(201).json(await response.json());
  } else {
    const response = await fetch("http://127.0.0.1:1323");
    res.status(200).json(await response.json());
  }
}
