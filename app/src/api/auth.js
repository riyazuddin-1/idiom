let accessToken = null;

const AUTH_BASE = "/api/v1";

export function getAccessToken() {
  return accessToken;
}

export function setAccessToken(token) {
  accessToken = token;
}

export function clearAccessToken() {
  accessToken = null;
}

export async function authenticate(code) {
  const response = await fetch(`${AUTH_BASE}/token`, {
    method: "POST",
    headers: { "Content-Type": "application/json" },
    body: JSON.stringify({ code, project_id: "idiom" }),
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
