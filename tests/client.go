package tests

import (
	"context"
	"fmt"
	"net/http"
	"sync"

	"github.com/boichique/eleven_test_task/contracts"
	"github.com/go-resty/resty/v2"
	"golang.org/x/time/rate"
)

type Client struct {
	client  *resty.Client
	baseURL string
}

func NewClient(url string) *Client {
	hc := &http.Client{}
	rc := resty.NewWithClient(hc)

	return &Client{
		client:  rc,
		baseURL: url,
	}
}

func (c *Client) path(f string, args ...any) string {
	return fmt.Sprintf(c.baseURL+f, args...)
}

func (c *Client) GenerateTestBatch(count int) contracts.Batch {
	var batch contracts.Batch
	var wg sync.WaitGroup
	var mu sync.Mutex

	for i := 0; i < count; i++ {
		i := i
		wg.Add(1)
		go func() {
			defer wg.Done()

			mu.Lock()
			defer mu.Unlock()

			item := contracts.Item{ID: i, Value: fmt.Sprintf("Value %d", i)}
			batch.Items = append(batch.Items, item)
		}()
	}
	wg.Wait()

	return batch
}

func (c *Client) GetLimits() (*contracts.Limits, error) {
	var resp contracts.Limits

	_, err := c.client.R().
		SetResult(&resp).
		Get(c.path("/api/items/limits"))

	return &resp, err
}

func (c *Client) ProcessItems(batch contracts.Batch) error {
	_, err := c.client.R().
		SetHeader("Content-Type", "application/json").
		SetBody(batch).
		Post(c.path("/api/items/process"))

	return err
}

func SendBatchParts(ctx context.Context, client *Client, batch contracts.Batch, maxItems int, limiter *rate.Limiter) error {
	for len(batch.Items) > 0 {
		// Wait for the time interval or until the context is canceled
		err := limiter.Wait(ctx)
		if err != nil {
			return err
		}

		select {
		case <-ctx.Done():
			return ctx.Err() // Context canceled or timed out
		default:
			// Divide the batch into chunks of at most maxItems length
			chunkSize := len(batch.Items)
			if chunkSize > maxItems {
				chunkSize = maxItems
			}

			// Send the part to the service
			err = client.ProcessItems(contracts.Batch{Items: batch.Items[:chunkSize]})
			if err != nil {
				return err
			}

			// Remove the processed items from the batch
			batch.Items = batch.Items[chunkSize:]
		}
	}

	return nil
}
