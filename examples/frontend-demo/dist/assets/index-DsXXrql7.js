(function(){const t=document.createElement("link").relList;if(t&&t.supports&&t.supports("modulepreload"))return;for(const i of document.querySelectorAll('link[rel="modulepreload"]'))a(i);new MutationObserver(i=>{for(const n of i)if(n.type==="childList")for(const c of n.addedNodes)c.tagName==="LINK"&&c.rel==="modulepreload"&&a(c)}).observe(document,{childList:!0,subtree:!0});function r(i){const n={};return i.integrity&&(n.integrity=i.integrity),i.referrerPolicy&&(n.referrerPolicy=i.referrerPolicy),i.crossOrigin==="use-credentials"?n.credentials="include":i.crossOrigin==="anonymous"?n.credentials="omit":n.credentials="same-origin",n}function a(i){if(i.ep)return;i.ep=!0;const n=r(i);fetch(i.href,n)}})();class f extends Error{constructor(t,r){super(r.error??r.message??"Request failed"),this.name="RhoApiError",this.status=t,this.body=r}}function m(e){return{signup:t=>e.post("/api/v1/auth/signup",t),login:t=>e.post("/api/v1/auth/login",t),refresh:t=>e.post("/api/v1/auth/refresh",t),logout:()=>e.post("/api/v1/auth/logout"),getProfile:()=>e.get("/api/v1/auth/me"),updateProfile:t=>e.put("/api/v1/auth/profile",t),changePassword:t=>e.put("/api/v1/auth/password",t),googleLogin:()=>e.get("/api/v1/auth/google/login")}}function y(e){return{get:()=>e.get("/api/v1/cart"),addItem:t=>e.post("/api/v1/cart/items",t),updateItem:(t,r)=>e.put(`/api/v1/cart/items/${t}`,r),removeItem:t=>e.delete(`/api/v1/cart/items/${t}`),clear:()=>e.delete("/api/v1/cart")}}function b(e){return{list:t=>e.get("/api/v1/categories",t),get:t=>e.get(`/api/v1/categories/${t}`),adminList:t=>e.get("/api/v1/categories/admin",t),create:t=>e.post("/api/v1/categories",t),update:(t,r)=>e.put(`/api/v1/categories/${t}`,r),delete:t=>e.delete(`/api/v1/categories/${t}`)}}function $(e){return{preview:t=>e.post("/api/v1/checkout/preview",t),create:t=>e.post("/api/v1/checkout",t)}}function L(e){return{apply:t=>e.post("/api/v1/discounts/apply",t),create:t=>e.post("/api/v1/discounts",t),get:t=>e.get(`/api/v1/discounts/${t}`),update:(t,r)=>e.put(`/api/v1/discounts/${t}`,r),delete:t=>e.delete(`/api/v1/discounts/${t}`)}}function S(e){return{create:t=>e.post("/api/v1/inventory",t),get:t=>e.get(`/api/v1/inventory/${t}`),getByProduct:t=>e.get(`/api/v1/inventory/product/${t}`),update:(t,r)=>e.put(`/api/v1/inventory/${t}`,r),adjust:(t,r)=>e.post(`/api/v1/inventory/product/${t}/adjust`,r),delete:t=>e.delete(`/api/v1/inventory/${t}`)}}function k(e){return{create:t=>e.post("/api/v1/orders",t),get:t=>e.get(`/api/v1/orders/${t}`),updateStatus:(t,r)=>e.put(`/api/v1/orders/${t}/status`,r)}}function q(e){return{initialize:t=>e.post("/api/v1/payments/initialize",t),get:t=>e.get(`/api/v1/payments/${t}`),verify:(t,r)=>e.post(`/api/v1/payments/${t}/verify`,r??{})}}function w(e){return{list:t=>e.get("/api/v1/products",t),get:t=>e.get(`/api/v1/products/${t}`),adminList:t=>e.get("/api/v1/products/admin",t),create:t=>e.post("/api/v1/products",t),update:(t,r)=>e.put(`/api/v1/products/${t}`,r),delete:t=>e.delete(`/api/v1/products/${t}`)}}function R(e){return{create:t=>e.post("/api/v1/shipping",t),get:t=>e.get(`/api/v1/shipping/${t}`),getByOrder:t=>e.get(`/api/v1/shipping/order/${t}`),update:(t,r)=>e.put(`/api/v1/shipping/${t}`,r),delete:t=>e.delete(`/api/v1/shipping/${t}`)}}function U(e){return{list:t=>e.get("/api/v1/users",t),get:t=>e.get(`/api/v1/users/${t}`),update:(t,r)=>e.put(`/api/v1/users/${t}`,r),delete:t=>e.delete(`/api/v1/users/${t}`),updateRole:(t,r)=>e.put(`/api/v1/users/${t}/role`,r)}}class P{constructor(t){this.auth=m(this),this.products=w(this),this.categories=b(this),this.cart=y(this),this.orders=k(this),this.payments=q(this),this.shipping=R(this),this.discounts=L(this),this.checkout=$(this),this.inventory=S(this),this.users=U(this),this.baseURL=t.baseURL.replace(/\/+$/,""),this.token=t.token,this.fetchFn=t.fetch??fetch,this.headers={"Content-Type":"application/json",...t.headers??{}}}setToken(t){this.token=t}clearToken(){this.token=void 0}buildUrl(t,r){const a=new URL(t.startsWith("http")?t:`${this.baseURL}${t}`);return r&&Object.entries(r).forEach(([i,n])=>{n!=null&&a.searchParams.set(i,String(n))}),a.toString()}async request(t,r,a,i,n){const c=this.buildUrl(r,i),s=await this.fetchFn(c,{method:t,headers:{...this.headers,...n,...this.token?{Authorization:`Bearer ${this.token}`}:{}},body:a!==void 0?JSON.stringify(a):void 0}),h=(s.headers.get("content-type")??"").includes("application/json")?await s.json():await s.text();if(!s.ok)throw new f(s.status,h);return h}get(t,r){return this.request("GET",t,void 0,r)}post(t,r,a){return this.request("POST",t,r,a)}put(t,r,a){return this.request("PUT",t,r,a)}patch(t,r,a){return this.request("PATCH",t,r,a)}delete(t,r){return this.request("DELETE",t,void 0,r)}}function p(e){return new P(e)}const o={baseURL:"http://localhost:8080",token:"",products:[],cart:null,output:"Ready to connect to Rho."},g=document.querySelector("#app");if(!g)throw new Error("App root not found");const d=p({baseURL:o.baseURL});function u(){g.innerHTML=`
    <div class="container">
      <header>
        <div>
          <p class="badge">Rho Demo</p>
          <h1>Frontend happy-path demo</h1>
        </div>
        <div class="inline">
          <input id="base-url" value="${o.baseURL}" aria-label="Base URL" />
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
                <input id="token" value="${o.token}" readonly />
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
        <pre id="output" class="output">${o.output}</pre>
      </section>
    </div>
  `;const e=document.querySelector("#base-url"),t=document.querySelector("#connect-btn"),r=document.querySelector("#signup-btn"),a=document.querySelector("#login-btn"),i=document.querySelector("#products-btn"),n=document.querySelector("#cart-btn"),c=document.querySelector("#token");t.onclick=()=>{o.baseURL=e.value.trim()||"http://localhost:8080",o.output=`Connected to ${o.baseURL}`,d.baseURL=o.baseURL,u()},r.onclick=async()=>{try{const s=await p({baseURL:o.baseURL}).auth.signup({email:document.querySelector("#email").value,password:document.querySelector("#password").value,firstName:document.querySelector("#first-name").value,lastName:document.querySelector("#last-name").value});o.output=JSON.stringify(s,null,2),o.token=s.accessToken,c.value=o.token,d.setToken(o.token),u()}catch(s){o.output=`Signup failed: ${s instanceof Error?s.message:String(s)}`,u()}},a.onclick=async()=>{try{const s=await p({baseURL:o.baseURL}).auth.login({email:document.querySelector("#email").value,password:document.querySelector("#password").value});o.output=JSON.stringify(s,null,2),o.token=s.accessToken,c.value=o.token,d.setToken(o.token),u()}catch(s){o.output=`Login failed: ${s instanceof Error?s.message:String(s)}`,u()}},i.onclick=async()=>{try{const s=await d.products.list({page:1,limit:10});o.products=s.items??[],o.output=JSON.stringify(s,null,2),v(),u()}catch(s){o.output=`Products failed: ${s instanceof Error?s.message:String(s)}`,u()}},n.onclick=async()=>{try{const s=await d.cart.get();o.cart=s,o.output=JSON.stringify(s,null,2),l(),u()}catch(s){o.output=`Cart failed: ${s instanceof Error?s.message:String(s)}`,u()}},v(),l()}function v(){const e=document.querySelector("#products");if(e){if(!o.products.length){e.innerHTML="<p>No products loaded yet.</p>";return}e.innerHTML=o.products.map(t=>`
        <div class="product">
          <div>
            <strong>${t.name}</strong><br />
            <small>${t.slug}</small>
          </div>
          <div class="inline">
            <span>${t.price/100} ${t.currency}</span>
            <button class="secondary" data-product-id="${t.id}">Add to cart</button>
          </div>
        </div>
      `).join(""),e.querySelectorAll("[data-product-id]").forEach(t=>{t.addEventListener("click",async()=>{const r=t.getAttribute("data-product-id");if(r)try{const a=await d.cart.addItem({productId:r,quantity:1});o.cart=a,o.output=JSON.stringify(a,null,2),l(),u()}catch(a){o.output=`Add-to-cart failed: ${a instanceof Error?a.message:String(a)}`,u()}})})}}function l(){const e=document.querySelector("#cart");if(!e)return;if(!o.cart){e.innerHTML="<p>Your cart is empty.</p>";return}const t=o.cart.items??[];if(!t.length){e.innerHTML="<p>Your cart is empty.</p>";return}e.innerHTML=t.map(r=>`
        <div class="cart-item">
          <div>
            <strong>Product ${r.productId}</strong><br />
            <small>Qty: ${r.quantity}</small>
          </div>
          <span>${r.totalPrice/100}</span>
        </div>
      `).join("")}u();
