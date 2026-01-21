package gojinn

import (
	"context"
	"fmt"

	"github.com/dustin/go-humanize"
	"github.com/tetratelabs/wazero"
	"github.com/tetratelabs/wazero/api"
	"github.com/tetratelabs/wazero/imports/wasi_snapshot_preview1"
)

// EnginePair mantém o Runtime e o Código Compilado juntos.
type EnginePair struct {
	Runtime wazero.Runtime
	Code    wazero.CompiledModule
}

func (r *Gojinn) createWorker(wasmBytes []byte) (*EnginePair, error) {
	ctxWazero := context.Background()
	rConfig := wazero.NewRuntimeConfig().WithCloseOnContextDone(true)

	if r.MemoryLimit != "" {
		bytes, err := humanize.ParseBytes(r.MemoryLimit)
		if err == nil && bytes > 0 {
			const wasmPageSize = 65536
			pages := uint32(bytes / wasmPageSize)
			if bytes%wasmPageSize != 0 {
				pages++
			}
			rConfig = rConfig.WithMemoryLimitPages(pages)
		}
	}

	engine := wazero.NewRuntimeWithConfig(ctxWazero, rConfig)

	// Host Module
	_, err := engine.NewHostModuleBuilder("gojinn").
		NewFunctionBuilder().
		WithGoModuleFunction(api.GoModuleFunc(func(ctx context.Context, mod api.Module, stack []uint64) {
			level := uint32(stack[0])
			ptr := uint32(stack[1])
			size := uint32(stack[2])
			msgBytes, ok := mod.Memory().Read(ptr, size)
			if !ok {
				return
			}
			msg := string(msgBytes)
			if level == 3 {
				r.logger.Error(msg)
			} else {
				r.logger.Info(msg)
			}
		}), []api.ValueType{api.ValueTypeI32, api.ValueTypeI32, api.ValueTypeI32}, []api.ValueType{}).
		Export("host_log").
		Instantiate(ctxWazero)

	if err != nil {
		return nil, fmt.Errorf("failed to instantiate host module: %w", err)
	}

	wasi_snapshot_preview1.MustInstantiate(ctxWazero, engine)

	code, err := engine.CompileModule(ctxWazero, wasmBytes)
	if err != nil {
		return nil, fmt.Errorf("failed to compile wasm binary: %w", err)
	}

	return &EnginePair{Runtime: engine, Code: code}, nil
}
