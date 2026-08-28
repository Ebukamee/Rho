import './styles.css';
import { createClient } from '@rho-commerce/sdk';

type AppState = {
  baseURL: string;
  token: string;
  products: any[];
  cart: any;
  output: string;
};

const state: AppState = {
  baseURL: 'http://localhost:8080',
  token: '',
  products: [],
  cart: null,
  output: 'Ready to connect to Rho.',
};

const app = document.querySelector('#app');
if (!app) throw new Error('App root not found');

const client = createClient({ baseURL: state.baseURL });

function render() {
  app.innerHTML = `
    <div class="container">
      <header>
        <div>
          <p class="badge">Rho Demo</p>
          <h1>Frontend happy-path demo</h1>
        </div>
        <div class="inline">
          <input id="base-url" value="${state.baseURL}" aria-label="Base URL" />
          <button id="connect-btn" class="secondary">Connect</button>
        </div>
      </header>

      <section class="card">
        <h2>Authentication</h2>
        <div class="controls">
          <div class="form-grid">
            <label>
              Email
              <input id="email" type="email" placeholder="user@example.com" value="demo@example.com" />
            </label>
            <label>
              First name
              <input id="first-name" value="Demo" />
            </label>
            <label>
              Last name
              <input id="last-name" value="User" />
            </label>
            <label>
              Password
              <input id="password" type="password" value="password123" />
            </label>
            <div class="inline">
              <button id="signup-btn">Sign up</button>
              <button id="login-btn" class="secondary">Log in</button>
            </div>
          </div>

          <div class="stack">
            <div class="form-grid">
              <label>
                Access token
                <input id="token" value="${state.token}" readonly />
              </label>
              <div class="inline">
                <button id="products-btn" class="secondary">Load products</button>
                <button id="cart-btn" class="secondary">Load cart</button>
              </div>
            </div>
          </div>
        </div>
      </section>

      <section class="card">
        <h2>Products</h2>
        <div id="products" class="list"></div>
      </section>

      <section class="card">
        <h2>Cart</h2>
        <div id="cart" class="list"></div>
      </section>

      <section class="card">
        <h2>Console output</h2>
        <pre id="output" class="output">${state.output}</pre>
      </section>
    </div>
  `;

  const baseUrlInput = document.querySelector('#base-url') as HTMLInputElement;
  const connectBtn = document.querySelector('#connect-btn') as HTMLButtonElement;
  const signupBtn = document.querySelector('#signup-btn') as HTMLButtonElement;
  const loginBtn = document.querySelector('#login-btn') as HTMLButtonElement;
  const productsBtn = document.querySelector('#products-btn') as HTMLButtonElement;
  const cartBtn = document.querySelector('#cart-btn') as HTMLButtonElement;
  const tokenInput = document.querySelector('#token') as HTMLInputElement;

  connectBtn.onclick = () => {
    state.baseURL = baseUrlInput.value.trim() || 'http://localhost:8080';
    state.output = `Connected to ${state.baseURL}`;
    client.baseURL = state.baseURL;
    render();
  };

  signupBtn.onclick = async () => {
    try {
      const result = await createClient({ baseURL: state.baseURL }).auth.signup({
        email: (document.querySelector('#email') as HTMLInputElement).value,
        password: (document.querySelector('#password') as HTMLInputElement).value,
        firstName: (document.querySelector('#first-name') as HTMLInputElement).value,
        lastName: (document.querySelector('#last-name') as HTMLInputElement).value,
      });

      state.output = JSON.stringify(result, null, 2);
      state.token = result.accessToken;
      tokenInput.value = state.token;
      client.setToken(state.token);
      render();
    } catch (error) {
      state.output = `Signup failed: ${error instanceof Error ? error.message : String(error)}`;
      render();
    }
  };

  loginBtn.onclick = async () => {
    try {
      const result = await createClient({ baseURL: state.baseURL }).auth.login({
        email: (document.querySelector('#email') as HTMLInputElement).value,
        password: (document.querySelector('#password') as HTMLInputElement).value,
      });

      state.output = JSON.stringify(result, null, 2);
      state.token = result.accessToken;
      tokenInput.value = state.token;
      client.setToken(state.token);
      render();
    } catch (error) {
      state.output = `Login failed: ${error instanceof Error ? error.message : String(error)}`;
      render();
    }
  };

  productsBtn.onclick = async () => {
    try {
      const response = await client.products.list({ page: 1, limit: 10 });
      state.products = response.items ?? [];
      state.output = JSON.stringify(response, null, 2);
      renderProducts();
      render();
    } catch (error) {
      state.output = `Products failed: ${error instanceof Error ? error.message : String(error)}`;
      render();
    }
  };

  cartBtn.onclick = async () => {
    try {
      const response = await client.cart.get();
      state.cart = response;
      state.output = JSON.stringify(response, null, 2);
      renderCart();
      render();
    } catch (error) {
      state.output = `Cart failed: ${error instanceof Error ? error.message : String(error)}`;
      render();
    }
  };

  renderProducts();
  renderCart();
}

function renderProducts() {
  const productsRoot = document.querySelector('#products');
  if (!productsRoot) return;

  if (!state.products.length) {
    productsRoot.innerHTML = '<p>No products loaded yet.</p>';
    return;
  }

  productsRoot.innerHTML = state.products
    .map(
      (product: any) => `
        <div class="product">
          <div>
            <strong>${product.name}</strong><br />
            <small>${product.slug}</small>
          </div>
          <div class="inline">
            <span>${product.price / 100} ${product.currency}</span>
            <button class="secondary" data-product-id="${product.id}">Add to cart</button>
          </div>
        </div>
      `
    )
    .join('');

  productsRoot.querySelectorAll('[data-product-id]').forEach((button) => {
    button.addEventListener('click', async () => {
      const productId = button.getAttribute('data-product-id');
      if (!productId) return;

      try {
        const response = await client.cart.addItem({ productId, quantity: 1 });
        state.cart = response;
        state.output = JSON.stringify(response, null, 2);
        renderCart();
        render();
      } catch (error) {
        state.output = `Add-to-cart failed: ${error instanceof Error ? error.message : String(error)}`;
        render();
      }
    });
  });
}

function renderCart() {
  const cartRoot = document.querySelector('#cart');
  if (!cartRoot) return;

  if (!state.cart) {
    cartRoot.innerHTML = '<p>Your cart is empty.</p>';
    return;
  }

  const items = state.cart.items ?? [];
  if (!items.length) {
    cartRoot.innerHTML = '<p>Your cart is empty.</p>';
    return;
  }

  cartRoot.innerHTML = items
    .map(
      (item: any) => `
        <div class="cart-item">
          <div>
            <strong>Product ${item.productId}</strong><br />
            <small>Qty: ${item.quantity}</small>
          </div>
          <span>${item.totalPrice / 100}</span>
        </div>
      `
    )
    .join('');
}

render();
