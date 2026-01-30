# 🎨 Gojinn — Design Principles & Invariants

This document defines the design philosophy, core invariants, and explicit non-goals of Gojinn.

It exists to answer why certain decisions were made, what must never change, and where the project deliberately draws hard boundaries.

If a proposal or pull request violates a principle in this document, it is out of scope by design, not by lack of interest.

## 1. Design Goals

Gojinn is designed to be a minimal, deterministic, in-process serverless runtime.

Its primary goals are:

- Zero network hop execution
- Deterministic performance and resource usage
- Strong isolation without containers
- Self-hosted and sovereign by default
- WASI-first, language-agnostic execution
- Operational simplicity over feature breadth

Gojinn intentionally optimizes for request-time execution, not infrastructure orchestration.

## 2. Core Runtime Invariants

The following invariants define Gojinn’s identity.
They are non-negotiable and must hold true across all future versions.

### 2.1 In-Process Execution

All user code executes inside the Caddy process.

- No sidecar containers
- No external worker daemons
- No internal service mesh
- No implicit IPC or networking

This invariant exists to guarantee:

- Microsecond latency
- Zero internal network overhead
- Predictable execution cost

Any design that introduces a network hop violates this invariant.

### 2.2 Request-Bound Execution Model

All execution is strictly bound to an HTTP request lifecycle.

- No background jobs
- No long-running workers
- No async task queues
- No detached execution contexts

When the request ends:

- Execution ends
- Memory is reclaimed
- State must already be persisted externally (DB / KV)

This ensures:

- Deterministic cleanup
- No hidden resource consumption
- Clear failure semantics

### 2.3 Explicit Capabilities Only

WASM modules have no ambient authority.

By default, modules have access to:

- Stdin / Stdout
- Explicit host-managed capabilities only

They do not have access to:

- File system
- Network
- Environment variables
- Host OS APIs

Any capability must be:

- Explicitly exposed
- Versioned
- Auditable

This follows the Principle of Least Authority (POLA).

### 2.4 Deterministic Resource Limits

CPU, memory, and execution time must be enforceable and predictable.

- Hard memory limits
- Instruction / execution budgets
- Immediate termination on violation

There are no “best effort” limits.

If a module exceeds its budget:

- It is terminated
- The request fails safely
- The host process remains unaffected

This invariant exists to prevent noisy neighbors and abuse.

### 2.5 Runtime Simplicity Over Abstraction

The runtime favors simplicity and clarity over general-purpose abstractions.

- Fewer knobs
- Fewer layers
- Fewer magic behaviors

Complexity is pushed out of the runtime and into:

- The user’s WASM code
- External databases
- Explicit configuration

This keeps the hot path fast and auditable.

## 3. Non-Goals (What Gojinn Will Never Be)

Gojinn explicitly rejects the following domains:

### 3.1 Not a Container Platform

Gojinn will never:

- Run Docker containers
- Emulate OCI semantics
- Replace Kubernetes workloads

Containers solve infrastructure isolation.
Gojinn solves request-time computation.

### 3.2 Not a Distributed Orchestrator

Gojinn will not implement:

- Multi-node scheduling
- Cluster consensus
- Leader election
- Autoscaling controllers

Node-level orchestration is delegated to:

- Systemd
- Nomad
- Kubernetes (if desired)

Gojinn operates inside a single process boundary.

### 3.3 Not a Background Job System

Gojinn will not support:

- Cron jobs
- Queues
- Event consumers
- Long-running workers

If you need background execution, use:

- A job runner
- A queue system
- A dedicated service

### 3.4 Not a General-Purpose PaaS

Gojinn is not:

- A replacement for Cloudflare Workers
- A hosted platform
- A managed runtime

It is a self-hosted runtime component, not a platform business.

## 4. API Stability Boundaries

Gojinn distinguishes between stable contracts and internal implementation details.

### Stable Contracts

The following are considered stable and versioned:

- WASM stdin/stdout JSON contract
- Official SDK APIs
- Caddy configuration surface (documented directives)

Breaking changes here require:

- Major version bump
- Clear migration path

### Unstable / Internal APIs

The following are not public APIs:

- Runtime internals
- Worker pool implementation
- Memory layout
- Host function wiring
- Scheduling strategies

These may change at any time to improve:

- Performance
- Security
- Maintainability

## 5. Why Gojinn Is Not “Kubernetes 2.0”

Kubernetes optimizes for:

- Long-running services
- Networked components
- Horizontal scaling via replicas

Gojinn optimizes for:

- Per-request execution
- In-memory communication
- Vertical efficiency

| Dimension | Kubernetes | Gojinn |
|---|---|---|
| Execution | Long-lived Pods | Request-bound |
| Communication | Network | Memory |
| Isolation | OS-level | WASM sandbox |
| Scaling | Horizontal | CPU-bound |
| Complexity | High | Minimal |

Attempting to merge these domains results in:

- Higher latency
- More abstraction
- Loss of determinism

Gojinn explicitly avoids this path.

## 6. Governance by Design

Gojinn is governed first by design constraints, not feature requests.

Decisions prioritize:

- Performance over flexibility
- Determinism over convenience
- Explicitness over magic
- Runtime simplicity over ecosystem breadth

Features that violate these principles are rejected even if they are popular.

## 7. Final Note

Gojinn is intentionally opinionated.

Its power comes not from doing everything —
but from doing one thing extremely well.

> “A runtime is only fast if it knows what it refuses to be.”