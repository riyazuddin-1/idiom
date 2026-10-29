const STORAGE_KEY = "idiom_mock_data";

const defaultData = {
  organizations: [
    {
      id: "org-1",
      name: "Northstar Studio",
      status: "active",
      createdAt: "2025-03-01T09:00:00.000Z",
    },
    {
      id: "org-2",
      name: "Launchpad Labs",
      status: "active",
      createdAt: "2025-04-11T15:30:00.000Z",
    },
  ],
  projects: {
    "org-1": [
      {
        "id": "prj_1a0180b99de",
        "org_id": "org_1a01805c89b",
        "name": "notelet",
        "slug": "notelet",
        "description": "",
        "status": "active",
        "auth_options": [
            "email"
        ],
        "redirect_urls": [
            {
                "device": "web",
                "url": "https://notelet-backend.vercel.app/auth/callback"
            },
            {
                "device": "mobile",
                "url": "https://notelet-backend.vercel.app/auth/callback"
            }
        ],
        "allowed_origins": [],
        "created_at": "2026-08-19T03:23:29.525282Z",
        "updated_at": "2026-08-19T03:23:29.525282Z"
    },
      {
        id: "proj-admin",
        name: "Admin Portal",
        slug: "noteletadmin",
        status: "active",
        createdAt: "2025-03-15T11:00:00.000Z",
      },
    ],
    "org-2": [
      {
        id: "proj-mobile",
        name: "Mobile App",
        slug: "noteletmobile",
        status: "active",
        createdAt: "2025-04-18T09:00:00.000Z",
      },
    ],
  },
  apiKeys: {
    "proj-web": [
      {
        id: "ak-1",
        name: "Staging",
        project_id: "proj-web",
        key_prefix: "ik_AbCdEfGhI",
        scopes: ["read", "write"],
        is_active: true,
        created_by: "user-1",
        created_at: "2025-03-03T10:00:00.000Z",
      },
      {
        id: "ak-2",
        name: "CI/CD Bot",
        project_id: "proj-web",
        key_prefix: "ik_012345678",
        scopes: ["read"],
        is_active: false,
        created_by: "user-1",
        created_at: "2025-06-10T12:00:00.000Z",
        revoked_at: "2025-08-01T09:00:00.000Z",
      },
    ],
    "proj-admin": [],
    "proj-mobile": [
      {
        id: "ak-3",
        name: "Production",
        project_id: "proj-mobile",
        key_prefix: "ik_ZyXwVuTsR",
        scopes: ["read", "write"],
        is_active: true,
        created_by: "user-1",
        created_at: "2025-04-20T09:00:00.000Z",
      },
    ],
  },
};

function readStorage() {
  try {
    const raw = localStorage.getItem(STORAGE_KEY);
    return raw ? JSON.parse(raw) : defaultData;
  } catch {
    return defaultData;
  }
}

function saveStorage(data) {
  localStorage.setItem(STORAGE_KEY, JSON.stringify(data));
}

export function getMockOrganizations() {
  const data = readStorage();
  return [...data.organizations];
}

export function getMockProjects(organizationId) {
  const data = readStorage();
  return [...(data.projects[organizationId] || [])];
}

export function getMockProject(organizationId, projectId) {
  const data = readStorage();
  const projects = data.projects[organizationId] || [];
  return projects.find((p) => p.id === projectId) || null;
}

export function addMockOrganization(name) {
  const data = readStorage();
  const org = {
    id: `org-${Date.now()}`,
    name,
    status: "active",
    createdAt: new Date().toISOString(),
  };

  data.organizations = [org, ...data.organizations];
  data.projects[org.id] = [];
  saveStorage(data);
  return org;
}

export function addMockProject(organizationId, name) {
  const data = readStorage();
  const project = {
    id: `proj-${Date.now()}`,
    name,
    status: "active",
    createdAt: new Date().toISOString(),
  };

  if (!data.projects[organizationId]) {
    data.projects[organizationId] = [];
  }

  data.projects[organizationId] = [project, ...data.projects[organizationId]];
  saveStorage(data);
  return project;
}

export function getMockApiKeys(projectId) {
  const data = readStorage();
  return [...(data.apiKeys[projectId] || [])];
}

export function addMockApiKey(projectId, name, scopes) {
  const data = readStorage();
  const keyId = `ak-${Date.now()}`;
  const key = `ik_${btoa(`${keyId}${Date.now()}`)
    .replace(/[^a-zA-Z0-9]/g, "")
    .slice(0, 32)}`;

  const apiKey = {
    id: keyId,
    name,
    project_id: projectId,
    key_prefix: key.slice(0, 11),
    scopes: scopes || [],
    is_active: true,
    created_by: "user-1",
    created_at: new Date().toISOString(),
  };

  if (!data.apiKeys) {
    data.apiKeys = {};
  }
  if (!data.apiKeys[projectId]) {
    data.apiKeys[projectId] = [];
  }

  data.apiKeys[projectId] = [apiKey, ...data.apiKeys[projectId]];
  saveStorage(data);
  return { api_key: apiKey, key };
}

export function mockRotateApiKey(projectId, keyId) {
  const data = readStorage();
  const keys = data.apiKeys[projectId] || [];
  const index = keys.findIndex((k) => k.id === keyId);
  if (index === -1) return null;

  const key = `ik_${btoa(`${keyId}rotated${Date.now()}`)
    .replace(/[^a-zA-Z0-9]/g, "")
    .slice(0, 32)}`;
  const apiKey = { ...keys[index], key_prefix: key.slice(0, 11) };
  data.apiKeys[projectId][index] = apiKey;
  saveStorage(data);
  return { api_key: apiKey, key };
}

export function mockRevokeApiKey(projectId, keyId) {
  const data = readStorage();
  const keys = data.apiKeys[projectId] || [];
  const index = keys.findIndex((k) => k.id === keyId);
  if (index === -1) return null;

  const apiKey = {
    ...keys[index],
    is_active: false,
    revoked_at: new Date().toISOString(),
  };
  data.apiKeys[projectId][index] = apiKey;
  saveStorage(data);
  return apiKey;
}

export function mockUpdateProject(organizationId, projectId, updates) {
  const data = readStorage();
  const projects = data.projects[organizationId] || [];
  const index = projects.findIndex((p) => p.id === projectId);
  if (index === -1) return null;

  const project = { ...projects[index], ...updates };
  data.projects[organizationId][index] = project;
  saveStorage(data);
  return project;
}

export function mockUpdateApiKey(projectId, keyId, updates) {
  const data = readStorage();
  const keys = data.apiKeys[projectId] || [];
  const index = keys.findIndex((k) => k.id === keyId);
  if (index === -1) return null;

  const apiKey = { ...keys[index], ...updates };
  data.apiKeys[projectId][index] = apiKey;
  saveStorage(data);
  return apiKey;
}

export function resetMockData() {
  localStorage.removeItem(STORAGE_KEY);
}
