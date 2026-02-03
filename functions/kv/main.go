package main

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
	"strings"
	"unsafe"
)

//go:wasmimport gojinn host_kv_set
func host_kv_set(keyPtr, keyLen, valPtr, valLen uint32)

//go:wasmimport gojinn host_kv_get
func host_kv_get(keyPtr, keyLen, outPtr, outLimit uint32) uint64

type Request struct {
	Body string `json:"body"`
}

func main() {
	inputBytes, _ := io.ReadAll(os.Stdin)
	var req Request
	if err := json.Unmarshal(inputBytes, &req); err != nil {
		req.Body = string(inputBytes)
	}

	cmd := strings.TrimSpace(req.Body)
	parts := strings.Split(cmd, " ")

	if len(parts) < 2 {
		// Retorna 400 Bad Request
		fmt.Printf(`{"status": 400, "body": "Usage: SET <key> <val> | GET <key>"}`)
		return
	}

	action := strings.ToUpper(parts[0])
	key := parts[1]

	if action == "SET" && len(parts) >= 3 {
		val := strings.Join(parts[2:], " ")
		ptrKey, lenKey := stringToPtr(key)
		ptrVal, lenVal := stringToPtr(val)

		host_kv_set(ptrKey, lenKey, ptrVal, lenVal)

		// SUCESSO: Status 200 (Inteiro)
		fmt.Printf(`{"status": 200, "headers": {"Content-Type": "text/plain"}, "body": "✅ SET OK: %s = %s"}`, key, val)

	} else if action == "GET" {
		buf := make([]byte, 1024)
		ptrKey, lenKey := stringToPtr(key)
		ptrOut := uint32(uintptr(unsafe.Pointer(&buf[0])))

		written := host_kv_get(ptrKey, lenKey, ptrOut, 1024)

		if written == 0xFFFFFFFFFFFFFFFF {
			// NÃO ENCONTRADO: Status 404 (Inteiro)
			fmt.Printf(`{"status": 404, "body": "Key not found: %s"}`, key)
		} else {
			val := string(buf[:written])
			// SUCESSO: Status 200 (Inteiro)
			fmt.Printf(`{"status": 200, "headers": {"Content-Type": "text/plain"}, "body": "Found: %s"}`, val)
		}
	} else {
		fmt.Printf(`{"status": 400, "body": "Unknown command"}`)
	}
}

func stringToPtr(s string) (uint32, uint32) {
	if len(s) == 0 {
		return 0, 0
	}
	return uint32(uintptr(unsafe.Pointer(unsafe.StringData(s)))), uint32(len(s))
}
