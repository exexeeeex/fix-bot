export interface SessionTokens {
  access_token: string;
  refresh_token: string;
}

export function isSessionTokens(value: unknown): value is SessionTokens {
  if (typeof value !== 'object' || value === null) return false;

  return (
    'access_token' in value &&
    typeof value.access_token === 'string' &&
    value.access_token.length > 0 &&
    'refresh_token' in value &&
    typeof value.refresh_token === 'string' &&
    value.refresh_token.length > 0
  );
}

export function clearSessionStorage() {
  if (typeof window === 'undefined') return;

  window.sessionStorage.removeItem('_tokens');
  window.sessionStorage.removeItem('_user-role');
  window.sessionStorage.removeItem('_organization-id');
}

export function getSessionTokens(): SessionTokens | null {
  if (typeof window === 'undefined') return null;

  const storedTokens = window.sessionStorage.getItem('_tokens');
  if (!storedTokens) return null;

  try {
    const tokens: unknown = JSON.parse(storedTokens);
    if (isSessionTokens(tokens)) return tokens;
  } catch {
    clearSessionStorage();
    return null;
  }

  clearSessionStorage();
  return null;
}

export function storeSessionTokens(tokens: SessionTokens) {
  if (typeof window === 'undefined') return;

  window.sessionStorage.setItem('_tokens', JSON.stringify(tokens));
}
