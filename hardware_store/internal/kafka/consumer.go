package kafka

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"time"

	"hardware_store/internal/config"
	"hardware_store/internal/model/product"
	"hardware_store/internal/web/handler/ws"
	"log/slog"
	"sync"

	"github.com/segmentio/kafka-go"
)

type Storage interface {
	UpdateProduct(ctx context.Context, id string, newCol int, newPrice float64) error
}

type Consumer struct {
	reader    *kafka.Reader
	db        Storage
	logger    *slog.Logger
	wsHandler *ws.WSHandler
	sseHub    *sseHub
}

type sseHub struct {
	sseClients map[chan product.ProductUpdate]bool
	sseMutex   sync.RWMutex
}

func NewHub() *sseHub {
	return &sseHub{
		sseClients: make(map[chan product.ProductUpdate]bool),
	}
}

func NewConsumer(cfg *config.Config, db Storage, logger *slog.Logger, wsHandler *ws.WSHandler) *Consumer {
	return &Consumer{
		reader: kafka.NewReader(kafka.ReaderConfig{
			Brokers:        cfg.KafkaBrokers,
			Topic:          cfg.KafkaTopic,
			GroupID:        "hardware-store-group",
			MinBytes:       10e3,
			MaxBytes:       10e6,
			CommitInterval: time.Second,            // Коммитить каждую секунду
			StartOffset:    kafka.LastOffset,       // Читать с последнего оффсета
			MaxWait:        10 * time.Second,       // Дольше ждём данных (макс 10 сек)
			ReadBackoffMin: 100 * time.Millisecond, // Минимальная задержка перед повтором
			ReadBackoffMax: 1 * time.Second,        // Максимальная задержка перед повтором
		}),
		db:        db,
		logger:    logger,
		sseHub:    NewHub(),
		wsHandler: wsHandler,
	}
}

func (c *Consumer) Start(ctx context.Context) {
	c.logger.Info("📥 starting kafka consumer",
		slog.Any("brokers", c.reader.Config().Brokers),
		slog.String("topic", c.reader.Config().Topic),
		slog.String("group_id", c.reader.Config().GroupID),
	)

	for {
		select {
		case <-ctx.Done():
			c.logger.Info("consumer context cancelled, stopping")
			return
		default:
			// Используем длительный timeout (120 сек) для получения сообщения
			// Это нормально, если timeout истечет — значит, сообщений нет
			fetchCtx, cancel := context.WithTimeout(context.Background(), 120*time.Second)
			msg, err := c.reader.FetchMessage(fetchCtx)
			cancel()

			if err != nil {
				// Проверяем, закрыт ли основной контекст
				if ctx.Err() != nil {
					c.logger.Info("context cancelled, stopping consumer")
					return
				}

				// Context deadline exceeded — это нормально, нет новых сообщений
				var deadlineExceeded interface{ Timeout() bool }
				if errors.As(err, &deadlineExceeded) && deadlineExceeded.Timeout() {
					c.logger.Debug("kafka fetch timeout (no messages available)",
						slog.Int64("wait_time_ms", 120000),
					)
				} else {
					// Это реальная ошибка
					c.logger.Error("failed to fetch message from kafka",
						slog.Any("error", err),
						slog.String("error_type", fmt.Sprintf("%T", err)),
					)
				}

				select {
				case <-ctx.Done():
					return
				case <-time.After(time.Second):
					continue
				}
			}

			// Парсим сообщение
			var event product.ProductUpdate
			if err := json.Unmarshal(msg.Value, &event); err != nil {
				c.logger.Error("failed to unmarshal event from kafka",
					slog.Any("error", err),
					slog.String("raw_value", string(msg.Value)),
				)
				// Коммитим сообщение даже если ошибка парсинга, чтобы не зависнуть
				if err := c.reader.CommitMessages(ctx, msg); err != nil {
					c.logger.Error("failed to commit message after unmarshal error",
						slog.Any("error", err),
					)
				}
				continue
			}

			c.logger.Debug("received event from kafka",
				slog.String("product_id", event.ProductID),
				slog.Float64("new_price", event.NewPrice),
				slog.Int("new_stock", event.NewStock),
			)

			if err := c.db.UpdateProduct(ctx, event.ProductID, event.NewStock, event.NewPrice); err != nil {
				c.logger.Error("failed to update product in db",
					slog.Any("error", err),
					slog.String("product_id", event.ProductID),
				)
				continue
			}

			if c.wsHandler != nil {
				c.wsHandler.Broadcast(&event)
			}

			if c.sseHub != nil {
				c.sseHub.broadcast(&event)
			}

			if err := c.reader.CommitMessages(ctx, msg); err != nil {
				c.logger.Error("failed to commit message",
					slog.Any("error", err),
					slog.String("product_id", event.ProductID),
				)

			}

			c.logger.Info("product updated successfully",
				slog.String("product_id", event.ProductID),
				slog.Float64("new_price", event.NewPrice),
				slog.Int("new_stock", event.NewStock),
			)
		}
	}
}

func (h *sseHub) RegisterSSE(ch chan product.ProductUpdate) {
	h.sseMutex.Lock()
	h.sseClients[ch] = true
	h.sseMutex.Unlock()
}

func (h *sseHub) UnregisterSSE(ch chan product.ProductUpdate) {
	h.sseMutex.Lock()
	delete(h.sseClients, ch)
	h.sseMutex.Unlock()
}

func (h *sseHub) broadcast(event *product.ProductUpdate) {
	h.sseMutex.RLock()
	for ch := range h.sseClients {
		select {
		case ch <- *event:
		default:
		}
	}
	h.sseMutex.RUnlock()
}

func (c *Consumer) GetHub() *sseHub {
	return c.sseHub
}

func (c *Consumer) Close() error {
	c.logger.Info("closing kafka consumer")
	if err := c.reader.Close(); err != nil {
		c.logger.Error("failed to close kafka reader", slog.Any("error", err))
		return err
	}
	c.logger.Info("kafka consumer closed successfully")
	return nil
}
