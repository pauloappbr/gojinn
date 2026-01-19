# 📦 Execution Environment (Sandbox)

Gojinn executes your code within a strict **WebAssembly Sandbox**. This means the environment where your code runs is completely isolated from the host operating system and the main Caddy process.

Understanding the limitations and resources of this environment is crucial to avoid common errors like `exit code 1` or `502 Bad Gateway`.

---

## 🧠 Memory Management

Unlike a normal process that requests memory from the OS as needed, WebAssembly operates in pre-allocated or expandable linear memory up to a limit.

### Why does my Go code give a memory error?

Languages with *Garbage Collector* (like standard Go, Java, C#) need to allocate a significant amount of memory right at startup to bring up the runtime.

- **TinyGo / Rust / C++**: Consume very little (KB or a few MB). Work well with low limits.
- **Go (Standard)**: The compiled binary is usually 2MB+ and the runtime needs heap space.
  - **Recommendation**: If using standard Go, set `memory_limit` to at least **64MB**
  - **Error symptom**: If the limit is too low, Caddy will return **502 Bad Gateway** and logs will show something like `sys_mmap failed` or `out of memory`

---

## ⏱️ Timeouts and CPU Cycles

In-process *serverless* shares the CPU with the web server. An infinite loop in your WASM code could theoretically lock up a Caddy thread.

To prevent this, Gojinn enforces a **Hard Timeout** via `context.Context`.

1. When the time set in `timeout` is reached, the current WASM instruction is interrupted
2. Memory is released immediately
3. The client receives an error

### Best Practices

- Keep your functions short and fast
- For long tasks, use the function to dispatch an event to a queue (work in progress) and return `202 Accepted`

---

## 🔒 File System and Network

By default, the WASI specification (WebAssembly System Interface) and Gojinn's implementation adopt a **"Deny by Default"** stance.

### File System (Filesystem)

Your WASM code **cannot see** the server's disk.

- Attempting to open `/etc/passwd` or `./config.json` will result in `file not found` or `permission denied` error
- **Solution**: Pass configuration via environment variables (`env`) or through the JSON input payload

### Network (Network)

Currently, Gojinn **does not allow** arbitrary TCP/UDP socket opening (e.g., connecting to a MySQL database or calling an external API via HTTP) directly from within WASM.

- **Reason**: Security and limitations of the current WASI spec (preview1)
- **Current State**: Code should be pure (computational)
- **Roadmap**: Support for *WASI Sockets* is planned to allow controlled external HTTP calls

---

## 📋 Environment Variables

Host environment variables **are not automatically inherited**. This prevents accidental leakage of server secrets to the function.

You must explicitly allow each variable in the `Caddyfile`:

```caddy
gojinn app.wasm {
    # The function will only see "ENV" and "API_KEY". 
    # It will NOT see "PATH", "HOME", etc.
    env ENV "production"
    env API_KEY {env.MY_HOST_SECRET}
}
```