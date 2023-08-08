package tests

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/boichique/eleven_test_task/contracts"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"golang.org/x/time/rate"
)

func TestClient_GetLimits(t *testing.T) {
	// Create a test server that returns a mock response
	mockResponse := contracts.Limits{
		MaxItems: 100,
		Interval: 2 * time.Minute,
	}
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		err := json.NewEncoder(w).Encode(mockResponse)
		require.NoError(t, err)
	}))
	defer ts.Close()

	// Create a new client with the test server's URL
	client := NewClient(ts.URL)

	// Call GetLimits and check the response
	limits, err := client.GetLimits()
	require.NoError(t, err)
	assert.Equal(t, mockResponse, *limits)
}

func TestClient_ProcessItems(t *testing.T) {
	// Create a test server that returns a mock response
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
	}))
	defer ts.Close()

	// Create a new client with the test server's URL
	client := NewClient(ts.URL)

	// Create a test batch of items
	batch := client.GenerateTestBatch(1000)

	// Call ProcessItems and check the response
	err := client.ProcessItems(batch)
	require.NoError(t, err)
}

func TestRateLimiting(t *testing.T) {
	// Create a test server that returns mock responses
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
	}))
	defer ts.Close()

	// Create a new client with the test server's URL
	client := NewClient(ts.URL)

	// Create a test batch of items
	batch := client.GenerateTestBatch(300)
	limits := &contracts.Limits{
		MaxItems: 100,
		Interval: 2 * time.Second,
	}

	limiter := rate.NewLimiter(rate.Every(limits.Interval), 1)

	// Use a context with a timeout of 7 seconds
	ctx, cancel := context.WithTimeout(context.Background(), 7*time.Second)
	defer cancel()

	// Call SendBatchParts and check if the requests are spaced correctly
	err := SendBatchParts(ctx, client, batch, limits.MaxItems, limiter)
	require.NoError(t, err)
}

func TestBatchSizeLimit(t *testing.T) {
	// Create a test server that returns mock responses
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
	}))
	defer ts.Close()

	// Create a new client with the test server's URL
	client := NewClient(ts.URL)

	// Create a test batch of items with size greater than the maxItems limit
	batch := client.GenerateTestBatch(300)
	limits := &contracts.Limits{
		MaxItems: 100,
		Interval: 1 * time.Second,
	}

	limiter := rate.NewLimiter(rate.Every(limits.Interval), 1)

	// Call sendBatchParts and check if the batch is divided correctly
	err := SendBatchParts(context.Background(), client, batch, limits.MaxItems, limiter)
	require.NoError(t, err)
}
