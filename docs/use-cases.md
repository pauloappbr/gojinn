# 💡 Use Cases and Architecture Patterns

Gojinn is not just an alternative to Docker; it's a solution for specific architectural problems where network latency and idle costs are prohibitive.

Below, we explore ideal scenarios for adopting **In-Process Serverless**.


## 1. Massive Multi-Tenant SaaS

**Scenario:** You operate a SaaS platform (like Shopify, Webflow, or Zapier) and want to allow your users to execute custom scripts or custom business rules.

* **The Traditional Problem:**
    * Running a Docker container per client is impractical (astronomical RAM costs).
    * Using AWS Lambda introduces network latency and high variable costs.
* **The Gojinn Solution:**
    * You can host **thousands of functions** (`.wasm` files) on a single server.
    * **Isolation:** Each execution is sandboxed; one client cannot access another's data.
    * **Density:** Since idle code is just a file on disk, you scale to 10,000 clients with the infrastructure cost of 1 server.


## 2. "Air-Gapped" Environments and Compliance

**Scenario:** Financial institutions, Government, Healthcare, or Industry 4.0 where data **cannot** leave the local infrastructure (On-Premise) for processing on the public cloud.

* **The Traditional Problem:**
    * Modern serverless solutions (Cloudflare Workers, Vercel) require traffic to pass through their infrastructure.
    * Installing a full Kubernetes on-premise just to run scripts is "overkill".
* **The Gojinn Solution:**
    * **Sovereignty:** The runtime runs entirely in your Caddy binary. No "phoning home".
    * **Simplicity:** Deployment is just copying a binary and a configuration file. Works offline on isolated networks.


## 3. High-Performance Middleware

**Scenario:** You need to validate complex payloads (JSON Schema), verify cryptographic signatures (HMAC/JWT), or transform data *before* it reaches your legacy backend.

* **The Traditional Problem:**
    * Adding a "Sidecar" or external API Gateway introduces network hops, increasing latency by precious milliseconds.
* **The Gojinn Solution:**
    * Logic runs **inside the process** of the web server.
    * **Zero-Copy Networking:** Data is passed from Caddy's memory directly to the function's memory.
    * Ideal for: *Advanced Rate Limiting, Custom WAF, XML to JSON Transformation.*


## 4. Edge Computing

**Scenario:** You have geographically distributed servers (own CDNs or cheap VPS) and want to process requests as close to the user as possible.

* **The Gojinn Solution:**
    * Due to the tiny size of WASM binaries (often < 2MB) and low resource consumption, you can run Gojinn on modest hardware (Raspberry Pi, $5 VPS) with enterprise performance.


## 5. Protection Against Infinite Resources

**Scenario:** Execute logic where you don't trust the code (e.g., third-party plugins).

* **The Gojinn Solution:**
    * Unlike a native binary, Gojinn allows you to define **Hard Limits**.
    * If a plugin enters an infinite loop (`while true`), Gojinn terminates it exactly at the configured time limit (e.g., 50ms), protecting your main server's CPU.