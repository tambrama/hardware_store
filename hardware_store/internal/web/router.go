package web

import (
	"hardware_store/internal/web/handler/category"
	"hardware_store/internal/web/handler/client"
	"hardware_store/internal/web/handler/images"
	"hardware_store/internal/web/handler/product"
	"hardware_store/internal/web/handler/supplier"

	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"

	_ "hardware_store/docs"
)

func NewRouter(client *client.ClientHandler, product *product.ProductHandler,
	image *images.ImageHandler,
	category *category.CategoryHandler, supplier *supplier.SupplierHandler) *gin.Engine {
	r := gin.Default()

	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
	api := r.Group("/api/v1")
	{
		client.Register(api)
		product.Register(api)
		image.Register(api)
		category.Register(api)
		supplier.Register(api)
	}
	return r
}
