# 🦀 Creating Functions in Rust

Rust and WebAssembly are a perfect combination. Rust offers memory safety and native performance without the overhead of a Garbage Collector, making it ideal for high-density functions on Gojinn.

---

## ✅ Prerequisites

You'll need the `wasm32-wasi` target:

```bash
rustup target add wasm32-wasi
```

### Suggested Dependencies (Cargo.toml)

```toml
[dependencies]
serde = { version = "1.0", features = ["derive"] }
serde_json = "1.0"
```

## 📋 The Pattern (Boilerplate)

```rust
use std::collections::HashMap;
use std::io::{self, Read, Write};
use serde::{Deserialize, Serialize};

// --- 1. Gojinn Contract ---

#[derive(Deserialize)]
struct GojinnRequest {
    method: String,
    headers: HashMap<String, Vec<String>>,
    body: String, // Escaped payload inside here
}

#[derive(Serialize)]
struct GojinnResponse {
    status: u16,
    headers: HashMap<String, Vec<String>>,
    body: String,
}

// --- 2. Your Payload ---

#[derive(Deserialize)]
struct MyPayload {
    name: String,
}

fn main() -> io::Result<()> {
    // A. Read Stdin
    let mut buffer = String::new();
    io::stdin().read_to_string(&mut buffer)?;

    // B. Unwrap Request
    let req: GojinnRequest = match serde_json::from_str(&buffer) {
        Ok(r) => r,
        Err(_) => return reply_error(400, "Invalid JSON Input"),
    };

    // C. Process Internal Payload
    // Note: In Rust, we safely handle empty strings or invalid JSON
    let name = if req.body.is_empty() {
        "Stranger".to_string()
    } else {
        match serde_json::from_str::<MyPayload>(&req.body) {
            Ok(p) => p.name,
            Err(_) => "Stranger".to_string(),
        }
    };

    // D. Logic (Log to Stderr)
    eprintln!("Rust processed request for: {}", name);

    // E. Respond
    let response_body = format!(r#"{{"message": "Hello from Rust, {}!"}}"#, name);
    
    reply(200, response_body)
}

fn reply(status: u16, body: String) -> io::Result<()> {
    let mut headers = HashMap::new();
    headers.insert("Content-Type".to_string(), vec!["application/json".to_string()]);

    let resp = GojinnResponse {
        status,
        headers,
        body,
    };

    let output = serde_json::to_string(&resp)?;
    io::stdout().write_all(output.as_bytes())?;
    Ok(())
}

fn reply_error(status: u16, msg: &str) -> io::Result<()> {
    let body = format!(r#"{{"error": "{}"}}"#, msg);
    reply(status, body)
}
```

## 🔧 Compilation

To compile your Rust function to WASI:

```bash
cargo build --target wasm32-wasi --release
```

The binary will be at `target/wasm32-wasi/release/your_function.wasm`.

---

## 🚀 Why Rust on Gojinn?

- **Memory Stability**: Unlike Go, Rust rarely suffers from initial OOM (Out of Memory) since it doesn't load a heavy runtime

- **Size**: Optimized Rust binaries can be extremely small, ideal for distribution at the Edge