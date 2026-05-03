package producer

import (
	"context"
	"encoding/json"
	"fmt"
	"log/slog"
	"time"

	"generator-service/internal/models"

	"github.com/segmentio/kafka-go"
)

type KafkaProducer struct {
	writer *kafka.Writer
	logger *slog.Logger
}

// TestMessage для проверки подключения к Kafka
type TestMessage struct {
	ID string `json:"id"`
}

// NewKafkaProducer создаёт новый Kafka producer
func NewKafkaProducer(brokers []string, topic string, logger *slog.Logger) *KafkaProducer {
	writer := &kafka.Writer{
		Addr:         kafka.TCP(brokers...),
		Topic:        topic,
		Balancer:     &kafka.LeastBytes{},
		WriteTimeout: 10 * time.Second,
		ReadTimeout:  10 * time.Second,
		RequiredAcks: kafka.RequireAll, // Дождаться подтверждения от всех реплик
		Compression:  kafka.Snappy,     // Сжатие сообщений для экономии
		MaxAttempts:  3,                // Количество попыток переотправки
	}

	return &KafkaProducer{
		writer: writer,
		logger: logger,
	}
}

// SendTestMessage отправляет тестовое сообщение для проверки подключения
func (kp *KafkaProducer) SendTestMessage(ctx context.Context, test *TestMessage) error {
	message, err := json.Marshal(test)
	if err != nil {
		kp.logger.Error("failed to marshal test message", slog.Any("error", err))
		return fmt.Errorf("marshal error: %w", err)
	}

	err = kp.writer.WriteMessages(ctx, kafka.Message{
		Key:   []byte("test"),
		Value: message,
	})

	if err != nil {
		kp.logger.Error("failed to send test message to kafka", slog.Any("error", err))
		return fmt.Errorf("send test message error: %w", err)
	}

	kp.logger.Debug("test message sent to kafka")
	return nil
}

// SendProductUpdate отправляет событие об изменении товара в Kafka с retry логикой
func (kp *KafkaProducer) SendProductUpdate(ctx context.Context, update *models.ProductUpdate) error {
	message, err := json.Marshal(update)
	if err != nil {
		kp.logger.Error("failed to marshal event", slog.Any("error", err))
		return fmt.Errorf("marshal error: %w", err)
	}

	// Добавляем контекст с таймаутом, если его нет
	deadline, ok := ctx.Deadline()
	if !ok {
		var cancel context.CancelFunc
		ctx, cancel = context.WithTimeout(ctx, 15*time.Second)
		defer cancel()
	} else {
		// Если уже установлен дедлайн, убедимся что его достаточно
		remaining := time.Until(deadline)
		if remaining < 5*time.Second {
			kp.logger.Warn("context deadline too soon", slog.Duration("remaining", remaining))
		}
	}

	// Отправляем в Kafka с retry логикой
	err = kp.writer.WriteMessages(ctx, kafka.Message{
		Key:   []byte(update.ProductID), // Ключ - ID товара для партиционирования
		Value: message,                  // Значение - JSON событие
	})

	if err != nil {
		kp.logger.Error("failed to send message to kafka",
			slog.Any("error", err),
			slog.String("product_id", update.ProductID),
		)
		return fmt.Errorf("send to kafka error: %w", err)
	}

	kp.logger.Info("product update sent to kafka",
		slog.String("product_id", update.ProductID),
		slog.String("product_name", update.ProductName),
		slog.Float64("new_price", update.NewPrice),
		slog.Int("new_stock", update.NewStock),
	)

	return nil
}

// Close закрывает подключение к Kafka
func (kp *KafkaProducer) Close() error {
	if err := kp.writer.Close(); err != nil {
		kp.logger.Error("failed to close kafka writer", slog.Any("error", err))
		return fmt.Errorf("close error: %w", err)
	}
	return nil
}
