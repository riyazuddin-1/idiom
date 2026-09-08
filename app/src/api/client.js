import { getAccessToken, isMockAuthEnabled, reAuthenticate } from "./auth";
import {
  addMockOrganization,
  addMockProject,
  getMockOrganizations,
  getMockProjects,
} from "./mockData";

function mockJsonResponse(data, status = 200) {
  return new Response(JSON.stringify({ success: true, data }), {
    status,
    headers: { "Content-Type": "application/json" },
  });
}

function mockApiFetch(url, options = {}) {
  const method = (options.method || "GET").toUpperCase();
  const path = String(url || "");

  if (path === "/api/v1/organizations" && method === "GET") {
    return mockJsonResponse(getMockOrganizations());
  }

  if (path === "/api/v1/organizations" && method === "POST") {
    const body = JSON.parse(options.body || "{}");
    return mockJsonResponse(addMockOrganization(body.name));
  }

  const projectsMatch = path.match(/^\/api\/v1\/organizations\/([^/]+)\/projects$/);
  if (projectsMatch && method === "GET") {
    const organizationId = projectsMatch[1];
    return mockJsonResponse(getMockProjects(organizationId));
  }

  if (projectsMatch && method === "POST") {
    const organizationId = projectsMatch[1];
    const body = JSON.parse(options.body || "{}");
    return mockJsonResponse(addMockProject(organizationId, body.name));
  }

  return mockJsonResponse([], 404);
}

export default async function apiFetch(url, options = {}, requireAuth = true) {
  if (isMockAuthEnabled() && typeof url === "string" && url.startsWith("/api/v1/")) {
    return mockApiFetch(url, options);
  }

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
