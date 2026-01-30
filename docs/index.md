# 🧞 Gojinn Documentation

> **A Sovereign, In-Process Serverless Runtime for Caddy**

Gojinn is an **opinionated WebAssembly runtime** embedded directly into the Caddy web server.

It enables you to execute **untrusted, deterministic, sandboxed code** inside the HTTP request lifecycle — without containers, orchestration layers, or external control planes.

This documentation assumes the reader is comfortable with:
- systems programming concepts
- HTTP internals
- WebAssembly / WASI fundamentals

---

## 🧠 Core Design Invariant

> **All user code executes inside a deterministic, isolated, ephemeral WASM sandbox and is never trusted by default.**

This invariant shapes **every architectural decision** in Gojinn.

If a feature violates this principle, it does not belong in the runtime.

See:
- [`GOVERNANCE.md`](../GOVERNANCE.md)
- [`MANIFESTO.md`](../MANIFESTO.md)

---

## 🚀 Why Gojinn Exists

Traditional serverless platforms optimize for **centralized scale**.  
Gojinn optimizes for **local correctness**.

| Capability | Description |
|-----------|-------------|
| ⚡ **In-Process Execution** | No network hops, no sidecars, no proxies. |
| 🧠 **Deterministic Runtime** | Explicit CPU and memory limits. |
| 🛡️ **Strong Isolation** | Every request runs in a fresh WASM sandbox. |
| 🗄️ **Host-Managed State** | Databases and KV via explicit capabilities. |
| 👁️ **First-Class Observability** | Metrics, tracing, and structured logs built-in. |

---

## 🧭 How to Navigate the Docs

### Getting Started
Start here if this is your first time using Gojinn.
- [Quickstart](./getting-started/quickstart.md)
- [Installation](./getting-started/installation.md)

### Guides
Practical usage and operations.
- [Golang SDK](./guides/golang.md)
- [Deployment & Operations](./guides/deployment.md)
- [Debugging & Observability](./guides/debugging.md)

### Concepts
Deep architectural understanding.
- [Architecture](./concepts/architecture.md)
- [Execution Contract](./concepts/contract.md)

### Reference
Precise definitions.
- [Benchmarks](./benchmark.md)
- [Caddyfile Reference](./reference/caddyfile.md)

---

> Gojinn favors **correctness over convenience**.  
> These docs reflect that philosophy.
