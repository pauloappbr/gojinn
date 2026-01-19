# ⚡ Quickstart: Hello World

In this guide, we'll create, compile, and run your first serverless function on Gojinn in less than 5 minutes.

We'll create a simple function that receives a name via JSON and returns a greeting.

---

## 1. Function Code (Go)

Create a file called `main.go`. 

This code includes the necessary structures to read Gojinn input and format the output correctly.

```go
package main

import (
	"encoding/json"
	"io"
	"os"
)

// --- 1. Gojinn Contract Structures ---

// Gojinn "wraps" the original HTTP request in this JSON
type GojinnRequest struct {
	Body string `json:"body"` // The actual payload comes here as a string
}

// Gojinn expects EXACTLY this response format
type GojinnResponse struct {
	Status  int                 `json:"status"`
	Headers map[string][]string `json:"headers"`
	Body    string              `json:"body"`
}

// --- 2. Your Custom Payload ---
type MyInput struct {
	Name string `json:"name"`
}

func main() {
	// A. Read STDIN (Input from Caddy)
	inputBytes, _ := io.ReadAll(os.Stdin)

	// B. Unwrap the Request
	var gojinnReq GojinnRequest
	json.Unmarshal(inputBytes, &gojinnReq)

	// C. Read the User Payload (which was inside the Body)
	var myInput MyInput
	json.Unmarshal([]byte(gojinnReq.Body), &myInput)

	// D. Business Logic
	message := "Hello, " + myInput.Name + "! Welcome to Gojinn."
	if myInput.Name == "" {
		message = "Hello, Stranger! Send a name in the JSON."
	}

	// E. Prepare Response (Must be JSON serialized as a string)
	responseJSON, _ := json.Marshal(map[string]string{
		"message": message,
	})

	// F. Send to STDOUT (Response to Caddy)
	reply(200, string(responseJSON))
}

// Helper to format the final response
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
```

## 2. 🔧 Compile to WebAssembly

Gojinn uses the WASI (WebAssembly System Interface). Let's compile our Go code to `.wasm`.

```bash
GOOS=wasip1 GOARCH=wasm go build -o hello.wasm main.go
```

> 💡 **Tip**: The generated file will be ~2MB. If you want tiny binaries (100KB), consider using [TinyGo](https://tinygo.org/). The code above is compatible with both.

## 3. 🐳 Configure Caddy

Create a file called `Caddyfile` in the same directory:

```caddy
{
    # Define the plugin execution order
    order gojinn before file_server
}

:8080 {
    # API route
    handle /api/hello {
        header Content-Type application/json
        
        # Run our WASM binary
        gojinn ./hello.wasm {
            timeout 1s
            memory_limit 64MB
            env ENV "dev"
        }
    }
}
```

## 4. 🚀 Run and Test

Start the server:

```bash
./caddy run
```

In another terminal, make a request:

```bash
curl -X POST http://localhost:8080/api/hello \
     -H "Content-Type: application/json" \
     -d '{"name": "Paulo"}'
```

### ✅ Expected Result

```json
{"message":"Hello, Paulo! Welcome to Gojinn."}
```

---

## 🎉 Congratulations!

You just executed Go code in-process inside your web server, without Docker and without containers.

## 📚 Next Steps

- [Understand the details of the JSON Contract](../concepts/contract.md)
- [See the Complete Caddyfile Reference](https://caddyserver.com/docs/caddyfile)