package main

import (
	"github.com/boichique/eleven_test_task/contracts"
)

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
