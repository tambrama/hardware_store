package di

import (
	"context"
	"generator-service/internal/config"
	"generator-service/internal/generator"
	"generator-service/internal/producer"
	"generator-service/internal/storage"
	"log"
	"log/slog"
	"os"
	"time"

	"go.uber.org/fx"
)

const (
	envLocal      = "local"
	envProduction = "production"
)

func setupLogger(cfg *config.Config) *slog.Logger {
	var log *slog.Logger
	switch cfg.Env {
	case envLocal:
		log = slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}))
	case envProduction:
		log = slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelInfo}))
	default:
		log = slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}))
	}
	return log
}

func newKafkaProducer(cfg *config.Config, logger *slog.Logger, lc fx.Lifecycle) *producer.KafkaProducer {
	p := producer.NewKafkaProducer(cfg.KafkaBrokers, cfg.KafkaTopic, logger)

	// Проверяем подключение к Kafka при инициализации
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	testUpdate := &producer.TestMessage{ID: "test"}
	if err := p.SendTestMessage(ctx, testUpdate); err != nil {
		log.Fatal("failed to connect to kafka on startup", slog.Any("error", err))
	}
	logger.Info("kafka producer initialized successfully", slog.Any("brokers", cfg.KafkaBrokers))

	lc.Append(fx.Hook{
		OnStop: func(ctx context.Context) error {
			logger.Info("closing kafka producer")
			return p.Close()
		},
	})
	return p
}

type ShopRepo interface{ generator.Storage }
type HwRepo interface{ generator.Storage }

func newShopRepo(cfg *config.Config, logger *slog.Logger, lc fx.Lifecycle) ShopRepo {
	logger.Info("🔌 connecting to shop_db", slog.String("dsn", cfg.ShopDB))
	pool := storage.NewDB(cfg.ShopDB, lc, logger)
	return storage.NewStorage(pool, logger)
}

func newHwRepo(cfg *config.Config, logger *slog.Logger, lc fx.Lifecycle) HwRepo {
	logger.Info("🔌 connecting to hardware_db", slog.String("dsn", cfg.HardwareDB))
	pool := storage.NewDB(cfg.HardwareDB, lc, logger)
	return storage.NewStorage(pool, logger)
}

func newGenerator(shopRepo ShopRepo, hwRepo HwRepo, kafkaProducer *producer.KafkaProducer, cfg *config.Config, logger *slog.Logger) *generator.GeneratorService {
	return generator.NewGeneratorService(shopRepo, hwRepo, kafkaProducer, cfg.GeneratorInterval, cfg.PriceChangePercent, cfg.StockChangeAmount, logger)
}

var Module = fx.Options(
	fx.Provide(
		config.NewConfig,
		setupLogger,
		// storage.NewDB,
		// fx.Annotate(
		// 	storage.NewStorage,
		// 	fx.As(new(generator.Storage)),
		// ),

		newKafkaProducer,
		newGenerator,
		newShopRepo,
		newHwRepo,
	),
	fx.Invoke(func(
		lc fx.Lifecycle,
		logger *slog.Logger,
		gen *generator.GeneratorService,
	) {
		runCtx, cancel := context.WithCancel(context.Background())

		lc.Append(fx.Hook{
			OnStart: func(ctx context.Context) error {
				logger.Info("starting")
				go gen.Start(runCtx)
				return nil
			},
			OnStop: func(ctx context.Context) error {
				logger.Info("stopping")
				cancel()
				return nil
			},
		})
	}),
)
