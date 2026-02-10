package gojinn

import (
	"fmt"
	"strings"
	"time"

	"github.com/nats-io/nats-server/v2/server"
	"github.com/nats-io/nats.go"
	"go.uber.org/zap"
)

func (g *Gojinn) startEmbeddedNATS() error {
	opts := &server.Options{
		Port:   g.NatsPort,
		NoLog:  true,
		NoSigs: true,
	}

	if len(g.NatsRoutes) > 0 {
		opts.Routes = server.RoutesFromStr(strings.Join(g.NatsRoutes, ","))
	}

	ns, err := server.NewServer(opts)
	if err != nil {
		return fmt.Errorf("failed to create NATS server: %w", err)
	}
	g.natsServer = ns

	go ns.Start()

	if !ns.ReadyForConnections(10 * time.Second) {
		return fmt.Errorf("nats server failed to start")
	}

	clientURL := ns.ClientURL()
	g.logger.Info("Embedded NATS Started", zap.String("url", clientURL))

	nc, err := nats.Connect(clientURL)
	if err != nil {
		return fmt.Errorf("failed to connect to local NATS: %w", err)
	}
	g.natsConn = nc

	return nil
}

func (g *Gojinn) ReloadWorkers() error {
	g.logger.Info("Hot Reload Initiated: Recycling Workers...")

	g.subsMu.Lock()
	defer g.subsMu.Unlock()

	for _, sub := range g.subs {
		if err := sub.Drain(); err != nil {
			g.logger.Warn("Failed to drain worker sub", zap.Error(err))
		}
	}
	g.subs = nil

	wasmBytes, err := g.loadWasmSecurely(g.Path)
	if err != nil {
		return fmt.Errorf("failed to reload wasm file: %w", err)
	}

	topic := g.getFunctionTopic()

	for i := 0; i < g.PoolSize; i++ {
		sub, err := g.startWorkerSubscriber(i, topic, wasmBytes)
		if err != nil {
			return fmt.Errorf("failed to start new worker %d: %w", i, err)
		}
		g.subs = append(g.subs, sub)
	}

	g.logger.Info("Hot Reload Complete", zap.Int("new_workers", len(g.subs)))
	return nil
}

func (g *Gojinn) getFunctionTopic() string {
	return fmt.Sprintf("gojinn.exec.%s", hashString(g.Path))
}
