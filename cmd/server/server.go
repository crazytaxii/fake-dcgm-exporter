package server

import (
	"cmp"
	"context"
	"fmt"
	"net/http"
	"os"
	"strconv"

	"github.com/spf13/cobra"
	"k8s.io/klog/v2"
	"sigs.k8s.io/controller-runtime/pkg/manager/signals"

	"github.com/crazytaxii/fake-dcgm-exporter/pkg/dcgm"
)

func NewCommand() *cobra.Command {
	cmd := &cobra.Command{
		Short: "fake-dcgm-exporter",
		Long:  "a fake GPU DCGM exporter",
		Run: func(_ *cobra.Command, _ []string) {
			if err := Run(signals.SetupSignalHandler()); err != nil {
				klog.Fatal(err)
			}
		},
	}
	return cmd
}

func Run(ctx context.Context) error {
	hostname, err := os.Hostname()
	if err != nil {
		return fmt.Errorf("failed to get hostname: %w", err)
	}
	cfg, err := LoadConfig(defaultConfigPath, cmp.Or(os.Getenv("NODE_NAME"), hostname), os.Getenv("POD_NAME"))
	if err != nil {
		return fmt.Errorf("failed to load config: %w", err)
	}

	srv, err := NewMetricsServer(cfg)
	if err != nil {
		return err
	}

	errCh := make(chan error)
	go func() {
		errCh <- srv.Start(ctx)
	}()

	select {
	case err := <-errCh:
		return fmt.Errorf("failed to start server: %w", err)
	case <-ctx.Done():
		return srv.Stop(ctx)
	}
}

type MetricsServer struct {
	server          http.Server
	metricsProvider *dcgm.MetricsRenderer
}

func NewMetricsServer(cfg *Config) (*MetricsServer, error) {
	provider, err := dcgm.NewGPUMetricsRenderer(cfg.FakeGPUConfig)
	if err != nil {
		return nil, err
	}
	ms := &MetricsServer{
		server: http.Server{
			Addr: fmt.Sprintf(":%d", cfg.Port),
		},
		metricsProvider: provider,
	}
	http.Handle("/metrics", http.HandlerFunc(ms.handleMetrics))
	http.Handle("/health", http.HandlerFunc(health))
	return ms, nil
}

func (s *MetricsServer) Start(ctx context.Context) error {
	klog.Info("starting server")
	return s.server.ListenAndServe()
}

func (s *MetricsServer) Stop(ctx context.Context) error {
	klog.Info("stopping server")
	return s.server.Shutdown(ctx)
}

func (s *MetricsServer) handleMetrics(w http.ResponseWriter, r *http.Request) {
	compressed := r.Header.Get("Accept-Encoding") == "gzip"
	data, err := s.metricsProvider.Render(compressed)
	if err != nil {
		klog.Errorf("failed to render fake metrics: %v", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Length", strconv.Itoa(len(data)))
	if compressed {
		w.Header().Set("Content-Encoding", "gzip")
	} else {
		w.Header().Set("Content-Type", "text/plain")
	}
	if _, err := w.Write(data); err != nil {
		klog.Errorf("failed to send response: %v", err)
	}
}

func health(w http.ResponseWriter, r *http.Request) {
	_, _ = w.Write([]byte("OK"))
	w.WriteHeader(http.StatusOK)
}
