# 🔍 Debugging and Troubleshooting

Developing for WebAssembly *in-process* can be challenging because you don't have a traditional debugger attached to the server. This guide lists the most common errors and how to resolve them.


## ❌ Error: `502 Bad Gateway`

This is the most common error. It means Gojinn couldn't get a valid response from your WASM code.

### Common Causes

**1. Panic / Crash**
- Your Go/Rust code exited unexpectedly (exit code != 0)
- **Solution**: Check the Caddy logs (Stderr). Gojinn prints the error there. Avoid using `os.Exit(1)` or `panic()`. Always handle errors and return JSON with status 500.

**2. Out of Memory (OOM)**
- The runtime of your language needed more memory than configured
- **Symptom**: Log `sys_mmap failed` or `failed to instantiate module`
- **Solution**: Increase the `memory_limit` in the Caddyfile. For standard Go, use at least `64MB` or `128MB`

**3. Invalid Protocol**
- Your code printed something to `Stdout` that is not the expected JSON
- **Cause**: Using `fmt.Println("debug")` to debug. This pollutes the response JSON
- **Solution**: Always use `fmt.Fprintf(os.Stderr, ...)` for debug logs. Stdout is exclusive for the response JSON

## ❌ Error: JSON error logs

If you see errors like:

```
json: cannot unmarshal string into Go struct field ... of type int
```

**Cause**: You're violating the [Function Contract](../concepts/contract.md)

- The `status` field must be a number (`200`), not a string (`"200"`)
- The `headers` field must be a map of lists (`{"Key": ["Val"]}`), not simple strings


## ❌ Error: "Operation not supported"

Your code tried to do something prohibited by the WASI Sandbox.

- **Network Attempt**: Trying to connect to a database or call an external API
- **File Attempt**: Trying to read a file outside the allowed context

**Solution**: Gojinn (Phase 1) is focused on pure computation. Pass necessary data via input (JSON) or environment variables. Socket support is planned for the future.


## 🛠️ Recommended Debugging Flow

When something goes wrong, follow this ritual:

### 1. Open two terminals

- **Terminal 1**: Running `./caddy run` (to see logs in real time)
- **Terminal 2**: To run `curl` commands and compile

### 2. Use Stderr Logs

- Fill your code with `eprintln!` (Rust) or `fmt.Fprintf(os.Stderr)` (Go)
- Example:
  ```go
  fmt.Fprintf(os.Stderr, "DEBUG: Payload received: %s\n", string(body))
  ```

### 3. Test in Isolation with Curl

- Don't just test by clicking a button on your website. The browser hides network errors
- Use:
  ```bash
  curl -v -X POST http://localhost:8080/api/function -d '{"test": true}'
  ```