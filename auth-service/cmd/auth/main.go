package main

import (
	"auth-service/internal/di"

	"go.uber.org/fx"
)

func main() {
	fx.New(
		di.Module,
	).Run()
}
