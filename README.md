# Keto Granola

Monorepo for Keto Granola e-commerce platform.

## App architecture
![high-level system architecture](apparchitecture.png)

Cloudflare acts as a full proxy caching static assets and SSR HTML, forwarding cache misses to the backend server. The server renders all routes via `html/template`, with React islands hydrating interactive components 
client-side.

In production, the frontend has no standalone runtime. Its build ouput is embedded into the Go binary and served directly by Echo. The only thing that runs in production is the single backend image.

For local development, a separate Vite dev server runs for hot reload.

## Structure

- `/frontend` — Vite + React islands. Embedded into the Go binary at build time. See [frontend/README.md](./frontend/README.md)
- `/backend` — Go + Echo. See [backend/README.md](./backend/README.md)

## Local Development

### Prerequisite file:
- .env (see [.env.example](.env.example))

**Important:**
- Any new env vars need to be added to `/backend/.env` and `.env.ci`
- Use dummy values for `env.ci` 

### Run everything via Docker:

- Start:
```
make up
```

- Force rebuild after making changes:
```
make up/build
```

- Stop:
```
make down
```

- Stop and remove volumes (wipes DB):
```
make down/vol
```

### Run everything without Docker

The postgres DB always runs via Docker.

Backend:
```
docker compose up -d postgres
cd backend
make dep
make run
```

See [`backend/README.md`](backend/README.md) for linting, testing, etc.

Frontend:
```
cd frontend
make dep
make run
```

See [`frontend/README.md`](frontend/README.md) for linting, testing, etc.


