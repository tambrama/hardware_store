package main

import (
	"hardware_store/internal/di"

	"go.uber.org/fx"
)

// @title Hardware Store API
// @version 1.0
// @description REST API для магазина бытовой техники
// @BasePath /api/v1

// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization
// @description Введите токен в формате: "Bearer <ваш_токен>"

// OAuth 2.0 (Google) 
// @securitydefinitions.oauth2.accessCode OAuth2
// @tokenUrl https://oauth2.googleapis.com/token
// @authorizationurl https://accounts.google.com/o/oauth2/v2/auth
// @scope.email Просмотр вашего email
// @scope.profile Просмотр вашего профиля

// @security BearerAuth
// @security OAuth2
func main() {
	fx.New(
		di.Module,
	).Run()
}
