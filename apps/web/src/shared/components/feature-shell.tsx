import type { PropsWithChildren } from "react";

export function FeatureShell({ children }: PropsWithChildren) {
  return <main className="app-shell">{children}</main>;
}
