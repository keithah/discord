package connector

import (
	"sync"
	"time"
)

const (
	gatewayReconnectInitialDelay = 5 * time.Second
	gatewayReconnectMaxDelay     = 2 * time.Minute
)

type reconnectThrottle struct {
	mu   sync.Mutex
	next time.Duration
}

func (rt *reconnectThrottle) Next() time.Duration {
	rt.mu.Lock()
	defer rt.mu.Unlock()

	if rt.next == 0 {
		rt.next = gatewayReconnectInitialDelay * 2
		return gatewayReconnectInitialDelay
	}

	delay := rt.next
	rt.next *= 2
	if rt.next > gatewayReconnectMaxDelay {
		rt.next = gatewayReconnectMaxDelay
	}
	return delay
}

func (rt *reconnectThrottle) Reset() {
	rt.mu.Lock()
	rt.next = 0
	rt.mu.Unlock()
}
