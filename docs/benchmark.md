# ⚡ Performance Benchmarks

> **"Trust, but verify."**

Gojinn's goal is to eliminate the overhead associated with traditional containerization for ephemeral tasks. Below we present comparative data between **Docker**-based architecture and **Gojinn**.

## 🥊 The Duel: Docker vs. Gojinn

Comparison performed by executing a simple Go calculation function.

- **Docker**: Alpine Linux container running a Go HTTP server
- **Gojinn**: Binary compiled to `wasm32-wasi` running on Caddy via Wazero

### 📊 Summary of Results

| Metric | Docker Container | Gojinn | Winner |
|---|---|---|---|
| **Cold Start (Initialization)** | ~1,500 ms | **~1 ms** | 🏆 **Gojinn** (1500x faster) |
| **Artifact Size** | 20.6 MB (Image) | **3.0 MB** (Binary) | 🏆 **Gojinn** (6.8x smaller) |
| **RAM Usage (Idle)** | ~10-20 MB (Daemon) | **0 MB** | 🏆 **Gojinn** |
| **Latency (Warm)** | ~9 ms | ~13 ms | 🤝 Technical Tie |

### 🔬 Detailed Analysis

#### 1️⃣ Cold Start (The "Serverless Killer")

The most critical metric.

- **Docker**: Needs to create namespace, cgroups, bring up the filesystem, and initialize the kernel/app process
- **Gojinn**: Only needs to allocate a block of memory and instantiate the VM. It's virtually instantaneous

#### 2️⃣ Density and Cost

- **Docker**: Each running container consumes operating system resources, even without receiving traffic
- **Gojinn**: Idle code is just a file on disk. You can have thousands of functions configured without consuming 1 byte of RAM until they're called


## 🧪 How to Reproduce

Transparency is key. You can run these tests on your own machine using the `benchmark/` directory from the repository.

### ✅ Prerequisites

- Docker installed
- Go installed
- `curl` and `time` (terminal utilities)

### Step 1: The Challenger (Docker)

Build and run the Docker image:

```bash
cd benchmark/docker
docker build -t benchmark-go .
docker run --rm -p 8081:8081 benchmark-go
```

### Step 2: The Defender (Gojinn)

Compile the WASM and start Caddy:

```bash
cd benchmark/wasm
GOOS=wasip1 GOARCH=wasm go build -ldflags="-s -w" -o tax.wasm main.go
cd ../..
./caddy run --config benchmark/Caddyfile.test
```

### Step 3: The Stress Test

In another terminal, fire off requests and measure total time (connection + processing).

**Docker Test:**

```bash
time curl -s -X POST http://localhost:8081/api/bench -d '{"id":"1", "valor":100}'
```

**Gojinn Test:**

```bash
time curl -s -X POST http://localhost:8080/api/bench -d '{"id":"1", "valor":100}'
```

## 📝 Conclusion

Gojinn **does not replace** Docker for long-running applications (databases, queues, complex backends).

However, for ephemeral functions, hooks, validations, and serverless logic, Gojinn offers startup performance and resource efficiency that traditional containerization **physically cannot achieve**.