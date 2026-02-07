package web

import (
	"hardware_store/internal/web/handler"

	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"

	_ "hardware_store/docs"
)

func NewRouter(client *handler.ClientHandler, product *handler.ProductHandler,
	address *handler.AddressHandler, image *handler.ImageHandler,
	category *handler.CategoryHandler, supplier *handler.SupplierHandler) *gin.Engine {
	r := gin.Default()
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
	api := r.Group("/api/v1")
	{
		client.Register(api)
		product.Register(api)
		address.Register(api)
		image.Register(api)
		category.Register(api)
		supplier.Register(api)
	}
	return r
}
