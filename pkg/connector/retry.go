package connector

import (
	"context"
	"errors"
	"net"
	"net/http"
	"time"

	"github.com/bwmarrin/discordgo"
	"github.com/rs/zerolog"
)

const discordTransientRetryMaxAttempts = 3

var transientRetrySleep = func(ctx context.Context, delay time.Duration) error {
	timer := time.NewTimer(delay)
	defer timer.Stop()
	select {
	case <-ctx.Done():
		return ctx.Err()
	case <-timer.C:
		return nil
	}
}

func executeWithTransientRetry[T any](
	ctx context.Context,
	maxAttempts int,
	retryMsg string,
	fn func() (T, error),
	isTransient func(error) bool,
) (T, error) {
	var zero T
	var result T
	var err error
	for attempt := 1; attempt <= maxAttempts; attempt++ {
		result, err = fn()
		if err == nil || !isTransient(err) || attempt == maxAttempts {
			return result, err
		}
		backoff := time.Duration(attempt) * time.Second
		zerolog.Ctx(ctx).Warn().
			Err(err).
			Int("attempt", attempt).
			Int("max_attempts", maxAttempts).
			Dur("retry_after", backoff).
			Msg(retryMsg)
		if sleepErr := transientRetrySleep(ctx, backoff); sleepErr != nil {
			return zero, sleepErr
		}
	}
	return result, err
}

func isTransientDiscordRESTError(err error) bool {
	if err == nil {
		return false
	}
	if errors.Is(err, context.DeadlineExceeded) {
		return true
	}

	var dnsErr *net.DNSError
	if errors.As(err, &dnsErr) {
		return true
	}

	var opErr *net.OpError
	if errors.As(err, &opErr) && opErr.Op == "dial" {
		return true
	}

	var netErr net.Error
	if errors.As(err, &netErr) && (netErr.Timeout() || netErr.Temporary()) {
		return true
	}

	var restErr *discordgo.RESTError
	if !errors.As(err, &restErr) {
		return false
	}
	if restErr.Response != nil {
		switch restErr.Response.StatusCode {
		case http.StatusBadGateway, http.StatusServiceUnavailable, http.StatusGatewayTimeout:
			return true
		}
	}
	return false
}
