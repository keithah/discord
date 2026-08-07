package connector

import (
	"context"
	"errors"
	"net"
	"net/http"
	"net/url"
	"testing"
	"time"

	"github.com/bwmarrin/discordgo"
)

func TestIsTransientDiscordRESTError(t *testing.T) {
	tests := []struct {
		name string
		err  error
		want bool
	}{
		{
			name: "service unavailable",
			err: &discordgo.RESTError{
				Response: &http.Response{StatusCode: http.StatusServiceUnavailable},
			},
			want: true,
		},
		{
			name: "bad gateway",
			err: &discordgo.RESTError{
				Response: &http.Response{StatusCode: http.StatusBadGateway},
			},
			want: true,
		},
		{
			name: "gateway timeout",
			err: &discordgo.RESTError{
				Response: &http.Response{StatusCode: http.StatusGatewayTimeout},
			},
			want: true,
		},
		{
			name: "bad request is permanent",
			err: &discordgo.RESTError{
				Response: &http.Response{StatusCode: http.StatusBadRequest},
				Message:  &discordgo.APIErrorMessage{Code: discordgo.ErrCodeInvalidFormBody},
			},
		},
		{
			name: "plain error is permanent",
			err:  errors.New("nope"),
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			if got := isTransientDiscordRESTError(tc.err); got != tc.want {
				t.Fatalf("unexpected transient classification: got %t, want %t", got, tc.want)
			}
		})
	}
}

func TestIsTransientDiscordRESTErrorCoversNetworkFailures(t *testing.T) {
	err := &url.Error{
		Op:  "Post",
		URL: "https://discord.com/api/v10/channels/123/messages",
		Err: &net.OpError{Op: "dial", Err: timeoutErr{}},
	}
	if !isTransientDiscordRESTError(err) {
		t.Fatal("expected dial timeout to be classified as transient")
	}
}

func TestExecuteWithTransientRetryRetriesTransientFailures(t *testing.T) {
	originalSleep := transientRetrySleep
	transientRetrySleep = func(context.Context, time.Duration) error { return nil }
	t.Cleanup(func() {
		transientRetrySleep = originalSleep
	})

	attempts := 0
	got, err := executeWithTransientRetry(context.Background(), 3, "test op", func() (string, error) {
		attempts++
		if attempts < 3 {
			return "", &discordgo.RESTError{Response: &http.Response{StatusCode: http.StatusServiceUnavailable}}
		}
		return "ok", nil
	}, isTransientDiscordRESTError)
	if err != nil {
		t.Fatalf("expected success, got error: %v", err)
	}
	if got != "ok" {
		t.Fatalf("unexpected result: %q", got)
	}
	if attempts != 3 {
		t.Fatalf("unexpected attempts: got %d, want 3", attempts)
	}
}

func TestExecuteWithTransientRetryStopsOnPermanentFailure(t *testing.T) {
	originalSleep := transientRetrySleep
	transientRetrySleep = func(context.Context, time.Duration) error { return nil }
	t.Cleanup(func() {
		transientRetrySleep = originalSleep
	})

	attempts := 0
	_, err := executeWithTransientRetry(context.Background(), 3, "test op", func() (string, error) {
		attempts++
		return "", errors.New("permanent failure")
	}, isTransientDiscordRESTError)
	if err == nil {
		t.Fatal("expected error")
	}
	if attempts != 1 {
		t.Fatalf("unexpected attempts: got %d, want 1", attempts)
	}
}
