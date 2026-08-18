# Console App — React + Go API

## Structure

```text
idiom/
├── api-services/              # Go monorepo
│   ├── cmd/
│   │   ├── auth/
│   │   │   └── main.go
│   │   └── console/
│   │       └── main.go
│   │
│   ├── web/
│   │   ├── auth/
│   │   │   └── main.go
│   │   └── console/
│   │       └── main.go
│   │
│   ├── api/
│   ├── domains/
│   └── packages/
│
└── app/                       # React/Vite console application
    ├── src/
    ├── public/
    ├── package.json
    ├── package-lock.json
    └── vite.config.js
```

`api-services` and `app` remain **separate development environments**.

---

# Development

## Go API

Run the console API normally:

```bash
cd idiom/api-services
go run ./cmd/console
```

Example:

```text
http://localhost:8081
```

The authentication service remains separate:

```bash
cd idiom/api-services
go run ./cmd/auth
```

---

## React

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

The React development server and Go API remain separate processes.

---

# React API URLs

**Use only the API endpoint path from React.**

Do not hardcode the full backend URL.

Use:

```js
fetch("/api/v1/organizations");
```

Not:

```js
fetch("http://localhost:8081/api/v1/organizations");
```

And not:

```js
fetch("https://idiom-console.onrender.com/api/v1/organizations");
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
Go Console API :8081
```

This means the React application does not need to know whether it is running locally or on Render.

---

# Production

The React application is built as part of the **console Render service build**.

The repository root is used as the build working directory because the console deployment needs access to both:

```text
app/
api-services/
```

The build performs:

```bash
npm ci --prefix app
npm run build --prefix app
go build -tags netgo -ldflags="-s -w" -o console-server ./api-services/cmd/console
```

The React build generates:

```text
app/dist/
```

The Go build generates:

```text
console-server
```

The resulting deployment contains:

```text
idiom/
├── app/
│   └── dist/
│       ├── index.html
│       └── assets/
│
├── api-services/
│   └── ...
│
└── console-server
```

`app/dist/` is generated during deployment and does **not** need to be committed to Git.

---

# Go Static File Serving

The Go console server serves the generated:

```text
app/dist/
```

directory.

The API routes must be registered before the SPA fallback:

```text
/api/...     → Go API handlers
/assets/...  → React static assets
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

This allows direct navigation and browser refreshes on React routes to work correctly.

---

# Production Routing

The production console service uses a single domain:

```text
https://idiom-console.onrender.com/
```

Routing is:

```text
https://idiom-console.onrender.com/
│
├── /dashboard
├── /organizations
├── /projects
│
├── /assets/*
│
└── /api/v1/...
```

The Go server handles the API and serves the React application.

There is **one Render Web Service and one domain** for the console.

---

# API Authentication

The React app should use the same-origin API:

```js
fetch("/api/v1/organizations", {
	credentials: "include",
});
```

For authenticated console APIs, the browser sends the appropriate authentication/refresh-token cookie automatically when applicable.

The React application does not communicate directly with the authentication service for console API operations.

For API-key-related management APIs, the React app talks to the **console API**:

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

The authentication service remains a separate Go service.

---

# Render Configuration

There are two Render Web Services:

```text
idiom-identity
idiom-console
```

Both use the same Git repository:

```text
riyazuddin-1 / idiom
```

and the `main` branch.

---

## `idiom-identity` — Auth Service

The auth service only needs the Go application.

### Root Directory

Set:

```text
api-services
```

This means Render executes the build from:

```text
idiom/api-services/
```

### Build Command

```bash
go build -tags netgo -ldflags="-s -w" -o app ./cmd/auth
```

### Start Command

```bash
./app
```

The resulting configuration is:

```text
Source:
riyazuddin-1 / idiom

Branch:
main

Root Directory:
api-services

Build Command:
go build -tags netgo -ldflags="-s -w" -o app ./cmd/auth

Start Command:
./app
```

The auth service continues to run independently from the console service.

---

# `idiom-console` — React + Go Service

The console service needs both the React application and the Go API.

Therefore, its Render Root Directory must remain **empty**.

Do **not** set the Root Directory to:

```text
api-services
```

because the React application is a sibling directory:

```text
idiom/
├── app/
└── api-services/
```

The console build must access both directories.

---

## Root Directory

Leave empty.

Render executes commands from:

```text
idiom/
```

---

## Build Command

```bash
npm ci --prefix app && npm run build --prefix app && go build -tags netgo -ldflags="-s -w" -o console-server ./api-services/cmd/console
```

This performs three steps:

```text
npm ci --prefix app
        │
        ▼
Install React dependencies

npm run build --prefix app
        │
        ▼
Generate app/dist/

go build ... ./api-services/cmd/console
        │
        ▼
Generate console-server
```

---

## Start Command

```bash
./console-server
```

The resulting configuration is:

```text
Source:
riyazuddin-1 / idiom

Branch:
main

Root Directory:
(empty)

Build Command:
npm ci --prefix app && npm run build --prefix app && go build -tags netgo -ldflags="-s -w" -o console-server ./api-services/cmd/console

Start Command:
./console-server
```

---

# Node.js Build Environment

The console Render service must have Node.js/npm available during the build because it runs:

```bash
npm ci
npm run build
```

The service itself is still a **Go Web Service** because the final production process is:

```text
./console-server
```

Node.js is only required during the build phase.

If the selected Render Go environment does not provide the required Node/npm version, configure the required Node version through the Render build environment or use an appropriate build setup. The runtime remains the Go server.

---

# Render Deployment Flow

When `idiom-console` deploys, Render effectively performs:

```text
Repository
    │
    ├── app/
    │     │
    │     ├── npm ci
    │     │
    │     └── npm run build
    │             │
    │             ▼
    │          app/dist/
    │
    └── api-services/
          │
          └── go build
                  │
                  ▼
             console-server
```

Then Render starts:

```bash
./console-server
```

The Go application serves:

```text
/api/*       → Go API
/assets/*    → React assets
/*           → React SPA
```

---

# Development Architecture

```text
idiom/
│
├── api-services/
│   │
│   ├── cmd/
│   │   ├── auth/
│   │   └── console/
│   │
│   └── web/
│       ├── auth/
│       └── console/
│
└── app/
    └── React/Vite
```

During development:

```text
React/Vite :5173
      │
      │ /api/*
      ▼
Vite proxy
      │
      ▼
Go Console API :8081
```

The auth service remains separate.

---

# Production Architecture

```text
                    idiom-console.onrender.com
                              │
                    ┌─────────┴─────────┐
                    │                   │
                  /api/*                /*
                    │                   │
                    ▼                   ▼
              Go Console API       React dist/
                    │
                    │
              console-server
```

The auth service is deployed separately:

```text
idiom-identity.onrender.com
          │
          ▼
      Go Auth API
```

Therefore:

```text
                    ┌──────────────────────┐
                    │   idiom-console      │
                    │   Render Web Service  │
                    │                      │
Browser ───────────►│ React + Go API       │
                    └──────────────────────┘
                              │
                              │ auth operations
                              ▼
                    ┌──────────────────────┐
                    │   idiom-identity     │
                    │   Render Web Service  │
                    │                      │
                    │   Go Auth API         │
                    └──────────────────────┘
```

---

# Important Rules

### React API rule

Always use endpoint paths:

```js
fetch("/api/v1/organizations");
```

Never use environment-specific full URLs:

```js
fetch("http://localhost:8081/api/v1/organizations");
```

or:

```js
fetch("https://idiom-console.onrender.com/api/v1/organizations");
```

The browser always communicates with the same-origin console domain.

---

### Render root-directory rule

For the auth service:

```text
Root Directory = api-services
```

For the console service:

```text
Root Directory = empty
```

The reason is that the auth build only needs:

```text
api-services/
```

while the console build needs:

```text
app/
api-services/
```

---

### React build rule

React is built during the Render console deployment:

```bash
npm ci --prefix app
npm run build --prefix app
```

which produces:

```text
app/dist/
```

The `dist/` directory does not need to be committed.

---

### Go server rule

The production console process is the Go server:

```bash
./console-server
```

The Go server is responsible for both:

```text
API
+
React static files / SPA fallback
```

---

# Final Repository

```text
idiom/
│
├── api-services/
│   │
│   ├── cmd/
│   │   ├── auth/
│   │   │   └── main.go
│   │   └── console/
│   │       └── main.go
│   │
│   ├── web/
│   │   ├── auth/
│   │   │   └── main.go
│   │   └── console/
│   │       └── main.go
│   │
│   ├── api/
│   │   ├── auth/
│   │   └── console/
│   │
│   ├── domains/
│   └── packages/
│
└── app/
    ├── src/
    ├── public/
    ├── package.json
    ├── package-lock.json
    ├── vite.config.js
    └── dist/                 # generated during build, gitignored
```

The key production model is:

```text
                    Git Repository
                          │
             ┌────────────┴────────────┐
             │                         │
      idiom-identity             idiom-console
             │                         │
    Root: api-services          Root: repository
             │                         │
       Go auth build          React build + Go build
             │                         │
             ▼                         ▼
       ./app binary             ./console-server
                                       │
                              ┌────────┴────────┐
                              │                 │
                           /api/*              /*
                              │                 │
                              ▼                 ▼
                         Go handlers        app/dist/
```

This keeps the **Go monorepo and React app separate for development**, while the production console remains a **single Render Web Service and single same-origin domain**.
