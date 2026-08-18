# Console App — React + Go API

## Structure

```text
idiom/
├── api-services/          # Go monorepo
│   ├── api/
│   │   ├── auth/
│   │   └── console/
│   ├── domains/
│   └── packages/
│
└── app/                   # React/Vite console application
    ├── src/
    ├── public/
    ├── package.json
    ├── package-lock.json
    └── vite.config.js
```

`api-services` and `app` remain **separate development environments**.

---

## Development

### Go API

Run the console API normally:

```bash
cd idiom/api-services
go run ./cmd/console
```

Example:

```text
http://localhost:8081
```

### React

Run the React app normally:

```bash
cd idiom/app
npm run dev
```

Example:

```text
http://localhost:5173
```

Configure Vite to proxy API requests:

```js
// app/vite.config.js

import { defineConfig } from "vite";
import react from "@vitejs/plugin-react";

export default defineConfig({
	plugins: [react()],
	server: {
		proxy: {
			"/api": {
				target: "http://localhost:8081",
				changeOrigin: true,
			},
		},
	},
});
```

---

## React API URLs

**Use only the API endpoint path from React.**

Do **not** hardcode the full backend URL.

Use:

```js
fetch("/api/v1/organizations");
```

Not:

```js
fetch("http://localhost:8080/api/v1/organizations");
```

And not:

```js
fetch("https://idiom-console-api.onrender.com/api/v1/organizations");
```

The browser should always use the same-origin path:

```text
/api/v1/...
```

During development, Vite proxies it:

```text
React
  │
  │ /api/v1/organizations
  ▼
Vite :5173
  │
  │ proxy
  ▼
Go :8080
```

---

# Production

The React application is built before the Go application is deployed.

```bash
cd idiom/app
npm ci
npm run build
```

This generates:

```text
idiom/app/dist/
```

`dist/` should remain in `.gitignore`.

Then build the Go console API.

The Go server serves:

```text
/api/*       → Go API
/assets/*    → React build
/*           → React index.html
```

Production becomes:

```text
https://idiom-console.onrender.com/
│
├── /dashboard
├── /organizations
├── /projects
│
└── /api/v1/...
```

There is **one Render Web Service and one domain**.

---

## Go Static File Serving

The Go console server should serve the generated:

```text
app/dist/
```

directory.

The API routes must be registered before the SPA fallback:

```text
/api/...     → Go handlers
/*           → React application
```

For React Router routes such as:

```text
/dashboard
/organizations
/projects
```

the Go server should return:

```text
app/dist/index.html
```

when the requested file does not exist.

---

## Render Build

The build needs to produce both the React assets and Go binary.

Conceptually:

```bash
cd app
npm ci
npm run build

cd ../api-services
go build -o console-server ./api/console
```

The resulting deployment contains:

```text
idiom/
├── app/
│   └── dist/              # generated during Render build
│
└── api-services/
    └── console-server
```

`dist/` does **not** need to be committed.

---

## API Authentication

The React app should use the same-origin API:

```js
fetch("/api/v1/organizations", {
	credentials: "include",
});
```

For authenticated console APIs, the browser sends the refresh-token cookie automatically when applicable.

For API-key-related management APIs, the React app is talking to the **console API**, not directly to the auth service.

```text
React
  │
  │ /api/v1/...
  ▼
Console API
  │
  ├── organizations
  ├── members
  ├── projects
  └── API keys
```

The auth service remains a separate Go API service.

---

## Final Architecture

```text
idiom/
│
├── api-services/
│   │
│   ├── api/
│   │   ├── auth/         # Authentication API
│   │   └── console/      # Console API
│   │
│   ├── domains/
│   └── packages/
│
└── app/                  # React/Vite console
    ├── src/
    ├── public/
    ├── package.json
    ├── vite.config.js
    └── dist/             # generated, gitignored
```

### Development

```text
React/Vite :5173
      │
      │ /api/*
      ▼
Vite proxy
      │
      ▼
Go Console API :8080
```

### Production

```text
                 idiom-console.onrender.com
                            │
              ┌─────────────┴─────────────┐
              │                           │
           /api/*                       /*
              │                           │
              ▼                           ▼
        Go Console API              React dist/
```

**React rule:** always use `/api/...` endpoint paths, never environment-specific full URLs.

Done:
2. Build React console pages — Console shell with sidebar nav, Dashboard, Organizations, Projects pages created.
3. Fix auth callback — Callback.jsx now uses `authenticate()` from auth.js with same-origin `/api/v1/token`.
4. Render build config — `render.yaml` created at project root, chains `npm ci && npm run build` in `app/` then `go build ./cmd/console` in `api-services/`.