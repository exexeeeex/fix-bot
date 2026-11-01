const isDev = import.meta.env.DEV;

export const getMax = (): MaxWebApp => {
  if (!isDev && window.WebApp) {
    return window.WebApp;
  }

  return {
    initData: import.meta.env.VITE_MOCK_INIT_DATA ?? '',
    platform: 'web',
    version: 'dev',
  };
};
