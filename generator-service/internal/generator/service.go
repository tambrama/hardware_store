package generator

import (
	"context"
	"fmt"
	"log/slog"
	"math"
	"math/rand"
	"time"

	"generator-service/internal/models"
	"generator-service/internal/producer"
)

type Storage interface {
	Get(ctx context.Context) (*models.Product, error)
}

type GeneratorService struct {
	shopRepo           Storage
	hwRepo             Storage
	kafkaProducer      *producer.KafkaProducer
	interval           time.Duration
	priceChangePercent float64
	stockChangeAmount  int
	logger             *slog.Logger
}

// NewGeneratorService создаёт новый сервис генератора
func NewGeneratorService(
	shopRepo           Storage,
	hwRepo             Storage,
	kafkaProducer *producer.KafkaProducer,
	interval time.Duration,
	priceChangePercent float64,
	stockChangeAmount int,
	logger *slog.Logger,
) *GeneratorService {
	return &GeneratorService{
		shopRepo:                 shopRepo,
		hwRepo:                   hwRepo,
		kafkaProducer:      kafkaProducer,
		interval:           interval,
		priceChangePercent: priceChangePercent,
		stockChangeAmount:  stockChangeAmount,
		logger:             logger,
	}
}

// Start запускает генератор событий
func (gs *GeneratorService) Start(ctx context.Context) {
	ticker := time.NewTicker(gs.interval)
	defer ticker.Stop()

	gs.logger.Info("generator service started", slog.Duration("interval", gs.interval))

	for {
		select {
		case <-ctx.Done():
			gs.logger.Info("generator service stopped")
			return
		case <-ticker.C:
			// Каждые N секунд генерируем событие
			if err := gs.generateEvent(ctx); err != nil {
				gs.logger.Error("failed to generate event", slog.Any("error", err))
			}
		}
	}
}

// generateEvent генерирует одно событие об изменении товара
func (gs *GeneratorService) generateEvent(ctx context.Context) error {
	useShop := rand.Float64() < 0.5

	var repo Storage
	var source models.ProductSource
	if useShop {
		repo = gs.shopRepo
		source = models.SourceShopAPI
	} else {
		repo = gs.hwRepo
		source = models.SourceHardware
	}

	product, err := repo.Get(ctx)
	if err != nil {
		gs.logger.Warn("failed to get product from source",
			slog.String("source", string(source)),
			slog.Any("error", err),
		)
		return nil
	}

	if product == nil {
		gs.logger.Warn("no products found", slog.String("source", string(source)))
		return nil
	}

	oldPrice := product.Price
	priceChange := oldPrice * (gs.priceChangePercent / 100.0)

	if rand.Float64() > 0.5 {
		product.Price += priceChange
	} else {
		product.Price -= priceChange
	}

	newPrice := math.Round(product.Price*100) / 100

	oldStock := product.Stock

	if rand.Float64() > 0.5 {
		product.Stock += gs.stockChangeAmount
	} else {
		product.Stock -= gs.stockChangeAmount
		if product.Stock < 0 {
			product.Stock = 0
		}
	}
	newStock := product.Stock

	event := &models.ProductUpdate{
		ProductID:   product.ID,
		ProductName: product.Name,
		OldPrice:    oldPrice,
		NewPrice:    newPrice,
		OldStock:    oldStock,
		NewStock:    newStock,
		Source:      string(source), 
		Timestamp:   time.Now(),
	}

	if err := gs.kafkaProducer.SendProductUpdate(ctx, event); err != nil {
		return fmt.Errorf("failed to send event to kafka: %w", err)
	}

	gs.logger.Debug("event sent",
		slog.String("product_id", product.ID),
		slog.String("source", string(source)),
		slog.Float64("new_price", newPrice),
		slog.Int("new_stock", newStock),
	)
	
	return nil
}
