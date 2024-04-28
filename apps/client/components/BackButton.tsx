import { useRouter } from "next/router";
import clsx from "clsx";

export default function BackButton({ className }: { className?: string }) {
  const { back } = useRouter();

  return (
    <button onClick={() => back()} className={clsx(className, "text-3xl")}>
      🔙
    </button>
  );
}
