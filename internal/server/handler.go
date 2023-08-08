package server

import (
	"errors"
	"net/http"
	"time"

	"github.com/boichique/eleven_test_task/contracts"
	"github.com/labstack/echo/v4"
)

var (
	ErrBlocked = errors.New("blocked")
	maxItems   = 100
	interval   = 5 * time.Second
)

type Handler struct{}

func NewHandler() *Handler {
	return &Handler{}
}

func (h *Handler) ProcessBatch(c echo.Context) error {
	var req contracts.ProcessRequest

	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, echo.Map{"error": "invalid request"})
	}

	if len(req.Batch.Items) > maxItems {
		return c.JSON(http.StatusTooManyRequests, ErrBlocked)
	}

	// Process the batch of items here...

	return c.NoContent(http.StatusOK)
}

func (h *Handler) GetLimits(c echo.Context) error {
	limits := contracts.Limits{
		MaxItems: maxItems,
		Interval: interval,
	}

	return c.JSON(http.StatusOK, limits)
}
