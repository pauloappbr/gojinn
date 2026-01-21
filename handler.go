package gojinn

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"sync"
	"time"

	"github.com/caddyserver/caddy/v2/modules/caddyhttp"
	"github.com/tetratelabs/wazero"
	"go.uber.org/zap"
)

var bufferPool = sync.Pool{
	New: func() interface{} {
		return new(bytes.Buffer)
	},
}

func (r *Gojinn) ServeHTTP(rw http.ResponseWriter, req *http.Request, next caddyhttp.Handler) error {
	start := time.Now()
	r.metrics.active.WithLabelValues(r.Path).Inc()
	defer r.metrics.active.WithLabelValues(r.Path).Dec()

	// 1. Acquire Worker
	pair := <-r.enginePool
	defer func() { r.enginePool <- pair }()

	ctx, cancel := context.WithTimeout(req.Context(), time.Duration(r.Timeout))
	defer cancel()

	// 2. Prepare Input
	bodyBytes, err := io.ReadAll(req.Body)
	if err != nil {
		return err
	}
	traceID := req.Header.Get("traceparent")

	reqPayload := struct {
		Method  string              `json:"method"`
		URI     string              `json:"uri"`
		Headers map[string][]string `json:"headers"`
		Body    string              `json:"body"`
		TraceID string              `json:"trace_id,omitempty"`
	}{
		Method:  req.Method,
		URI:     req.RequestURI,
		Headers: req.Header,
		Body:    string(bodyBytes),
		TraceID: traceID,
	}

	inputJSON, err := json.Marshal(reqPayload)
	if err != nil {
		return caddyhttp.Error(http.StatusInternalServerError, err)
	}

	stdoutBuf := bufferPool.Get().(*bytes.Buffer)
	stdoutBuf.Reset()
	defer bufferPool.Put(stdoutBuf)

	config := wazero.NewModuleConfig().
		WithStdout(stdoutBuf).
		WithStderr(os.Stderr).
		WithStdin(bytes.NewReader(inputJSON)).
		WithArgs(r.Args...)

	for k, v := range r.Env {
		config = config.WithEnv(k, v)
	}

	// 3. Fast Instantiation
	instance, err := pair.Runtime.InstantiateModule(ctx, pair.Code, config)

	duration := time.Since(start).Seconds()
	statusLabel := "200"

	if err != nil {
		statusLabel = "500"
		if ctx.Err() == context.DeadlineExceeded {
			statusLabel = "504"
		}
		r.metrics.duration.WithLabelValues(r.Path, statusLabel).Observe(duration)
		r.logger.Error("wasm execution failed", zap.Error(err))
		return caddyhttp.Error(http.StatusInternalServerError, err)
	}
	defer instance.Close(ctx)

	// 4. Process Output
	if stdoutBuf.Len() == 0 {
		statusLabel = "500"
		r.metrics.duration.WithLabelValues(r.Path, statusLabel).Observe(duration)
		return caddyhttp.Error(http.StatusInternalServerError, fmt.Errorf("empty response"))
	}

	var respPayload struct {
		Status  int                 `json:"status"`
		Headers map[string][]string `json:"headers"`
		Body    string              `json:"body"`
	}

	if err := json.Unmarshal(stdoutBuf.Bytes(), &respPayload); err != nil {
		statusLabel = "502"
		r.metrics.duration.WithLabelValues(r.Path, statusLabel).Observe(duration)
		return caddyhttp.Error(http.StatusBadGateway, fmt.Errorf("invalid json output"))
	}

	if respPayload.Status == 0 {
		respPayload.Status = 200
	}
	statusLabel = fmt.Sprintf("%d", respPayload.Status)
	r.metrics.duration.WithLabelValues(r.Path, statusLabel).Observe(duration)

	for k, v := range respPayload.Headers {
		for _, val := range v {
			rw.Header().Add(k, val)
		}
	}
	rw.WriteHeader(respPayload.Status)
	rw.Write([]byte(respPayload.Body))

	return nil
}
