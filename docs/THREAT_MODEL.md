# Gojinn - Threat Model & Security Architecture
*Version: 1.0.0 | Methodology: STRIDE*

This document outlines the security posture, trust boundaries, and mitigation strategies implemented in the Gojinn Sovereign Cloud Platform.

## 1. System Architecture & Trust Boundaries
Gojinn operates as a multi-tenant, edge-ready WASM execution environment. The core trust boundaries are:
- **Public Network -> Caddy Edge:** Handled via TLS and Rate Limiting.
- **Tenant Context -> Worker Execution:** Handled via Strict WASM Sandboxing & Context Cancellation.
- **Worker Execution -> Host OS:** Handled via Wasmtime/Wazero memory isolation (No host syscalls allowed unless explicitly granted via Host Functions).
- **Node -> Node (Cluster):** Handled via NATS JetStream native authentication and encrypted subjects.

## 2. STRIDE Threat Analysis

### S - Spoofing (Identity Deception)
* **Threat:** A malicious actor attempts to impersonate a legitimate tenant or administrator to execute workloads or read private state.
* **Mitigation:** Gojinn enforces mandatory API Keys (`X-API-Key` or `Authorization: Bearer`) per tenant. At the cluster level, nodes communicate via pre-shared TLS certificates and Nkey/Seed cryptography. Internal API endpoints (`/_sys/`) require localized or strongly authenticated access.

### T - Tampering (Data Modification)
* **Threat:** An attacker tries to alter worker outputs, modify the distributed state, or corrupt the physical disk holding the KV store.
* **Mitigation:** - **At-Rest:** All NATS JetStream and KV data is natively encrypted on disk using AES-GCM (Phase 27). 
    - **In-Transit:** Cluster replication uses TLS.
    - **Auditability:** Every tenant execution generates an HMAC-SHA256 Signed Audit Log injected immutably into the tenant's isolated KV bucket, guaranteeing Tamper-Evident tracking (Phase 28).

### R - Repudiation (Denying Actions)
* **Threat:** A tenant denies executing a transaction that altered distributed state.
* **Mitigation:** The combination of NATS JetStream WAL (Write-Ahead Logging) and the Cryptographically Signed Audit Logs ensures non-repudiation. The `job_id`, `timestamp`, and `tenant_id` are cryptographically bound to the execution output.

### I - Information Disclosure (Data Leaks)
* **Threat:** Tenant A gains access to Tenant B's execution queue, environment variables, or KV state.
* **Mitigation:** Hard Multi-Tenant Isolation (Phase 28). Gojinn provisions distinct physical streams (`WORKER_{TENANT}`) and KV Buckets (`STATE_{TENANT}`) per tenant. Worker pools are dynamically allocated and strictly bound to specific tenant subjects (`gojinn.tenant.{id}.>`). Shared memory spaces do not exist.

### D - Denial of Service (DoS)
* **Threat:** A tenant uploads an infinite loop (`for {}`) or attempts to allocate massive amounts of RAM, crashing the host server (OOM).
* **Mitigation:** - **Edge:** Caddy implements strict rate limiting (`rate_limit`, `rate_burst`) per IP/Tenant.
    - **Runtime:** CPU executions are bound by `context.WithTimeout`.
    - **I/O & Memory:** A rigorous `cappedWriter` immediately kills the WASM execution context via `context.CancelFunc` if output exceeds the predefined `MaxOutputBytes` (e.g., 5MB limit).

### E - Elevation of Privilege
* **Threat:** A WASM guest module attempts a sandbox escape to execute shell commands (`/bin/sh`) or read host files (`/etc/passwd`).
* **Mitigation:** Gojinn utilizes the `wazero` runtime. WASM modules execute in a default-deny capability model. Guests have zero access to network sockets, environmental host variables, or host filesystems unless strictly mounted via the Caddyfile `mounts` directive (and even then, constrained by WASI virtual filesystems).

## 3. Supply Chain Security (SLSA Alignment)
Gojinn guarantees build integrity through:
- **Reproducible Builds:** Compiled using `-trimpath` and `-buildvcs=false`.
- **SBOM:** CycloneDX JSON generated via Syft.
- **Transparency Logs:** SHA-256 Checksums provided for all release binaries.