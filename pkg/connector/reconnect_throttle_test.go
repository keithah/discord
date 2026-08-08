package connector

import "testing"

func TestReconnectThrottleProgressionAndReset(t *testing.T) {
	var throttle reconnectThrottle

	if got := throttle.Next(); got != gatewayReconnectInitialDelay {
		t.Fatalf("first delay: got %s, want %s", got, gatewayReconnectInitialDelay)
	}
	if got := throttle.Next(); got != gatewayReconnectInitialDelay*2 {
		t.Fatalf("second delay: got %s, want %s", got, gatewayReconnectInitialDelay*2)
	}
	if got := throttle.Next(); got != gatewayReconnectInitialDelay*4 {
		t.Fatalf("third delay: got %s, want %s", got, gatewayReconnectInitialDelay*4)
	}

	throttle.Reset()
	if got := throttle.Next(); got != gatewayReconnectInitialDelay {
		t.Fatalf("delay after reset: got %s, want %s", got, gatewayReconnectInitialDelay)
	}
}

func TestReconnectThrottleCapsAtMaxDelay(t *testing.T) {
	var throttle reconnectThrottle
	var got = throttle.Next()
	for got < gatewayReconnectMaxDelay {
		got = throttle.Next()
	}
	if got != gatewayReconnectMaxDelay {
		t.Fatalf("cap delay: got %s, want %s", got, gatewayReconnectMaxDelay)
	}
	if got = throttle.Next(); got != gatewayReconnectMaxDelay {
		t.Fatalf("delay past cap: got %s, want %s", got, gatewayReconnectMaxDelay)
	}
}
