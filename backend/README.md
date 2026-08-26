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

### Generate mocks:
```
make mocks
```

### Migrations

Create:
```
make migrate/create name=<migration_name>
```

Run:
1. Restart the server to apply migrations automatically.
2. Update `./internal/store/db/schema.sql` so sqlc has the latest db schema.

## MCP Server

The MCP server is mounted on the existing Echo server and exposes store admin operations as tools for AI agents. 

- **Endpoint**: `/mcp`
- **Transport**: Streamable HTTP
- **Auth**: Bearer token via `MCP_AUTH_TOKEN` env var see [`.env.example`](.env.example)

### Available Tools

| Tool                     | Description                                    |
_____________________________________________________________________________
| `get_low_stock_products` | Returns products below their restock threshold |

### Connecting a client

Add to your MCP client config (e.g. Claude Code):

```json
{
  "mcpServers": {
    "keto-granola-store": {
      "url": "http://localhost:${SERVER_PORT}/mcp",
      "headers": {
        "Authorization": "Bearer ${MCP_AUTH_TOKEN}"
      }
    }
  }
}
```

### Testing without an LLM client

Streamable HTTP is session-based so capture the `Mcp-Session-Id` header from the first step and pass it on every subsequent request:

1. Initialise

```bash
curl -s -i http://localhost:$SERVER_PORT/mcp \
  -H "Authorization: Bearer $MCP_AUTH_TOKEN" \
  -H "Content-Type: application/json" \
  -d '{"jsonrpc":"2.0","id":1,"method":"initialize","params":{"protocolVersion":"2025-06-18","capabilities":{},"clientInfo":{"name":"curl-test","version":"1.0"}}}'
```

2. Use tool

```bash
curl -s http://localhost:$SERVER_PORT/mcp \
  -H "Authorization: Bearer $MCP_AUTH_TOKEN" \
  -H "Content-Type: application/json" \
  -H "Mcp-Session-Id: <session_id>" \
  -d '{"jsonrpc":"2.0","id":2,"method":"tools/call","params":{"name":"get_low_stock_products","arguments":{}}}'
```

