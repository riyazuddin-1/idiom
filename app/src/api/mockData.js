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
        id: "proj-web",
        name: "Web App",
        status: "active",
        createdAt: "2025-03-02T10:00:00.000Z",
      },
      {
        id: "proj-admin",
        name: "Admin Portal",
        status: "active",
        createdAt: "2025-03-15T11:00:00.000Z",
      },
    ],
    "org-2": [
      {
        id: "proj-mobile",
        name: "Mobile App",
        status: "active",
        createdAt: "2025-04-18T09:00:00.000Z",
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

export function resetMockData() {
  localStorage.removeItem(STORAGE_KEY);
}
