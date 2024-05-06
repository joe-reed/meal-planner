import { useRouter } from "next/router";
import clsx from "clsx";

export default function BackButton({
  className,
  destination,
}: {
  className?: string;
  destination?: string;
}) {
  const { back, push } = useRouter();

  const onClick = destination ? () => push(destination) : () => back();

  return (
    <button onClick={onClick} className={clsx(className, "text-3xl")}>
      🔙
    </button>
  );
}
