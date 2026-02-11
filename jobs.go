package gojinn

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/nats-io/nats.go"
	"go.uber.org/zap"
)

func (r *Gojinn) runBackgroundJob(wasmFile string) {
	cronPayload := `{"event_type": "cron", "source": "gojinn_scheduler"}`
	r.runAsyncJob(wasmFile, cronPayload)
}

func (r *Gojinn) runAsyncJob(wasmFile, payload string) {
	if r.js == nil {
		r.logger.Error("Cannot queue async job: JetStream not ready", zap.String("file", wasmFile))
		return
	}

	topic := fmt.Sprintf("gojinn.exec.%s", hashString(wasmFile))

	jobPayload := struct {
		Method  string              `json:"method"`
		URI     string              `json:"uri"`
		Headers map[string][]string `json:"headers"`
		Body    string              `json:"body"`
	}{
		Method: "ASYNC",
		URI:    "internal://async/job",
		Headers: map[string][]string{
			"X-Source": {"internal"},
		},
		Body: payload,
	}

	jobBytes, err := json.Marshal(jobPayload)
	if err != nil {
		r.logger.Error("Failed to marshal async job", zap.Error(err))
		return
	}

	msgID := fmt.Sprintf("job_%d", time.Now().UnixNano())

	_, err = r.js.Publish(topic, jobBytes, nats.MsgId(msgID))
	if err != nil {
		r.logger.Error("Failed to persist async job",
			zap.String("file", wasmFile),
			zap.Error(err))
		return
	}

	r.logger.Info("Async Job Persisted & Queued",
		zap.String("file", wasmFile),
		zap.String("msg_id", msgID))
}
