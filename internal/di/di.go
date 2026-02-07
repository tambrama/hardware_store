package di

import (
	"hardware_store/internal/app"
	"hardware_store/internal/config"
	"hardware_store/internal/logger"
	"hardware_store/internal/server"
	addressservice "hardware_store/internal/service/address_service"
	categoryservice "hardware_store/internal/service/category_service"
	clientservice "hardware_store/internal/service/client_service"
	imagesservice "hardware_store/internal/service/images_service"
	productservice "hardware_store/internal/service/product_service"
	supplierservice "hardware_store/internal/service/supplier_service"
	"hardware_store/internal/storage/postgres"
	"hardware_store/internal/web"
	"hardware_store/internal/web/handler"
	"net/http"

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

var Module = fx.Options(
	fx.Provide(
		config.NewConfig,
		ProvideEnv,
		logger.NewLog,
		postgres.NewDB,
		postgres.NewTxManager,
		NewValidator,
		/////////////
		postgres.NewClientRepository,
		postgres.NewProductRepository,
		postgres.NewSupplierRepository,
		postgres.NewImagesRepository,
		postgres.NewAddressRepository,
		postgres.NewCategoryRepository,
		/////////////
		clientservice.NewClientService,
		productservice.NewProductService,
		supplierservice.NewSupplierService,
		imagesservice.NewImageService,
		addressservice.NewAddressService,
		categoryservice.NewCategoryService,
		/////////////
		handler.NewClientHandler,
		handler.NewAddressHandler,
		handler.NewImageHandler,
		handler.NewProductHandler,
		handler.NewCategoryHandler,
		handler.NewSupplierHandler,
		////////////hendlers
		web.NewRouter,
		func(engine *gin.Engine) http.Handler {
			return engine
		},
		server.NewServer,
	),
	fx.Invoke(app.NewApp,
		postgres.AddDBLifecycle),
)
