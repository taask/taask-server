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

func (hc *healthChecker) startHealthChecking(syncer syncSource) {
	for {
		<-time.After(time.Duration(time.Second * 20))

		if err := syncer.Send(&UpdateRequest{IsHealthCheck: true}); err != nil {
			hc.IsHealthy = false
			hc.UnhealthyChan <- err
			break
		}
	}
}
