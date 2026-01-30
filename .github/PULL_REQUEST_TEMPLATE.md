# Pull Request

Thank you for contributing to **Gojinn** 🧞  
Before submitting, please ensure this PR aligns with the project’s design principles.

---

## 📌 What does this PR do?

<!--
Provide a clear and concise description of the change.
What problem does it solve? What behavior does it introduce or modify?
-->

---

## 🎯 Motivation

<!--
Why is this change necessary?
Link to an Issue or describe the real-world use case.
-->

Fixes # (issue)

---

## 🧠 Design Alignment

Please confirm that this change respects Gojinn’s core design invariants
(as defined in DESIGN.md):

- [ ] In-process execution (no new network hops)
- [ ] Request-bound execution model
- [ ] Explicit capabilities only (no ambient authority)
- [ ] Deterministic resource usage
- [ ] Runtime simplicity over abstraction

If any invariant is affected, **explain why** and what trade-offs were considered.

---

## 🔒 Security Considerations

- [ ] This change does **not** introduce new implicit permissions
- [ ] Any new capability is explicit, scoped, and documented
- [ ] Failure modes are safe (panic isolation, proper cleanup)

If applicable, describe potential security implications.

---

## ⚙️ Performance Impact

- [ ] No impact on the hot path
- [ ] Performance impact is negligible
- [ ] Performance impact is measurable (include benchmarks)

If relevant, include:
- Before / After benchmarks
- Allocation or latency changes

---

## 🧪 Testing

Describe how this change was tested:

- [ ] Unit tests
- [ ] Integration tests
- [ ] Manual testing
- [ ] Not tested (explain why)

Provide relevant details or commands used.

---

## 📚 Documentation

- [ ] Documentation updated (docs / comments)
- [ ] No documentation changes required

If docs were updated, list affected files.

---

## 🚫 Out of Scope Check

This PR does **not** attempt to introduce:

- [ ] Background jobs / async workers
- [ ] Distributed orchestration logic
- [ ] Container or Kubernetes abstractions
- [ ] Implicit networking or IPC
- [ ] Platform-specific or vendor-locked APIs

---

## ✅ Checklist

- [ ] Code builds successfully (`make all`)
- [ ] Tests pass (`make test`)
- [ ] Commit messages follow Conventional Commits
- [ ] PR is scoped and focused (no unrelated changes)

---

## 💬 Additional Context

<!--
Anything else reviewers should know?
Trade-offs, future follow-ups, or known limitations.
-->
