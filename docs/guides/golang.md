# 🐹 Creating Functions in Go

Go is the "native" language of the Cloud Native ecosystem and an excellent choice for Gojinn. Due to the nature of WebAssembly (WASI), writing functions for Gojinn is very similar to writing command-line tools (CLI).

---

## 📋 The Pattern (Boilerplate)

To avoid deserialization errors, we recommend copying and maintaining these base structures in your functions.

```go
package main

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
)

// --- 1. Gojinn Structures (The Contract) ---

// Input Wrapper
type GojinnRequest struct {
	Method  string              `json:"method"`
	Headers map[string][]string `json:"headers"`
	Body    string              `json:"body"` // User payload comes here as a string
}

// Output Wrapper
type GojinnResponse struct {
	Status  int                 `json:"status"`
	Headers map[string][]string `json:"headers"`
	Body    string              `json:"body"`
}

// --- 2. Your Data Structures ---

type MyUserPayload struct {
	Name  string `json:"name"`
	Email string `json:"email"`
}

func main() {
	// A. Read Input (Stdin)
	inputBytes, err := io.ReadAll(os.Stdin)
	if err != nil {
		replyError(500, "Failed to read stdin")
		return
	}

	// B. Unwrap Gojinn Request
	var req GojinnRequest
	if err := json.Unmarshal(inputBytes, &req); err != nil {
		replyError(400, "Invalid input format")
		return
	}

	// C. Process your Payload (which was inside the Body string)
	var payload MyUserPayload
	// Tip: If body is empty, handle that!
	if req.Body != "" {
		json.Unmarshal([]byte(req.Body), &payload)
	}

	// D. Business Logic (Logs go to Stderr)
	fmt.Fprintf(os.Stderr, "Processing user: %s\n", payload.Name)

	// E. Respond
	responseData := map[string]string{"message": "Processed successfully"}
	responseJSON, _ := json.Marshal(responseData)

	reply(200, string(responseJSON))
}

// Helpers
func reply(status int, body string) {
	resp := GojinnResponse{
		Status: status,
		Headers: map[string][]string{
			"Content-Type": {"application/json"},
		},
		Body: body,
	}
	json.NewEncoder(os.Stdout).Encode(resp)
}

func replyError(status int, msg string) {
	// Creates a simple error JSON
	errJSON := fmt.Sprintf(`{"error": "%s"}`, msg)
	reply(status, errJSON)
}
```

## 🔧 Compilation

You have two main options for compiling your Go code to WASM.

### Option 1: Standard Compiler (Go Toolchain)

Full compatibility, but larger binaries (~2MB+) and higher memory usage.

```bash
GOOS=wasip1 GOARCH=wasm go build -o function.wasm main.go
```

**Gojinn requirement**: Configure `memory_limit` of 64MB or higher in the Caddyfile.

### Option 2: TinyGo (Recommended for Performance)

Produces tiny binaries (~100KB - 500KB) and ultra-fast startup.

```bash
tinygo build -o function.wasm -target=wasi main.go
```

- **Advantage**: Allows much lower `memory_limit` (e.g., 10MB)
- **Limitation**: Some standard libraries (like full `net/http` or heavy reflection) may not work perfectly

---

## 💡 Golden Tips

### Logging

Use `fmt.Fprintf(os.Stderr, ...)` for logs. Gojinn redirects WASM Stderr to Caddy logs. 

> ⚠️ **Never use** `fmt.Println` for logs, as it goes to Stdout and corrupts the JSON response.

### Panic

Avoid `panic()`. If your code panics, Gojinn will return error 502. Handle errors and return JSON with status 400/500.