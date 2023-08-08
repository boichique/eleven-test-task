package main

import (
	"fmt"
	"sync"

	"github.com/boichique/eleven_test_task/contracts"
)

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
