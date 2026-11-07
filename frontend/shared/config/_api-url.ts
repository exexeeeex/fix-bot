export const env = import.meta.env.MODE;

export const BASE_API_URL =
  import.meta.env.VITE_API_BASE_URL ??
  (env === 'development' ? 'http://localhost:8080/api/v1/' : '/api/');
