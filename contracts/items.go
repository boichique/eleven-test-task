package contracts

import "time"

type Item struct {
	ID    int    `json:"id"`
	Value string `json:"value"`
}

type Batch struct {
	Items []Item `json:"items"`
}

type Limits struct {
	MaxItems int           `json:"maxItems"`
	Interval time.Duration `json:"interval"`
}

type ProcessRequest struct {
	Batch Batch `json:"batch"`
}
