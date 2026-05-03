package ws

import (
	"encoding/json"
	"fmt"
	"hardware_store/internal/model/product"
	"log/slog"
	"net/http"
	"sync"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool { return true },
}

type WSHandler struct {
	client map[*websocket.Conn]bool
	mu     sync.RWMutex
	logger *slog.Logger
}

func NewWSHandler(logger *slog.Logger) *WSHandler {
	return &WSHandler{
		client: make(map[*websocket.Conn]bool),
		logger: logger,
	}
}

func (h *WSHandler) HandleWS(c *gin.Context) {
	conn, err := upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		h.logger.Error("upgrade failed", slog.Any("error", err))
		return
	}

	h.mu.Lock()
	h.client[conn] = true
	h.mu.Unlock()

	h.logger.Info("🔌 WebSocket connected", slog.String("ip", c.ClientIP()))

	defer func() {
		h.mu.Lock()
		delete(h.client, conn)
		h.mu.Unlock()
		conn.Close()
		h.logger.Info("🔌 WebSocket disconnected", slog.String("ip", c.ClientIP()))
	}()


	for {
		_, _, err := conn.ReadMessage()
		if err != nil {
			h.logger.Debug("WebSocket read error", slog.Any("error", err), slog.String("ip", c.ClientIP()))
			break
		}
	}
}

func (h *WSHandler) Broadcast(event *product.ProductUpdate) {
	h.mu.RLock()
	defer h.mu.RUnlock()

	data, err := json.Marshal(event)
	if err != nil {
		h.logger.Error("failed to marshal event", slog.Any("error", err))
		return
	}

	for conn := range h.client {
		if err := conn.WriteMessage(websocket.TextMessage, data); err != nil {
			h.logger.Debug("failed to write message to WebSocket",
				slog.Any("error", err),
				slog.String("error_type", fmt.Sprintf("%T", err)))
		}
	}
}
