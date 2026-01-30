# Gojinn SDK 🧞‍♂️

The official SDK for building Serverless functions in Go that run inside Gojinn.

## 📦 Installation

```bash
go get github.com/pauloappbr/gojinn/sdk
```

## 🚀 Features

### 1. Input and Output (I/O)

```go
package main

import "github.com/pauloappbr/gojinn/sdk"

func main() {
    // Reads the request body, headers and method
    req, _ := sdk.Parse()
    
    sdk.Log("Received method: %s", req.Method)

    // Send JSON response
    sdk.SendJSON(map[string]string{"hello": "world"})
    
    // OR Send HTML response (For HTMX)
    // sdk.SendHTML("<h1>Hello World</h1>")
}
```

### 2. Database (SQL)

Gojinn uses the Host's (Caddy) connection pool. You don't need to open a connection, just run the query. Supports: **Postgres**, **MySQL**, and **SQLite** (depending on your Caddyfile configuration).

```go
func main() {
    // Returns []map[string]interface{}
    rows, err := sdk.DB.Query("SELECT id, name FROM users WHERE active = 1")
    
    if err != nil {
        sdk.SendError(500, err.Error())
        return
    }
    
    sdk.SendJSON(rows)
}
```

### 3. Key-Value Store (In-Memory)

Ultra-fast in-memory storage on the server's RAM. Shared across all executions. Great for counters and caching.

```go
func main() {
    // Store
    sdk.KV.Set("last_access", "2026-01-30")

    // Retrieve
    val, found := sdk.KV.Get("last_access")
    if found {
        sdk.Log("Retrieved value: %s", val)
    }
}
```

### 4. Logs and Debug

Use `sdk.Log` instead of `fmt.Println`. If the request has the `X-Gojinn-Debug` header with the correct password, these logs will appear in the HTTP response header.

```go
sdk.Log("Starting complex processing...")
```
