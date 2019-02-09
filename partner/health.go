package partner

import (
	"time"

	"github.com/taask/taask-server/service"
)

type healthChecker struct {
	IsHealthy bool
	stream    service.PartnerService_StreamUpdatesClient
}

func newHealthChecker(stream service.PartnerService_StreamUpdatesClient) *healthChecker {
	return &healthChecker{
		IsHealthy: true,
		stream:    stream,
	}
}

func (hc *healthChecker) startHealthChecking() {
	for {
		<-time.After(time.Duration(time.Second * 30))

		if err := hc.stream.Send(&service.UpdateRequest{IsHealthCheck: true}); err != nil {
			hc.IsHealthy = false
			break
		}
	}
}
