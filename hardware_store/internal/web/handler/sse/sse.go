package sse

import (
	"encoding/json"
	"fmt"
	"hardware_store/internal/kafka"
	"hardware_store/internal/model/product"
	"io"
	"log/slog"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

type SSEHandler struct {
	consumer *kafka.Consumer
	logger   *slog.Logger
}

func NewSSEHandler(consumer *kafka.Consumer, logger *slog.Logger) *SSEHandler {
	return &SSEHandler{
		consumer: consumer,
		logger:   logger,
	}
}

func (h *SSEHandler) HandleSSE(c *gin.Context) {
	c.Writer.Header().Set("Content-Type", "text/event-stream")
	c.Writer.Header().Set("Cache-Control", "no-cache")
	c.Writer.Header().Set("Connection", "keep-alive")
	c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
	c.Writer.Header().Set("X-Accel-Buffering", "no")
	c.Writer.WriteHeader(http.StatusOK)

	fmt.Fprint(c.Writer, ": connected\n\n")
	c.Writer.Flush()

	ch := make(chan product.ProductUpdate, 10)

	h.consumer.GetHub().RegisterSSE(ch)
	h.logger.Info("📻 SSE client connected", slog.String("ip", c.ClientIP()))

	defer func() {
		h.consumer.GetHub().UnregisterSSE(ch)
		h.logger.Info("📻 SSE client disconnected", slog.String("ip", c.ClientIP()))
	}()

	ping := time.NewTicker(10 * time.Second)
	defer ping.Stop()

	c.Stream(func(w io.Writer) bool {
		select {
		case event, ok := <-ch:
			if !ok {
				return false
			}
			data, _ := json.Marshal(event)
			fmt.Fprintf(c.Writer, "event: product_update\n")
			fmt.Fprintf(c.Writer, "data: %s\n\n", data)
			c.Writer.Flush()
			return true
		case <-ping.C:
			fmt.Fprintf(c.Writer, "event: ping\ndata: %d\n\n", time.Now().Unix())
			c.Writer.Flush()
			return true
		case <-c.Request.Context().Done():
			return false
		}
	})
}
