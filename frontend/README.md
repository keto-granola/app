# Keto Granola client

React client for keto granola e-commerce platform using:

- React
- Typescript
- Vite
- TanStack (routing & queries)
- shadcn/ui (UI framework)
- Zod (validation)

## Local Development

### Prerequisite files:
- .env

### Setup:
```
make dep
```

### Run:
```
make run
```

**Note:** There's no `index.html`, so opening `localhost:5173` directly in a browser won't show anything. Just view the app through the Go server (`localhost:3001`) after Vite is running.

### Lint:
```
make lint
```

- Fix lint errors:
```
make lint/fix
```

- Check Typescript errors:
```
make typecheck
```

### Build:
```
make build
```