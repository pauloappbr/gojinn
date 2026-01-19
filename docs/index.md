# 🧞 Gojinn

> **"The fastest function is the one that never leaves the process."**

[![Go Reference](https://pkg.go.dev/badge/github.com/caddyserver/caddy/v2.svg)](https://pkg.go.dev/github.com/pauloappbr/gojinn)
[![Build Status](https://img.shields.io/badge/build-passing-brightgreen)]()
[![Wasm Engine](https://img.shields.io/badge/engine-wazero-purple)](https://wazero.io)

**Gojinn** is an *in-process serverless* runtime for the [Caddy](https://caddyserver.com) web server.

It allows you to execute **Go**, **Rust**, **Zig**, and **C++** code (compiled to WebAssembly) directly in the HTTP request flow, eliminating network latency and the *cold start* associated with traditional Docker containers or Lambda functions.

---

## 🚀 Why Gojinn?

Traditional serverless and microservices architecture introduces complexity and latency. Gojinn brings computation closer to the data.

| Feature | Description |
| :--- | :--- |
| **🏗️ Zero Infra** | No Docker daemon, no Kubernetes orchestration, no complex sidecars. It's just a Caddy plugin. |
| **⚡ Ultra-Low Latency** | Cold starts of **~1.3ms** (compared to 1500ms+ for Docker). Execution is immediate. |
| **🛡️ Security** | Each request runs in a strict Sandbox via [Wazero](https://wazero.io). If the code fails, the server remains intact. |
| **💰 Zero Idle Cost** | RAM consumption is zero when there are no active requests. No resources allocated waiting for traffic. |

---

## Getting Started

To start using Gojinn, you need to compile Caddy with the plugin and create your first WASM function.

* [Installation and Build](./getting-started/installation.md)
* [Understand the JSON Contract](./concepts/contract.md)
* [Quick Start Guide (Quickstart)](./getting-started/quickstart.md)