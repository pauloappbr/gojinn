package gojinn

import (
	"context"
	"fmt"

	"github.com/dustin/go-humanize"
	"github.com/tetratelabs/wazero"
	"github.com/tetratelabs/wazero/api"
	"github.com/tetratelabs/wazero/imports/wasi_snapshot_preview1"
)

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
		NewFunctionBuilder().
		WithGoModuleFunction(api.GoModuleFunc(func(ctx context.Context, mod api.Module, stack []uint64) {
			queryPtr := uint32(stack[0])
			queryLen := uint32(stack[1])
			outPtr := uint32(stack[2])
			outMaxLen := uint32(stack[3])

			qBytes, ok := mod.Memory().Read(queryPtr, queryLen)
			if !ok {
				stack[0] = 0
				return
			}
			query := string(qBytes)

			jsonBytes, err := r.executeQueryToJSON(query)
			if err != nil {
				jsonBytes = []byte(fmt.Sprintf(`[{"error": "%s"}]`, err.Error()))
			}

			bytesToWrite := uint32(len(jsonBytes))

			if bytesToWrite > outMaxLen {
				bytesToWrite = outMaxLen
				jsonBytes = jsonBytes[:bytesToWrite]
			}

			if !mod.Memory().Write(outPtr, jsonBytes) {
				stack[0] = 0
				return
			}

			stack[0] = uint64(bytesToWrite)
		}),
			[]api.ValueType{api.ValueTypeI32, api.ValueTypeI32, api.ValueTypeI32, api.ValueTypeI32},
			[]api.ValueType{api.ValueTypeI32}).
		Export("host_db_query").
		NewFunctionBuilder().
		WithGoModuleFunction(api.GoModuleFunc(func(ctx context.Context, mod api.Module, stack []uint64) {
			keyPtr := uint32(stack[0])
			keyLen := uint32(stack[1])
			valPtr := uint32(stack[2])
			valLen := uint32(stack[3])

			kBytes, ok := mod.Memory().Read(keyPtr, keyLen)
			if !ok {
				return
			}
			key := string(kBytes)

			vBytes, ok := mod.Memory().Read(valPtr, valLen)
			if !ok {
				return
			}
			val := string(vBytes)

			r.kvStore.Store(key, val)

		}), []api.ValueType{api.ValueTypeI32, api.ValueTypeI32, api.ValueTypeI32, api.ValueTypeI32}, []api.ValueType{}).
		Export("host_kv_set").
		NewFunctionBuilder().
		WithGoModuleFunction(api.GoModuleFunc(func(ctx context.Context, mod api.Module, stack []uint64) {
			keyPtr := uint32(stack[0])
			keyLen := uint32(stack[1])
			outPtr := uint32(stack[2])
			outMaxLen := uint32(stack[3])

			kBytes, ok := mod.Memory().Read(keyPtr, keyLen)
			if !ok {
				stack[0] = 0
				return
			}
			key := string(kBytes)

			val, ok := r.kvStore.Load(key)
			if !ok {
				stack[0] = uint64(0xFFFFFFFF)
				return
			}

			valueStr := val.(string)
			valBytes := []byte(valueStr)
			valLen := uint32(len(valBytes))

			bytesToWrite := valLen
			if bytesToWrite > outMaxLen {
				bytesToWrite = outMaxLen
			}

			if !mod.Memory().Write(outPtr, valBytes[:bytesToWrite]) {
				stack[0] = 0
				return
			}

			stack[0] = uint64(bytesToWrite)
		}),
			[]api.ValueType{api.ValueTypeI32, api.ValueTypeI32, api.ValueTypeI32, api.ValueTypeI32},
			[]api.ValueType{api.ValueTypeI32}).
		Export("host_kv_get").
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
