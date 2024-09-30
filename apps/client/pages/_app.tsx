import { AppProps } from "next/app";
import Head from "next/head";
import { QueryClient, QueryClientProvider } from "@tanstack/react-query";
import "../styles.css";

const queryClient = new QueryClient();

export default function App({ Component, pageProps }: AppProps) {
  return (
    <>
      <Head>
        <title>Meal planner</title>
      </Head>
      <QueryClientProvider client={queryClient}>
        <main className="flex justify-center p-5">
          <div className="w-full md:w-2/3 lg:w-1/2">
            <Component {...pageProps} />
          </div>
        </main>
      </QueryClientProvider>
    </>
  );
}
