package partner

import (
	"time"
)

type healthChecker struct {
	IsHealthy     bool
	UnhealthyChan chan error
}

func newHealthChecker() *healthChecker {
	return &healthChecker{
		IsHealthy:     true,
		UnhealthyChan: make(chan error),
	}
}

func (hc *healthChecker) startHealthCheckingWithClient(stream PartnerService_StreamUpdatesClient) {
	for {
		<-time.After(time.Duration(time.Second * 20))

		if err := stream.Send(&UpdateRequest{IsHealthCheck: true}); err != nil {
			hc.IsHealthy = false
			hc.UnhealthyChan <- err
			break
		}

		hc.IsHealthy = true
	}
}

func (hc *healthChecker) startHealthCheckingWithServer(stream PartnerService_StreamUpdatesServer) {
	for {
		<-time.After(time.Duration(time.Second * 30))

		if err := stream.Send(&UpdateRequest{IsHealthCheck: true}); err != nil {
			hc.IsHealthy = false
			hc.UnhealthyChan <- err
			break
		}

		hc.IsHealthy = true
	}
}
