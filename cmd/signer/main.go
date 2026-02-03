package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/pauloappbr/gojinn/pkg/sovereign"
)

func main() {
	action := flag.String("action", "", "gen-keys | sign")
	name := flag.String("name", "gojinn", "Nome da chave (para gen-keys)")
	keyFile := flag.String("key", "", "Caminho da chave privada (para sign)")
	wasmFile := flag.String("file", "", "Arquivo wasm para assinar")
	flag.Parse()

	switch *action {
	case "gen-keys":
		err := sovereign.GenerateKeys(*name)
		if err != nil {
			panic(err)
		}
		fmt.Printf("✅ Chaves geradas: %s.pub e %s.priv\n", *name, *name)
		pubBytes, _ := os.ReadFile(*name + ".pub")
		fmt.Printf("📋 PUBLIC KEY (Copie para o Caddyfile):\n%s\n", string(pubBytes))

	case "sign":
		if *keyFile == "" || *wasmFile == "" {
			panic("Precisa de --key e --file")
		}

		// Ler chave privada
		privHex, err := os.ReadFile(*keyFile)
		if err != nil {
			panic(fmt.Errorf("falha ao ler chave privada: %w", err))
		}
		privKeyBytes, err := sovereign.ParsePrivateKey(string(privHex))
		if err != nil {
			panic(err)
		}

		// Ler Wasm (CORREÇÃO AQUI: Tratamento de Erro)
		wasmBytes, err := os.ReadFile(*wasmFile)
		if err != nil {
			panic(fmt.Errorf("ARQUIVO WASM NÃO ENCONTRADO ou ilegível: %w", err))
		}

		if len(wasmBytes) == 0 {
			panic("O arquivo WASM está vazio! Verifique o build.")
		}

		// Assinar
		signedBytes, err := sovereign.SignWasm(wasmBytes, privKeyBytes)
		if err != nil {
			panic(err)
		}

		// Sobrescrever
		err = os.WriteFile(*wasmFile, signedBytes, 0644)
		if err != nil {
			panic(err)
		}
		fmt.Printf("🔐 Arquivo assinado com sucesso: %s (Tamanho: %d bytes)\n", *wasmFile, len(signedBytes))

	default:
		fmt.Println("Use: --action=gen-keys ou --action=sign")
	}
}
