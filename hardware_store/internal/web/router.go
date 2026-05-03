package web

import (
	"hardware_store/internal/web/handler/auth"
	"hardware_store/internal/web/handler/category"
	"hardware_store/internal/web/handler/client"
	"hardware_store/internal/web/handler/images"
	"hardware_store/internal/web/handler/product"
	"hardware_store/internal/web/handler/sse"
	"hardware_store/internal/web/handler/supplier"
	"hardware_store/internal/web/handler/ws"
	"hardware_store/internal/web/middleware"
	"log/slog"

	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"

	_ "hardware_store/docs"
)

func NewRouter(log *slog.Logger, client *client.ClientHandler, product *product.ProductHandler,
	image *images.ImageHandler,
	category *category.CategoryHandler, supplier *supplier.SupplierHandler, auth auth.TokenService, authHandler *auth.AuthHandler,
	wsHandler *ws.WSHandler, sseHandler *sse.SSEHandler) *gin.Engine {
	r := gin.Default()

	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	// WebSocket endpoint
	r.GET("/ws", wsHandler.HandleWS)

	// SSE
	r.GET("/sse", sseHandler.HandleSSE)

	api := r.Group("/api/v1")
	{
		authHandler.Register(api)

		protected := api.Group("")
		protected.Use(middleware.NewMiddleware(log, auth))
		{
			protected.POST("/auth/logout", authHandler.Logout)
			client.Register(protected)
			product.Register(protected)
			image.Register(protected)
			category.Register(protected)
			supplier.Register(protected)
		}

	}
	return r
}
