package main

import (
	"encoding/json"
	"fmt"
	"os"
)

type RequestPayload struct {
	Method  string              `json:"method"`
	URI     string              `json:"uri"`
	Headers map[string][]string `json:"headers"`
	Body    string              `json:"body"`
}

type ResponsePayload struct {
	Status  int                 `json:"status"`
	Headers map[string][]string `json:"headers"`
	Body    string              `json:"body"`
}

func main() {
	var req RequestPayload
	if err := json.NewDecoder(os.Stdin).Decode(&req); err != nil {
		sendError(400, "Invalid JSON input: "+err.Error())
		return
	}

	auth := req.Headers["Authorization"]
	if len(auth) == 0 || auth[0] != "secret" {
		sendError(401, "Unauthorized: Missing or wrong secret")
		return
	}

	resp := ResponsePayload{
		Status: 200,
		Headers: map[string][]string{
			"Content-Type": {"application/json"},
			"X-Gojinn":     {"Phase-1"},
		},
		Body: fmt.Sprintf(`{"message": "Hello from Wasm!", "your_method": "%s", "your_path": "%s"}`, req.Method, req.URI),
	}

	json.NewEncoder(os.Stdout).Encode(resp)
}

func sendError(code int, msg string) {
	resp := ResponsePayload{
		Status: code,
		Body:   fmt.Sprintf(`{"error": "%s"}`, msg),
	}
	json.NewEncoder(os.Stdout).Encode(resp)
}
