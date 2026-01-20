# ⚡ Performance Benchmarks

> **"Trust, but verify."**

Gojinn's goal is to eliminate the overhead associated with traditional containerization for ephemeral tasks. With the release of **v0.3.0**, we introduced a JIT (Just-In-Time) Caching Engine that makes WebAssembly execution nearly indistinguishable from native code.

Below we present comparative data between **Docker**-based architecture and **Gojinn v0.3.0**.

## 🥊 The Duel: Docker vs. Gojinn

Comparison performed by executing a simple Go calculation function.

- **Docker**: Alpine Linux container running a Go HTTP server (Native Binary)
- **Gojinn**: Binary compiled to `wasm32-wasi` running inside Caddy via Wazero (JIT)

### 📊 Summary of Results

| Metric | Docker Container | Gojinn v0.3.0 | Winner |
|---|---|---|---|
| **Cold Start (First Request)** | ~1,500 ms (Boot OS) | **< 1 ms** (Pre-compiled) | 🏆 **Gojinn** (1500x faster) |
| **Warm Latency (Execution)** | ~4 ms | **~0.3 ms** (Microseconds) | 🏆 **Gojinn** (10x faster) |
| **Artifact Size** | 20.6 MB (Image) | **3.0 MB** (Binary) | 🏆 **Gojinn** (6.8x smaller) |
| **Idle RAM Usage** | ~20 MB (per container) | **~0 MB** (per function)* | 🏆 **Gojinn** |

*\*Note: Gojinn keeps the compiled bytecode in shared memory (KB), but allocates zero active execution memory when idle.*

### 🔬 Detailed Analysis

#### 1️⃣ Cold Start (The "Serverless Killer")
The most critical metric for scaling to zero.
- **Docker**: Needs to create namespace, cgroups, bring up the filesystem, and initialize the kernel/app process.
- **Gojinn**: Since v0.3.0, the module is pre-compiled during Caddy startup (Provisioning). The "Cold Start" for a request is just a memory allocation.

#### 2️⃣ Warm Latency (The "Hot Path")
- **Docker**: Native binary is fast, but network overhead (veth pairs, bridge, NAT) adds latency.
- **Gojinn**: Using **Buffer Pooling** and **JIT Caching**, the code runs in-process. There is no OS network stack overhead between Caddy and the function.

#### 3️⃣ Density
- **Docker**: Limited by RAM. Running 50 containers usually maxes out a small VPS.
- **Gojinn**: Limited by CPU. You can have **thousands** of functions configured. They are just files on disk until triggered.

---

## 🧪 How to Reproduce

Transparency is key. You can run these tests on your own machine using the provided examples.

### ✅ Prerequisites

- Docker installed
- Go installed
- `curl` (terminal utility)
- `xcaddy` (to build the optimized binary)

### Step 1: The Challenger (Docker)

Build and run the Docker image:

```bash
cd benchmark/docker
docker build -t benchmark-go .
docker run --rm -p 8081:8081 benchmark-go
```

### Step 2: The Defender (Gojinn v0.3.0)

Build the optimized Caddy binary and run the example:

```bash
# Build binary with plugin
xcaddy build --with github.com/pauloappbr/gojinn

# Run the Strangler Fig example (contains a heavy calc function)
./caddy run --config examples/legacy-integration/Caddyfile
```

### Step 3: The Benchmark

In another terminal, fire off requests and measure total time.

#### Docker Test

```bash
# Average of 10 requests
for i in {1..10}; do curl -s -w "%{time_total}\n" -o /dev/null http://localhost:8081/api/bench; done
```

#### Gojinn Test

```bash
# Average of 10 requests
for i in {1..10}; do curl -s -w "%{time_total}\n" -o /dev/null http://localhost:8080/api/calc; done
```

---

## 📝 Conclusion

Gojinn does not replace Docker for long-running stateful applications (databases, queues).

However, for ephemeral functions, hooks, validations, and API glue code, Gojinn v0.3.0 offers performance that traditional containerization physically cannot achieve due to OS overhead.