import { NextResponse } from "next/server";

export async function POST(req: Request) {
  const formData = await req.formData();

  const response = await fetch(`${process.env.API_BASE_URL}/meals/upload`, {
    method: "POST",
    body: formData,
  });

  return new NextResponse(response.body, {
    status: response.status,
    headers: response.headers,
  });
}
