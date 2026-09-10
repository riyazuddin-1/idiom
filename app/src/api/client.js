import { getAccessToken, isMockAuthEnabled, reAuthenticate } from "./auth";
import {
  addMockApiKey,
  addMockOrganization,
  addMockProject,
  getMockApiKeys,
  getMockOrganizations,
  getMockProject,
  getMockProjects,
  mockRevokeApiKey,
  mockRotateApiKey,
  mockUpdateApiKey,
  mockUpdateProject,
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

  const projectMatch = path.match(
    /^\/api\/v1\/organizations\/([^/]+)\/projects\/([^/]+)$/,
  );
  if (projectMatch && method === "GET") {
    const organizationId = projectMatch[1];
    const projectId = projectMatch[2];
    const project = getMockProject(organizationId, projectId);
    if (!project) return mockJsonResponse({ message: "Project not found" }, 404);
    return mockJsonResponse(project);
  }

  if (projectMatch && method === "PUT") {
    const organizationId = projectMatch[1];
    const projectId = projectMatch[2];
    const body = JSON.parse(options.body || "{}");
    const project = mockUpdateProject(organizationId, projectId, body);
    if (!project) return mockJsonResponse({ message: "Project not found" }, 404);
    return mockJsonResponse(project);
  }

  const apiKeysListMatch = path.match(
    /^\/api\/v1\/organizations\/[^/]+\/projects\/([^/]+)\/api-keys$/,
  );
  if (apiKeysListMatch && method === "GET") {
    const projectId = apiKeysListMatch[1];
    return mockJsonResponse(getMockApiKeys(projectId));
  }

  if (apiKeysListMatch && method === "POST") {
    const projectId = apiKeysListMatch[1];
    const body = JSON.parse(options.body || "{}");
    return mockJsonResponse(addMockApiKey(projectId, body.name, body.scopes), 201);
  }

  const apiKeyRotateMatch = path.match(
    /^\/api\/v1\/organizations\/[^/]+\/projects\/([^/]+)\/api-keys\/([^/]+)\/rotate$/,
  );
  if (apiKeyRotateMatch && method === "POST") {
    const projectId = apiKeyRotateMatch[1];
    const keyId = apiKeyRotateMatch[2];
    const result = mockRotateApiKey(projectId, keyId);
    if (!result) return mockJsonResponse({ message: "Key not found" }, 404);
    return mockJsonResponse(result);
  }

  const apiKeyRevokeMatch = path.match(
    /^\/api\/v1\/organizations\/[^/]+\/projects\/([^/]+)\/api-keys\/([^/]+)$/,
  );
  if (apiKeyRevokeMatch && method === "PUT") {
    const projectId = apiKeyRevokeMatch[1];
    const keyId = apiKeyRevokeMatch[2];
    const body = JSON.parse(options.body || "{}");
    const result = mockUpdateApiKey(projectId, keyId, body);
    if (!result) return mockJsonResponse({ message: "API key not found" }, 404);
    return mockJsonResponse(result);
  }

  if (apiKeyRevokeMatch && method === "DELETE") {
    const projectId = apiKeyRevokeMatch[1];
    const keyId = apiKeyRevokeMatch[2];
    mockRevokeApiKey(projectId, keyId);
    return new Response(null, { status: 204 });
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
