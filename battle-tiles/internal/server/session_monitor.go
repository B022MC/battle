package server

import (
	"battle-tiles/internal/service"
	"context"

	"github.com/go-kratos/kratos/v2/transport"
)

// MonitorServer adapts SessionMonitor to kratos transport.Server lifecycle
type MonitorServer struct {
	monitor *service.SessionMonitor
}

func NewMonitorServer(m *service.SessionMonitor) transport.Server {
	return &MonitorServer{monitor: m}
}

func (s *MonitorServer) Start(ctx context.Context) error {
	s.monitor.Start(ctx)
	return nil
}

func (s *MonitorServer) Stop(ctx context.Context) error {
	s.monitor.Stop()
	return nil
}
