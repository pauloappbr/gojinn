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

// --- HELPERS ---

// compileTestWasm é um helper que compila código Go para WASM em tempo de teste.
// Isso garante que estamos testando binários reais compatíveis com a arquitetura atual.
func compileTestWasm(t *testing.T, sourceCode, outName string) string {
	tmpDir := t.TempDir()
	srcPath := filepath.Join(tmpDir, "main.go")
	wasmPath := filepath.Join(tmpDir, outName)

	err := os.WriteFile(srcPath, []byte(sourceCode), 0644)
	if err != nil {
		t.Fatal(err)
	}

	// Compila usando 'go build'
	// GOOS=wasip1 é fundamental para o Go 1.21+ funcionar no Wazero
	cmd := exec.Command("go", "build", "-o", wasmPath, srcPath)
	cmd.Env = append(os.Environ(), "GOOS=wasip1", "GOARCH=wasm")
	out, err := cmd.CombinedOutput()
	if err != nil {
		t.Fatalf("Failed to compile wasm: %v\nOutput: %s", err, out)
	}

	return wasmPath
}

// --- TESTES DE PROVISIONAMENTO (INTEGRAÇÃO) ---

func TestProvision_FullLifecycle(t *testing.T) {
	// 1. Setup: Cria um WASM válido
	code := `package main; func main() {}`
	wasmPath := compileTestWasm(t, code, "lifecycle.wasm")

	// 2. Config: Define o Gojinn com parâmetros explícitos
	r := &Gojinn{
		Path:        wasmPath,
		MemoryLimit: "10MB",
		Timeout:     caddy.Duration(5 * time.Second),
		PoolSize:    2, // Força 2 workers para facilitar a contagem
	}

	// Mock do Contexto do Caddy (necessário para logs e métricas)
	ctx, _ := caddy.NewContext(caddy.Context{Context: context.Background()})

	// 3. Execution: Roda o Provision (que chama setupMetrics e initWorkers internamente)
	err := r.Provision(ctx)

	// 4. Assertions
	assert.NoError(t, err)

	// Verifica se a componentização funcionou:
	// - Metrics (metrics.go)
	assert.NotNil(t, r.metrics, "Metrics struct should be initialized by metrics.go logic")

	// - Worker Pool (worker.go + gojinn.go)
	assert.NotNil(t, r.enginePool, "Worker pool should be initialized")
	assert.Equal(t, 2, len(r.enginePool), "Should have exactly 2 pre-warmed workers")

	// - Worker Integrity (worker.go)
	// Retira um worker do pool para inspecionar
	worker := <-r.enginePool
	assert.NotNil(t, worker.Runtime, "Wazero Runtime should be active")
	assert.NotNil(t, worker.Code, "Compiled Module (JIT) should be present")

	// Devolve para o pool para o Cleanup funcionar
	r.enginePool <- worker

	// 5. Cleanup
	err = r.Cleanup()
	assert.NoError(t, err)
}

func TestProvision_DefaultPoolSize(t *testing.T) {
	// Testa a lógica de auto-scaling no gojinn.go
	code := `package main; func main() {}`
	wasmPath := compileTestWasm(t, code, "autoscaling.wasm")

	r := &Gojinn{
		Path:     wasmPath,
		PoolSize: 0, // Deve acionar a lógica: max(NumCPU*4, 50)
	}

	ctx, _ := caddy.NewContext(caddy.Context{Context: context.Background()})
	err := r.Provision(ctx)
	assert.NoError(t, err)

	// Garante que criou workers (mínimo de 50 pela regra atual)
	assert.Greater(t, len(r.enginePool), 0)

	_ = r.Cleanup()
}

func TestProvision_FileNotFound(t *testing.T) {
	// Testa validação básica no gojinn.go
	r := &Gojinn{
		Path: "./arquivo_fantasma.wasm",
	}

	ctx, _ := caddy.NewContext(caddy.Context{Context: context.Background()})

	err := r.Provision(ctx)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "failed to read wasm file")
}

// Teste de resiliência: Configuração de memória inválida não deve quebrar o boot
func TestProvision_GracefulInvalidConfig(t *testing.T) {
	code := `package main; func main() {}`
	wasmPath := compileTestWasm(t, code, "graceful.wasm")

	r := &Gojinn{
		Path:        wasmPath,
		MemoryLimit: "BATATA", // Valor inválido
		PoolSize:    1,
	}

	ctx, _ := caddy.NewContext(caddy.Context{Context: context.Background()})

	// Não deve retornar erro, apenas ignora o limite e loga (se tivéssemos capturando logs)
	err := r.Provision(ctx)
	assert.NoError(t, err)
	assert.Equal(t, 1, len(r.enginePool))

	_ = r.Cleanup()
}
