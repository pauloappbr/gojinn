# 🧞 Gojinn

> **High-Performance Serverless Runtime for Caddy**

[![Go Reference](https://pkg.go.dev/badge/github.com/caddyserver/caddy/v2.svg)](https://pkg.go.dev/github.com/pauloappbr/gojinn)
[![Build Status](https://img.shields.io/badge/build-passing-brightgreen)]()
[![Wasm Engine](https://img.shields.io/badge/engine-wazero-purple)](https://wazero.io)
[![Version](https://img.shields.io/badge/version-v0.3.0-blue)]()

**Gojinn** is an *in-process serverless* runtime for the [Caddy](https://caddyserver.com) web server.

It allows you to execute **Go**, **Rust**, **Zig**, and **C++** code (compiled to WebAssembly) directly in the HTTP request flow.

With the release of **v0.3.0**, Gojinn features a **JIT (Just-In-Time) Caching Engine**, eliminating the *cold start* problem entirely and achieving **microsecond latency**.

---

## 🚀 Why Gojinn?

Traditional serverless introduces network latency. Gojinn brings computation closer to the data.

| Feature | Description |
| :--- | :--- |
| **⚡ Microsecond Latency** | JIT Caching & Buffer Pooling ensure execution in **< 1ms**. (vs 1500ms+ for Docker). |
| **🏗️ Zero Infra** | No Docker daemon, no Kubernetes sidecars. It's just a Caddy plugin. |
| **👁️ Observable** | Native support for **Prometheus Metrics**, **OpenTelemetry Tracing**, and Structured Logging. |
| **🛡️ Secure** | Each request runs in a strict Sandbox via [Wazero](https://wazero.io). If the code crashes, Caddy stays alive. |
| **💰 Zero Idle Cost** | Functions are just files on disk (or cached RAM). They consume **zero CPU** when idle. |

---

## Use Cases

* **Strangler Fig Pattern:** Replace legacy monolith endpoints one by one.
* **High-Performance Middleware:** Custom authentication, WAF, or transformations.
* **Multi-Tenant Platforms:** Run untrusted user code safely.

---

## Getting Started

To start using Gojinn, you need to compile Caddy with the plugin and create your first WASM function.

* [⚡ Quick Start Guide](./getting-started/quickstart.md)
* [🏗 Architecture Deep Dive](./concepts/architecture.md)
* [📊 Performance Benchmarks](./benchmark.md)
* [🔌 JSON Contract](./concepts/contract.md)