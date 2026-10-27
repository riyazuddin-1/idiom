const AUTH_BASE = "https://idiom-identity.onrender.com/api/v1";
const ACCESS_TOKEN_KEY = "idiom_access_token";
const MOCK_TOKEN = "mock-dev-token";

let accessToken = sessionStorage.getItem(ACCESS_TOKEN_KEY);

export function isMockAuthEnabled() {
  return import.meta.env.DEV && String(import.meta.env.VITE_USE_MOCK_AUTH || "").toLowerCase() === "true";
}

export function isAuthenticated() {
  return !!getAccessToken();
}

export function getAccessToken() {
  if (isMockAuthEnabled()) {
    return sessionStorage.getItem(ACCESS_TOKEN_KEY) || MOCK_TOKEN;
  }

  return accessToken;
}

export function setAccessToken(token) {
  accessToken = token;

  if (token) {
    sessionStorage.setItem(ACCESS_TOKEN_KEY, token);
  } else {
    sessionStorage.removeItem(ACCESS_TOKEN_KEY);
  }
}

export function clearAccessToken() {
  setAccessToken(null);
}

export async function authenticate(code) {
  if (isMockAuthEnabled()) {
    setAccessToken(MOCK_TOKEN);
    return MOCK_TOKEN;
  }

  const response = await fetch(`${AUTH_BASE}/token`, {
    method: "POST",
    headers: { "Content-Type": "application/json", "X-Project-ID": "idiom" },
    body: JSON.stringify({ code }),
    credentials: "include",
  });

  if (!response.ok) {
    throw new Error("Failed to authenticate");
  }

  const { data } = await response.json();
  setAccessToken(data.accessToken);
  return accessToken;
}

export async function reAuthenticate() {
  if (isMockAuthEnabled()) {
    setAccessToken(MOCK_TOKEN);
    return MOCK_TOKEN;
  }

  const response = await fetch(`${AUTH_BASE}/token/refresh`, {
    method: "POST",
    credentials: "include",
  });

  if (!response.ok) {
    clearAccessToken();
    throw new Error("Re-authentication failed");
  }

  const { data } = await response.json();
  setAccessToken(data.accessToken);
  return data.accessToken;
}
