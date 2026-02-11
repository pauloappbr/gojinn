package gojinn

import (
	"context"
	"fmt"

	"github.com/dustin/go-humanize"
	"github.com/tetratelabs/wazero"
	"github.com/tetratelabs/wazero/imports/wasi_snapshot_preview1"
)

type EnginePair struct {
	Runtime wazero.Runtime
	Code    wazero.CompiledModule
}

func (r *Gojinn) createWazeroRuntime(wasmBytes []byte) (*EnginePair, error) {
	ctxWazero := context.Background()
	rConfig := wazero.NewRuntimeConfig().WithCloseOnContextDone(true)

	if r.MemoryLimit != "" {
		bytes, err := humanize.ParseBytes(r.MemoryLimit)
		if err == nil && bytes > 0 {
			const wasmPageSize = 65536
			pages := uint32(bytes / wasmPageSize) //nolint:gosec
			if bytes%wasmPageSize != 0 {
				pages++
			}
			rConfig = rConfig.WithMemoryLimitPages(pages)
		}
	}

	engine := wazero.NewRuntimeWithConfig(ctxWazero, rConfig)

	if err := r.buildHostModule(ctxWazero, engine); err != nil {
		return nil, fmt.Errorf("failed to instantiate host module: %w", err)
	}

	wasi_snapshot_preview1.MustInstantiate(ctxWazero, engine)

	code, err := engine.CompileModule(ctxWazero, wasmBytes)
	if err != nil {
		return nil, fmt.Errorf("failed to compile wasm binary: %w", err)
	}

	return &EnginePair{Runtime: engine, Code: code}, nil
}
