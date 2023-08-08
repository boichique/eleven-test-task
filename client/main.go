package main

import (
	"context"
	"flag"
	"fmt"
	"os"

	"github.com/boichique/eleven_test_task/contracts"
	"golang.org/x/exp/slog"
	"golang.org/x/time/rate"
)

func main() {
	addr := flag.String("addr", "http://localhost:10000", "Address of the service")
	flag.Parse()

	client := NewClient(*addr)

	// get limits from service
	limits, err := client.GetLimits()
	failOnError(err, "getting limits")

	limiter := rate.NewLimiter(rate.Every(limits.Interval), 1)

	// create batch of 1000 items
	batch := client.GenerateTestBatch(1000)

	err = SendBatchParts(context.Background(), client, batch, limits.MaxItems, limiter)
	failOnError(err, "sending batch")

	fmt.Println("Done")
}

func failOnError(err error, message string) {
	if err != nil {
		slog.Error(err.Error(), message)
		os.Exit(1)
	}
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
