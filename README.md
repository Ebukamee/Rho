Rho — Self-hosted E-commerce Backend

Overview

Rho is a Go-based, self-hosted backend that provides production-ready e-commerce APIs (auth, products, cart, checkout, payments, inventory, shipping, discounts). This repository contains the server, a TypeScript SDK, and a small frontend demo.

Quick start (developer running locally)

1. Build & run the backend

- Set environment variables (create a `.env` or export):

  ```env
  DATABASE_URL=postgres://user:pass@localhost:5432/rho?sslmode=disable
  JWT_SECRET=your_jwt_secret
  PORT=8080
  RATE_LIMIT_REQUESTS=60
  RATE_LIMIT_WINDOW_SECONDS=60
  ```

- Run migrations (using `psql` or your migration tool):

  ```bash
  cd migrations
  # example using psql
  psql "$DATABASE_URL" < 000001_create_users.up.sql
  psql "$DATABASE_URL" < 000002_create_products.up.sql
  # ... apply remaining .up.sql files in sequence
  ```

- Build and run the Go server:

  ```bash
  go build ./...
  go run cmd/api/main.go
  ```

2. Build the TypeScript SDK (optional)

```bash
cd sdk
npm install
npm run build
```

3. Run the demo frontend (optional)

```bash
cd examples/frontend-demo
npm install
npm run dev
# open http://localhost:4173
```

Repository layout

- `cmd/api` — server entry point
- `internal/*` — application modules (auth, product, cart, order, payment, shipping, etc.)
- `migrations/` — SQL migrations
- `sdk/` — TypeScript SDK for frontend developers
- `examples/frontend-demo` — tiny Vite demo app demonstrating the SDK

Support & contributions

If you want to contribute, open a PR and include tests for changes to the server logic.
