import { AppProps } from "next/app";
import Head from "next/head";
import { QueryClient, QueryClientProvider } from "react-query";
import "../styles.css";

const queryClient = new QueryClient();

export default function App({ Component, pageProps }: AppProps) {
  return (
    <>
      <Head>
        <title>Meal planner</title>
      </Head>
      <QueryClientProvider client={queryClient}>
        <main className="p-5 flex justify-center">
          <div className="w-full md:w-1/3">
          <Component {...pageProps} />
          </div>
        </main>
      </QueryClientProvider>
    </>
  );
}
