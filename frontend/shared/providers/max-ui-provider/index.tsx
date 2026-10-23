import { useEffect, useState, type ReactNode } from "react";

type MaxUIComponent = typeof import("@maxhub/max-ui").MaxUI;

export function MaxUIProvider({ children }: { children: ReactNode }) {
  const [MaxUI, setMaxUI] = useState<MaxUIComponent | null>(null);

  useEffect(() => {
    import("@maxhub/max-ui").then(({ MaxUI }) => {
      setMaxUI(() => MaxUI);
    });
  }, []);

  if (!MaxUI) {
    return null;
  }

  return <MaxUI className="bg-background" platform='ios'>{children}</MaxUI>;
}