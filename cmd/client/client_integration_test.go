//go:build integration
// +build integration

package client

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// TestClient_Integration_FullFlow tests the complete flow of connecting to a
// meter, polling data, and verifying the output against a mock HTTP server.
func TestClient_Integration_FullFlow(t *testing.T) {
	// Set up a mock server that simulates a power meter endpoint
	handler := http.NewServeMux()
	handler.HandleFunc("/status", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(`{
			"power": 1234.5,
			"voltage": 230.1,
			"current": 5.36,
			"frequency": 50.0,
			"energy": 987.6
		}`))
	})

	server := httptest.NewServer(handler)
	defer server.Close()

	client, err := New(Config{
		Host:     server.URL,
		Interval: 1 * time.Second,
		Timeout:  5 * time.Second,
	})
	require.NoError(t, err)
	require.NotNil(t, client)

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	// Run the client and collect at least one reading
	readings := make(chan Reading, 10)
	errCh := make(chan error, 1)

	go func() {
		errCh <- client.Run(ctx, readings)
	}()

	var received []Reading
	collect:
	for {
		select {
		case r := <-readings:
			received = append(received, r)
			if len(received) >= 2 {
				cancel()
				break collect
			}
		case err := <-errCh:
			// context cancellation is expected
			if err != nil && err != context.Canceled && err != context.DeadlineExceeded {
				t.Fatalf("unexpected error from client: %v", err)
			}
			break collect
		case <-ctx.Done():
			break collect
		}
	}

	require.NotEmpty(t, received, "expected at least one reading")

	// Validate the contents of the first reading
	first := received[0]
	assert.InDelta(t, 1234.5, first.Power, 0.01)
	assert.InDelta(t, 230.1, first.Voltage, 0.01)
	assert.InDelta(t, 5.36, first.Current, 0.01)
	assert.InDelta(t, 50.0, first.Frequency, 0.01)
	assert.InDelta(t, 987.6, first.Energy, 0.01)
}

// TestClient_Integration_ServerUnavailable verifies that the client handles
// a server that becomes unavailable gracefully without panicking.
func TestClient_Integration_ServerUnavailable(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusServiceUnavailable)
	}))
	defer server.Close()

	client, err := New(Config{
		Host:     server.URL,
		Interval: 500 * time.Millisecond,
		Timeout:  2 * time.Second,
	})
	require.NoError(t, err)

	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	readings := make(chan Reading, 10)
	errCh := make(chan error, 1)

	go func() {
		errCh <- client.Run(ctx, readings)
	}()

	select {
	case err := <-errCh:
		if err != nil && err != context.Canceled && err != context.DeadlineExceeded {
			// Errors from unavailable server are acceptable
			t.Logf("client returned expected error: %v", err)
		}
	case <-ctx.Done():
		// timeout is also acceptable
	}
}
