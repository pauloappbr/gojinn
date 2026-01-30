# Frequently Asked Questions (FAQ)

This document answers common questions about Gojinn’s design, scope, and philosophy.

---

## ❓ Why not just use Docker or Kubernetes?

Because they solve a different problem.

Docker and Kubernetes manage **processes and infrastructure**.
Gojinn executes **untrusted code inside an existing process** with strict isolation.

If you need:
- container orchestration
- service meshes
- pod scheduling

Gojinn is not the right tool.

If you need:
- deterministic execution
- ultra-low latency
- zero idle cost
- sovereign infrastructure

Gojinn is purpose-built for that.

---

## ❓ Is Gojinn a Kubernetes replacement?

No.

Gojinn is not:
- an orchestrator
- a scheduler
- a control plane

Gojinn executes **code + events**, nothing more.

---

## ❓ Why WASM instead of native plugins?

Because WASM provides:

- Memory safety
- Deterministic execution
- Language neutrality
- Strong sandboxing

Native plugins require trust.
Gojinn assumes **zero trust**.

---

## ❓ Can a function keep state between requests?

Not implicitly.

All execution is ephemeral by default.

State is only allowed through:
- Explicit host-managed KV
- Explicit database bindings

This is intentional and enforced by design.

---

## ❓ Does Gojinn support background jobs or queues?

Not by default.

Async and event-driven execution may exist as **explicit, opt-in host capabilities**.
There is no implicit background execution.

---

## ❓ Why is Gojinn so opinionated?

Because flexibility without constraints leads to:

- undefined behavior
- insecure defaults
- unmaintainable APIs

Gojinn prefers **few correct ways** over many incorrect ones.

---

## ❓ Why not support feature X from AWS Lambda / Cloudflare Workers?

Gojinn does not aim for cloud parity.

If a feature:
- increases abstraction
- hides execution cost
- introduces implicit state
- weakens determinism

It will not be accepted, regardless of popularity.

---

## ❓ Is Gojinn production-ready?

Yes, for its intended scope.

Gojinn is suitable for:
- edge logic
- API endpoints
- request-bound computation
- stateful serverless with explicit storage

It is not suitable for:
- batch processing
- long-running jobs
- uncontrolled concurrency

---

## ❓ Can Gojinn be used commercially?

Yes.

Gojinn is licensed under **Apache 2.0**.
You may use it commercially, modify it, and embed it in proprietary systems.

---

## ❓ Who controls the project?

The Project Founder.

Gojinn follows a BDFL governance model to preserve architectural integrity and long-term vision.

See: GOVERNANCE.md
