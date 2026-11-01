interface MaxWebApp {
  initData: string;
  platform: string;
  version: string;
}

interface Window {
  WebApp?: MaxWebApp;
}
