Rho Frontend Demo

This demo is a tiny Vite app that exercises the public SDK (signup/login, list products, add to cart).

Run locally

1. Ensure the backend server is running (defaults to `http://localhost:8080`).
2. From the demo folder:

```bash
cd examples/frontend-demo
npm install
npm run dev
```

3. Open http://localhost:4173 and use the UI. Set the base URL to your backend if different.

Notes

- The demo depends on the local SDK package (file:../../sdk). Build the SDK first if you've made changes: `cd sdk && npm install && npm run build`.
- This demo is intentionally minimal to show the SDK usage; use it as a starting point for a real storefront.
