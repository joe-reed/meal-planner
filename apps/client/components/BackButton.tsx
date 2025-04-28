import clsx from "clsx";
import Link from "next/link";

export default function BackButton({
  destination,
  className,
}: {
  destination: string;
  className?: string;
}) {
  return (
    <Link href={destination} className={clsx(className, "text-3xl")}>
      🔙
    </Link>
  );
}
