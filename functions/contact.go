package main

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
)

type Request struct {
	Method string `json:"method"`
	Body   string `json:"body"`
}

type Response struct {
	Status  int                 `json:"status"`
	Headers map[string][]string `json:"headers"`
	Body    string              `json:"body"`
}

func main() {
	fmt.Fprintf(os.Stderr, "DEBUG [WASM]: Iniciando a funcao de contato...\n")
	fmt.Fprintf(os.Stderr, "DEBUG [WASM]: Verificando variaveis de ambiente...\n")

	inputData, _ := io.ReadAll(os.Stdin)

	fmt.Fprintf(os.Stderr, "DEBUG [WASM]: Recebi payload de tamanho: %d bytes\n", len(inputData))

	var req Request
	json.Unmarshal(inputData, &req)

	resp := Response{
		Status: 200,
		Headers: map[string][]string{
			"Content-Type": {"application/json"},
			"X-Powered-By": {"Gojinn WASM"},
		},
		Body: `{"message": "Sucesso! O Debug Mode funciona!", "original_method": "` + req.Method + `"}`,
	}

	json.NewEncoder(os.Stdout).Encode(resp)
}
