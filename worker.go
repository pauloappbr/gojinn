package gojinn

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/nats-io/nats.go"
	"github.com/tetratelabs/wazero"
	"go.uber.org/zap"
)

const (
	MaxRetries = 5
)

func (r *Gojinn) startWorkerSubscriber(id int, topic string, wasmBytes []byte) (*nats.Subscription, error) {
	pair, err := r.createWazeroRuntime(wasmBytes)
	if err != nil {
		return nil, fmt.Errorf("failed to create wazero runtime for worker %d: %w", id, err)
	}

	queueGroup := fmt.Sprintf("WORKERS_%s", hashString(r.Path))

	sub, err := r.js.QueueSubscribe(topic, queueGroup, func(m *nats.Msg) {
		meta, err := m.Metadata()
		if err != nil {
			r.logger.Error("Failed to get msg metadata", zap.Error(err))
			_ = m.Nak()
			return
		}

		deliverCount := meta.NumDelivered
		r.logger.Debug("Worker picked up job",
			zap.Int("worker_id", id),
			zap.Uint64("seq", meta.Sequence.Stream),
			zap.Uint64("attempt", deliverCount))

		_ = m.InProgress()

		ctx, cancel := context.WithTimeout(context.Background(), time.Duration(r.Timeout))
		defer cancel()

		stdoutBuf := bufferPool.Get().(*bytes.Buffer)
		stdoutBuf.Reset()
		defer bufferPool.Put(stdoutBuf)

		stderrBuf := bufferPool.Get().(*bytes.Buffer)
		stderrBuf.Reset()
		defer bufferPool.Put(stderrBuf)

		fsConfig := wazero.NewFSConfig()
		for host, guest := range r.Mounts {
			fsConfig = fsConfig.WithDirMount(host, guest)
		}

		modConfig := wazero.NewModuleConfig().
			WithStdout(stdoutBuf).
			WithStderr(stderrBuf).
			WithStdin(bytes.NewReader(m.Data)).
			WithSysWalltime().
			WithSysNanotime().
			WithFSConfig(fsConfig)

		for k, v := range r.Env {
			modConfig = modConfig.WithEnv(k, v)
		}

		mod, err := pair.Runtime.InstantiateModule(ctx, pair.Code, modConfig)
		if err != nil {
			errMsg := fmt.Sprintf("Wasm Error: %v | Stderr: %s", err, stderrBuf.String())
			r.logger.Warn("Job Failed",
				zap.Int("worker_id", id),
				zap.Uint64("attempt", deliverCount),
				zap.String("error", errMsg))

			if deliverCount >= MaxRetries {
				r.logger.Error("Max Retries Reached. Moving to DLQ/CrashDump.", zap.Uint64("seq", meta.Sequence.Stream))

				snapshot := CrashSnapshot{
					Timestamp: time.Now(),
					Error:     errMsg,
					Input:     json.RawMessage(m.Data),
					Env:       r.Env,
					WasmFile:  r.Path,
				}
				dumpBytes, _ := json.MarshalIndent(snapshot, "", "  ")
				filename := fmt.Sprintf("crash_%s_seq%d.json", time.Now().Format("20060102-150405"), meta.Sequence.Stream)
				r.saveCrashDump(filename, dumpBytes)

				if ackErr := m.Ack(); ackErr != nil {
					r.logger.Error("Failed to Ack Poisoned Message", zap.Error(ackErr))
				}
				return
			}

			backoff := time.Duration(deliverCount) * time.Second
			_ = m.NakWithDelay(backoff)
			return
		}

		mod.Close(ctx)

		if err := m.Ack(); err != nil {
			r.logger.Error("Failed to ACK message", zap.Error(err))
		} else {
			r.logger.Info("Job Completed Successfully",
				zap.Int("worker_id", id),
				zap.Uint64("seq", meta.Sequence.Stream))
		}

	},
		nats.ManualAck(),
		nats.BindStream("GOJINN_WORKER"),
		nats.MaxDeliver(MaxRetries+1),
	)

	if err != nil {
		return nil, err
	}

	r.logger.Info("Worker Durable Consumer Ready",
		zap.Int("id", id),
		zap.String("topic", topic),
		zap.String("queue", queueGroup))

	return sub, nil
}
