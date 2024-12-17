import "../global.css";
import { Metadata } from "next";
import ReactQueryClientProvider from "../components/ReactQueryClientProvider";
import React from "react";

export const metadata: Metadata = {
  title: "Meal planner",
};

export default function RootLayout({
  children,
}: {
  children: React.ReactNode;
}) {
  return (
    <html lang="en">
      <body className="flex justify-center p-5">
        <ReactQueryClientProvider>
          <main className="w-full md:w-2/3 lg:w-1/2">{children}</main>
        </ReactQueryClientProvider>
      </body>
    </html>
  );
}
