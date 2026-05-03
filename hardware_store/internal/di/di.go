package di

import (
	"context"
	"hardware_store/internal/app"
	"hardware_store/internal/clients/auth/grpc"
	"hardware_store/internal/clients/photo"
	"hardware_store/internal/config"
	"hardware_store/internal/kafka"
	"hardware_store/internal/logger"
	"hardware_store/internal/redis"
	"hardware_store/internal/server"
	addressservice "hardware_store/internal/service/address"
	authservice "hardware_store/internal/service/auth"
	categoryservice "hardware_store/internal/service/category"
	clientservice "hardware_store/internal/service/client"
	imagesservice "hardware_store/internal/service/images"
	productservice "hardware_store/internal/service/product"
	supplierservice "hardware_store/internal/service/supplier"
	"hardware_store/internal/storage/cache"
	"hardware_store/internal/storage/postgres"
	"hardware_store/internal/storage/postgres/address"
	"hardware_store/internal/storage/postgres/category"
	"hardware_store/internal/storage/postgres/client"
	"hardware_store/internal/storage/postgres/images"
	"hardware_store/internal/storage/postgres/product"
	"hardware_store/internal/storage/postgres/supplier"
	"hardware_store/internal/storage/postgres/tx"
	"hardware_store/internal/web"
	authhandler "hardware_store/internal/web/handler/auth"
	categoryhandler "hardware_store/internal/web/handler/category"
	clienthandler "hardware_store/internal/web/handler/client"
	imageshandler "hardware_store/internal/web/handler/images"
	producthandler "hardware_store/internal/web/handler/product"
	"hardware_store/internal/web/handler/sse"
	supplierhandler "hardware_store/internal/web/handler/supplier"
	"hardware_store/internal/web/handler/ws"
	"net/http"
	"time"

	"log/slog"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"go.uber.org/fx"
)

func NewValidator() *validator.Validate {
	return validator.New()
}
func ProvideEnv(cfg *config.Config) string {
	return cfg.Env
}

func ProvideRedisConfig(cfg *config.Config) redis.Config {
	return redis.Config{
		Addr:     cfg.RedisHost + ":" + cfg.RedisPort,
		Password: "",
		DB:       0,
	}
}

func ProvideImageCache(client *redis.Client) *cache.ImageCache {
	return cache.NewImageCache(client, "image:", 24*time.Hour)
}

// Auth service providers - используют один экземпляр для разных интерфейсов
func ProvideAuthService(auth *authservice.Auth) authhandler.AuthService {
	return auth
}

func ProvideUserService(auth *authservice.Auth) authhandler.UserService {
	return auth
}

func ProvideTokenService(auth *authservice.Auth) authhandler.TokenService {
	return auth
}

func ProvideGoogleOAuthService(auth *authservice.Auth) authhandler.GoogleOAuthService {
	return auth
}

var Module = fx.Options(
	fx.Provide(
		config.NewConfig,
		ProvideEnv,
		logger.NewLog,
		postgres.NewDB,
		tx.NewTxManager,
		NewValidator,
		ProvideRedisConfig,
		redis.NewClient,
		ProvideImageCache,

		func(cfg *config.Config, log *slog.Logger) *photo.PhotoServiceClient {
			log.Debug("Creating PhotoServiceClient",
				slog.String("base_url", cfg.PhotoServiceURL),
			)
			return photo.NewPhotoServiceClient(cfg.PhotoServiceURL, log)
		},
		/////////////
		fx.Annotate(client.NewClientRepository, fx.As(new(clientservice.ClientRepository))),
		fx.Annotate(address.NewAddressRepository, fx.As(new(addressservice.AddressRepository))),
		fx.Annotate(product.NewProductRepository, fx.As(new(productservice.ProductRepository))),
		fx.Annotate(product.NewProductRepository, fx.As(new(kafka.Storage))),
		fx.Annotate(supplier.NewSupplierRepository, fx.As(new(supplierservice.SupplierRepository))),
		fx.Annotate(images.NewImagesRepository, fx.As(new(imagesservice.ImagesRepository))),
		fx.Annotate(category.NewCategoryRepository, fx.As(new(categoryservice.CategoryRepository))),
		/////////////
		fx.Annotate(
			clientservice.NewClientService,
			fx.As(new(clientservice.ClientService)),
		),
		fx.Annotate(productservice.NewProductService,
			fx.As(new(productservice.ProductService)),
		),
		fx.Annotate(supplierservice.NewSupplierService,
			fx.As(new(supplierservice.SupplierService)),
		),
		fx.Annotate(imagesservice.NewImageService,
			fx.As(new(imagesservice.ImageService)),
		),
		fx.Annotate(addressservice.NewAddressService,
			fx.As(new(addressservice.AddressService)),
		),
		fx.Annotate(categoryservice.NewCategoryService,
			fx.As(new(categoryservice.CategoryService)),
		),
		/////////////
		clienthandler.NewClientHandler,
		imageshandler.NewImageHandler,
		producthandler.NewProductHandler,
		categoryhandler.NewCategoryHandler,
		supplierhandler.NewSupplierHandler,
		////////////
		grpc.NewClient,
		authservice.NewAuth,
		ProvideAuthService,
		ProvideUserService,
		ProvideTokenService,
		ProvideGoogleOAuthService,
		authhandler.NewAuthHandler,

		///////
		ws.NewWSHandler,
		kafka.NewConsumer,
		sse.NewSSEHandler,
		///////
		web.NewRouter,
		func(engine *gin.Engine) http.Handler {
			return engine
		},
		server.NewServer,
	),
	fx.Invoke(app.NewApp,
		postgres.AddDBLifecycle,
		func(lc fx.Lifecycle, log *slog.Logger, consumer *kafka.Consumer) {
			var cancel context.CancelFunc

			lc.Append(fx.Hook{
				OnStart: func(ctx context.Context) error {
					runCtx, c := context.WithCancel(context.Background())
					cancel = c
					go consumer.Start(runCtx)
					log.Info("kafka consumer started")
					return nil
				},
				OnStop: func(ctx context.Context) error {
					if cancel != nil {
						cancel()
					}
					if err := consumer.Close(); err != nil {
						log.Error("failed to close kafka consumer", slog.Any("error", err))
						return err
					}
					log.Info("kafka consumer stopped")
					return nil
				},
			})
		},
		func(lc fx.Lifecycle, log *slog.Logger, redisClient *redis.Client) {
			lc.Append(fx.Hook{
				OnStart: func(ctx context.Context) error {
					log.Info("redis client initialized")
					return nil
				},
				OnStop: func(ctx context.Context) error {
					if err := redisClient.Close(); err != nil {
						log.Error("failed to close redis client", slog.Any("error", err))
						return err
					}
					log.Info("redis client closed")
					return nil
				},
			})
		},
	),
)
