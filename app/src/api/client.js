import { getAccessToken, reAuthenticate } from "./auth";

export default async function apiFetch(url, options = {}, requireAuth = true) {
  let token = getAccessToken();

  let response = await fetch(url, {
    ...options,
    credentials: "include",
    headers: {
      ...options.headers,
      ...(requireAuth && token ? { Authorization: `Bearer ${token}` } : {}),
    },
  });

  if (response.status === 401 && requireAuth) {
    try {
      token = await reAuthenticate();
      response = await fetch(url, {
        ...options,
        credentials: "include",
        headers: {
          ...options.headers,
          ...(token ? { Authorization: `Bearer ${token}` } : {}),
        },
      });
    } catch {
      // reAuthenticate failed, return the original 401
    }
  }

  return response;
}
