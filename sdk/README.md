@rho-commerce/sdk — TypeScript SDK

Overview

This SDK provides a lightweight, typed client for the Rho API. It exposes modules for authentication, products, cart, orders, payments, inventory, shipping, discounts and checkout.

Install (local development)

From the repo root you can build the SDK and use it locally:

```bash
cd sdk
npm install
npm run build
```

API surface

Import the client and use modules:

```ts
import { createClient } from '@rho-commerce/sdk';

const rho = createClient({ baseURL: 'http://localhost:8080' });

await rho.auth.signup({ email, password, firstName, lastName });
await rho.auth.login({ email, password });
await rho.products.list({ page: 1, limit: 10 });
```

Build/Typecheck

```bash
npm run typecheck
npm run build
```

Publishing

When you're ready to publish the package to npm, update `package.json` version and follow your normal release process. This repo includes a local `file:` dependency usage in the demo to load the built SDK.
