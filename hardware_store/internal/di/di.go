package di

import (
	"hardware_store/internal/app"
	"hardware_store/internal/config"
	"hardware_store/internal/logger"
	"hardware_store/internal/server"
	addressservice "hardware_store/internal/service/address"
	categoryservice "hardware_store/internal/service/category"
	clientservice "hardware_store/internal/service/client"
	imagesservice "hardware_store/internal/service/images"
	productservice "hardware_store/internal/service/product"
	supplierservice "hardware_store/internal/service/supplier"
	"hardware_store/internal/storage/postgres"
	"hardware_store/internal/storage/postgres/address"
	"hardware_store/internal/storage/postgres/category"
	"hardware_store/internal/storage/postgres/client"
	"hardware_store/internal/storage/postgres/images"
	"hardware_store/internal/storage/postgres/product"
	"hardware_store/internal/storage/postgres/supplier"
	"hardware_store/internal/storage/postgres/tx"
	"hardware_store/internal/web"
	categoryhandler "hardware_store/internal/web/handler/category"
	clienthandler "hardware_store/internal/web/handler/client"
	imageshandler "hardware_store/internal/web/handler/images"
	producthandler "hardware_store/internal/web/handler/product"
	supplierhandler "hardware_store/internal/web/handler/supplier"
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
		tx.NewTxManager,
		NewValidator,
		/////////////
		fx.Annotate(client.NewClientRepository, fx.As(new(clientservice.ClientRepository))),
		fx.Annotate(address.NewAddressRepository, fx.As(new(addressservice.AddressRepository))),
		fx.Annotate(product.NewProductRepository, fx.As(new(productservice.ProductRepository))),
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
		web.NewRouter,
		func(engine *gin.Engine) http.Handler {
			return engine
		},
		server.NewServer,
	),
	fx.Invoke(app.NewApp,
		postgres.AddDBLifecycle),
)
