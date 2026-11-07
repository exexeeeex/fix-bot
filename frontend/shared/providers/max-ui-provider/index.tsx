import { useEffect, useState, type ReactNode } from 'react';

type MaxUIComponent = typeof import('@maxhub/max-ui').MaxUI;

export function MaxUIProvider({ children }: { children: ReactNode }) {
  const [MaxUI, setMaxUI] = useState<MaxUIComponent | null>(null);

  useEffect(() => {
    if (!document.getElementById('max-web-app-sdk')) {
      const script = document.createElement('script');
      script.id = 'max-web-app-sdk';
      script.src = 'https://st.max.ru/js/max-web-app.js';
      script.async = true;
      script.addEventListener('error', () => console.error('Failed to load the MAX WebApp SDK'), {
        once: true,
      });
      document.head.appendChild(script);
    }

    import('@maxhub/max-ui').then(({ MaxUI }) => {
      setMaxUI(() => MaxUI);
    });
  }, []);

  if (!MaxUI) {
    return null;
  }

  return (
    <MaxUI className="!bg-background" platform="ios">
      {children}
    </MaxUI>
  );
}
