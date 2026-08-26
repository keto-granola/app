# Keto Granola server 

Keto granola server using:

- Echo (HTTP routing)
- PostgreSQL
- sqlc (type-safe database queries)

## Local Development

### Prerequisites:
- `.env` file
- postgres DB running via docker. See [root README.md](../README.md)

**Important**:
Any new env vars need to be added to both root files:
- `.env`
- `.env.ci` (use dummy values)

### Setup:
```
make dep
```

### Run:
```
make run
```

### Lint:
```
make lint
```

- Fix lint errors:
```
make lint/fix
```

### Tests:
```
make test
```

- Running unit tests only:
`make test/unit`

- Running e2e tests only:
`make test/e2e`

### Generate db queries:
```
make sqlc
```

### Create a db migration:
```
make migrate/create name=<migration_name>
```

### Run db migrations:

Restart the server to apply migrations automatically.

### Generate mocks:
```
make mocks
```
