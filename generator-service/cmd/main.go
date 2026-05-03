package main

import (
	"generator-service/internal/di"

	"go.uber.org/fx"
)

func main() {
	fx.New(
		di.Module,
	).Run()
}
