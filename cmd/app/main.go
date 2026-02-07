package main

import (
	"hardware_store/internal/di"

	"go.uber.org/fx"
)

// @title Hardware Store API
// @version 1.0
// @description REST API для магазина бытовой техники
// @host localhost:8081
// @BasePath /api/v1
func main() {
	fx.New(
		di.Module,
	).Run()
}
