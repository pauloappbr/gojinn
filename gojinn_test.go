package gojinn

import (
	"context"
	"os"
	"os/exec"
	"path/filepath"
	"testing"
	"time"

	"github.com/caddyserver/caddy/v2"
	"github.com/stretchr/testify/assert"
)

// Helper para compilar um WASM simples para o teste
func compileTestWasm(t *testing.T, sourceCode, outName string) string {
	tmpDir := t.TempDir()
	srcPath := filepath.Join(tmpDir, "main.go")
	wasmPath := filepath.Join(tmpDir, outName)

	err := os.WriteFile(srcPath, []byte(sourceCode), 0644)
	if err != nil {
		t.Fatal(err)
	}

	// Compila usando go build (exige Go instalado no ambiente de teste)
	cmd := exec.Command("go", "build", "-o", wasmPath, srcPath)
	cmd.Env = append(os.Environ(), "GOOS=wasip1", "GOARCH=wasm")
	out, err := cmd.CombinedOutput()
	if err != nil {
		t.Fatalf("Failed to compile wasm: %v\nOutput: %s", err, out)
	}

	return wasmPath
}

func TestProvision_ValidatesConfig(t *testing.T) {
	// Código Go mínimo
	code := `package main; func main() {}`
	wasmPath := compileTestWasm(t, code, "empty.wasm")

	r := &Gojinn{
		Path:        wasmPath,
		MemoryLimit: "10MB",
		Timeout:     caddy.Duration(5 * time.Second),
	}

	// Mock do Contexto do Caddy
	ctx, _ := caddy.NewContext(caddy.Context{Context: context.Background()})

	// Deve passar sem erro
	err := r.Provision(ctx)
	assert.NoError(t, err)
	assert.NotNil(t, r.engine)
}

func TestProvision_InvalidMemoryLimit(t *testing.T) {
	code := `package main; func main() {}`
	wasmPath := compileTestWasm(t, code, "empty.wasm")

	r := &Gojinn{
		Path:        wasmPath,
		MemoryLimit: "BATATA", // Valor inválido
	}

	ctx, _ := caddy.NewContext(caddy.Context{Context: context.Background()})

	// Deve falhar no provision
	err := r.Provision(ctx)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "invalid memory_limit")
}

func TestProvision_FileNotFound(t *testing.T) {
	r := &Gojinn{
		Path: "./arquivo_que_nao_existe.wasm",
	}

	ctx, _ := caddy.NewContext(caddy.Context{Context: context.Background()})

	err := r.Provision(ctx)
	assert.Error(t, err)
}

// Nota: Testar o ServeHTTP completo exige mockar http.ResponseWriter e *http.Request
// Para a Fase 0, garantir que o Provision valida as configs de segurança é o mais crítico.
