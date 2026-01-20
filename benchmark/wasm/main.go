package main

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"
)

type GojinnRequest struct {
	Method string `json:"method"`
	Body   string `json:"body"`
}

type GojinnResponse struct {
	Status  int               `json:"status"`
	Headers map[string]string `json:"headers"`
	Body    string            `json:"body"`
}

type Order struct {
	ID    string  `json:"id"`
	Value float64 `json:"value"`
}

func main() {
	var req GojinnRequest
	if err := json.NewDecoder(os.Stdin).Decode(&req); err != nil {
		sendError(400, "Invalid Input Protocol")
		return
	}

	var order Order
	if err := json.NewDecoder(strings.NewReader(req.Body)).Decode(&order); err != nil {
		sendError(400, "Invalid JSON Body")
		return
	}

	tax := order.Value * 0.15
	total := order.Value + tax

	outputPayload := fmt.Sprintf(`{"order_id": "%s", "total_final": %.2f, "engine": "gojinn-v0.3.0"}`, order.ID, total)

	resp := GojinnResponse{
		Status: 200,
		Headers: map[string]string{
			"Content-Type": "application/json",
		},
		Body: outputPayload,
	}

	json.NewEncoder(os.Stdout).Encode(resp)
}

func sendError(status int, msg string) {
	resp := GojinnResponse{
		Status: status,
		Body:   fmt.Sprintf(`{"error": "%s"}`, msg),
	}
	json.NewEncoder(os.Stdout).Encode(resp)
}
