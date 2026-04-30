import type { ButtonHTMLAttributes } from "react";

type ButtonProps = ButtonHTMLAttributes<HTMLButtonElement>;

export function Button({ className = "", ...props }: ButtonProps) {
  return (
    <button
      className={`inline-flex items-center justify-center rounded-full bg-rose-400 px-4 py-2 text-sm font-semibold text-slate-950 transition hover:bg-rose-300 ${className}`}
      {...props}
    />
  );
}
